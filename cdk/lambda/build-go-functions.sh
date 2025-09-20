#!/bin/bash

# Build script for Go Lambda functions

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Move out of GOPATH if we're in it
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org,direct

# Lambda関数リスト
FUNCTIONS="users accounts transactions categories budgets reports auth notifications"

echo "Building Go Lambda functions..."

for func in $FUNCTIONS; do
    echo "Building $func..."
    
    # ディレクトリが存在することを確認
    if [ -d "$func" ]; then
        cd "$func"
        
        # main.goが存在する場合のみビルド
        if [ -f "main.go" ]; then
            # go.modがない場合は作成
            if [ ! -f "go.mod" ]; then
                go mod init finsight/$func
            fi
            
            # 共通ライブラリへのreplace追加
            if ! grep -q "replace finsight/common" go.mod; then
                echo "" >> go.mod
                echo "replace finsight/common => ../common" >> go.mod
            fi
            
            # 依存関係の更新
            GO111MODULE=on go mod tidy || true
            
            # ビルド
            GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
            echo "✓ Built $func"
        else
            # main.goが存在しない場合は、プレースホルダーを作成
            echo "Creating placeholder for $func..."
            cat > main.go << EOF
package main

import (
    "context"
    "encoding/json"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    response := map[string]string{
        "message": "$func endpoint placeholder",
    }
    
    body, _ := json.Marshal(response)
    
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Headers: map[string]string{
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "*",
        },
        Body: string(body),
    }, nil
}

func main() {
    lambda.Start(handler)
}
EOF
            # go.modファイルを作成
            go mod init finsight/$func
            go get github.com/aws/aws-lambda-go/events
            go get github.com/aws/aws-lambda-go/lambda
            
            # ビルド
            GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
            echo "✓ Created and built placeholder for $func"
        fi
        
        cd ..
    fi
done

echo "All Lambda functions built successfully!"