#!/bin/bash

set -e

echo "🔍 Onion Architecture Dependency Checker"
echo "========================================"
echo ""

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# プロジェクトルートから実行されているか確認
if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found. Please run this script from the backend directory.${NC}"
    exit 1
fi

# アーキテクチャレイヤーの定義
DOMAIN_PKG="financetracker/internal/domain"
APP_PKG="financetracker/internal/application"
INFRA_PKG="financetracker/internal/infrastructure"
INTERFACE_PKG="financetracker/internal/interface"

violations=0
violation_details=""

# 進捗表示関数
print_checking() {
    echo -e "${YELLOW}Checking:${NC} $1"
}

print_success() {
    echo -e "${GREEN}✅${NC} $1"
}

print_error() {
    echo -e "${RED}❌${NC} $1"
}

# 依存関係をチェックする関数
check_dependencies() {
    local layer_name=$1
    local layer_path=$2
    local forbidden_patterns=$3
    local layer_description=$4
    
    print_checking "$layer_name layer dependencies..."
    
    # パッケージが存在するか確認
    if [ ! -d "$layer_path" ]; then
        echo "  ⚠️  $layer_name layer not found at $layer_path"
        return 0
    fi
    
    # 依存関係を取得
    local deps=$(go list -f '{{join .Deps "\n"}}' ./$layer_path/... 2>/dev/null | grep -E "$forbidden_patterns" || true)
    
    if [ -n "$deps" ]; then
        print_error "$layer_name layer has forbidden dependencies:"
        echo "$deps" | while read -r dep; do
            echo "    └─ $dep"
            # どのファイルが問題のあるインポートをしているか特定
            local files=$(grep -r "\"$dep\"" $layer_path --include="*.go" 2>/dev/null | cut -d: -f1 | sort -u || true)
            if [ -n "$files" ]; then
                echo "$files" | while read -r file; do
                    echo "       📄 $file"
                done
            fi
        done
        violations=$((violations + 1))
        violation_details="${violation_details}\n${layer_name}: ${layer_description}"
    else
        print_success "$layer_name layer: No forbidden dependencies found"
    fi
    echo ""
}

# 循環依存をチェックする関数
check_circular_dependencies() {
    print_checking "Circular dependencies..."
    
    # go mod graphで循環依存を検出（内部パッケージのみ）
    local circular=$(go list -f '{{.ImportPath}} -> {{join .Imports " "}}' ./... 2>/dev/null | \
        grep "internal/" | \
        awk '{for(i=3;i<=NF;i++) if($i ~ /internal\//) print $1 " -> " $i}' | \
        sort -u || true)
    
    if [ -n "$circular" ]; then
        # 簡易的な循環依存チェック（A->B かつ B->A のパターンを検出）
        local has_circular=false
        while IFS= read -r line; do
            pkg1=$(echo "$line" | cut -d' ' -f1)
            pkg2=$(echo "$line" | cut -d' ' -f3)
            reverse_check=$(echo "$circular" | grep "^$pkg2 -> $pkg1$" || true)
            if [ -n "$reverse_check" ]; then
                if [ "$has_circular" = false ]; then
                    print_error "Circular dependencies detected:"
                    has_circular=true
                    violations=$((violations + 1))
                fi
                echo "    └─ $pkg1 ⟷ $pkg2"
            fi
        done <<< "$circular"
        
        if [ "$has_circular" = false ]; then
            print_success "No circular dependencies found"
        fi
    else
        print_success "No circular dependencies found"
    fi
    echo ""
}

# パッケージ構造の概要を表示
show_architecture_overview() {
    echo "📊 Architecture Overview"
    echo "========================"
    echo ""
    echo "Onion Architecture Layers:"
    echo "  1. Domain (Core) - Business logic and entities"
    echo "  2. Application - Use cases and application services"
    echo "  3. Infrastructure - External concerns (DB, APIs)"
    echo "  4. Interface - Controllers and middleware"
    echo ""
    echo "Dependency Rules:"
    echo "  • Domain → No dependencies on other layers"
    echo "  • Application → Can depend on Domain only"
    echo "  • Infrastructure → Can depend on Domain and Application"
    echo "  • Interface → Can depend on all inner layers"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
}

# メイン処理
main() {
    show_architecture_overview
    
    echo "🔍 Running Architecture Checks..."
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # Domain層のチェック（他の層に依存してはいけない）
    check_dependencies \
        "Domain" \
        "internal/domain" \
        "application|infrastructure|interface" \
        "Should not depend on Application, Infrastructure, or Interface layers"
    
    # Application層のチェック（Infrastructure/Interface層に依存してはいけない）
    check_dependencies \
        "Application" \
        "internal/application" \
        "infrastructure|interface" \
        "Should not depend on Infrastructure or Interface layers"
    
    # Infrastructure層のチェック（Interface層に依存してはいけない）
    # interface{}型は除外
    check_dependencies \
        "Infrastructure" \
        "internal/infrastructure" \
        "internal/interface[^{}]" \
        "Should not depend on Interface layer"
    
    # 循環依存のチェック
    check_circular_dependencies
    
    # 結果のサマリー
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📋 Summary"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    if [ $violations -eq 0 ]; then
        echo -e "${GREEN}✅ All architecture checks passed!${NC}"
        echo ""
        echo "Your code follows the Onion Architecture principles correctly."
    else
        echo -e "${RED}❌ Found $violations architecture violation(s)${NC}"
        echo ""
        echo "Violations found:"
        echo -e "$violation_details"
        echo ""
        echo "Please fix these violations to maintain clean architecture."
        exit 1
    fi
}

# スクリプト実行
main