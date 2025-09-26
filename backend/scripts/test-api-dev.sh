#!/bin/bash

# FinanceTracker API 開発環境テストスクリプト
# 使用方法: ./test-api-dev.sh

BASE_URL="http://localhost:8080"
API_URL="${BASE_URL}/api/v1"
DEV_USER_ID="auth0|dev-test-user-123"

# 色付き出力
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🚀 FinanceTracker API 開発環境テスト開始"
echo "========================================="
echo "開発用ユーザーID: ${DEV_USER_ID}"
echo ""

# 1. ヘルスチェック
echo -e "\n${YELLOW}1. ヘルスチェック${NC}"
curl -s ${BASE_URL}/health | jq .

# 2. 現在のユーザー情報取得
echo -e "\n${YELLOW}2. 現在のユーザー情報（初回は404）${NC}"
curl -s -H "X-Dev-User-ID: ${DEV_USER_ID}" ${API_URL}/users/me | jq .

# 3. ユーザー情報更新（自動作成）
echo -e "\n${YELLOW}3. ユーザー情報更新（自動作成）${NC}"
curl -s -X PUT -H "X-Dev-User-ID: ${DEV_USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "開発テストユーザー",
    "email": "dev@example.com"
  }' \
  ${API_URL}/users/me | jq .

# 4. カテゴリマスター一覧
echo -e "\n${YELLOW}4. カテゴリマスター一覧（収入）${NC}"
curl -s -H "X-Dev-User-ID: ${DEV_USER_ID}" \
  "${API_URL}/categories/master?category_type=income" | jq .

echo -e "\n${YELLOW}5. カテゴリマスター一覧（支出）${NC}"
curl -s -H "X-Dev-User-ID: ${DEV_USER_ID}" \
  "${API_URL}/categories/master?category_type=expense" | jq .

# 6. 口座作成
echo -e "\n${YELLOW}6. 口座作成${NC}"
ACCOUNT_RESPONSE=$(curl -s -X POST -H "X-Dev-User-ID: ${DEV_USER_ID}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "メイン口座",
    "account_type": "checking",
    "currency": "JPY",
    "initial_balance": 100000
  }' \
  ${API_URL}/accounts)
echo "$ACCOUNT_RESPONSE" | jq .
ACCOUNT_ID=$(echo "$ACCOUNT_RESPONSE" | jq -r .id)

# 7. 口座一覧
echo -e "\n${YELLOW}7. 口座一覧${NC}"
curl -s -H "X-Dev-User-ID: ${DEV_USER_ID}" ${API_URL}/accounts | jq .

# 8. カテゴリ作成（マスターカテゴリから）
echo -e "\n${YELLOW}8. カテゴリ作成準備${NC}"
# まずカテゴリマスターのIDを取得
MASTER_RESPONSE=$(curl -s -H "X-Dev-User-ID: ${DEV_USER_ID}" \
  "${API_URL}/categories/master?category_type=expense")
MASTER_ID=$(echo "$MASTER_RESPONSE" | jq -r '.category_masters[0].id')

if [ "$MASTER_ID" != "null" ]; then
  echo "マスターカテゴリID: $MASTER_ID"
  
  # カテゴリ作成
  CATEGORY_RESPONSE=$(curl -s -X POST -H "X-Dev-User-ID: ${DEV_USER_ID}" \
    -H "Content-Type: application/json" \
    -d "{
      \"category_master_id\": \"$MASTER_ID\",
      \"custom_name\": \"ランチ代\"
    }" \
    ${API_URL}/categories)
  echo "$CATEGORY_RESPONSE" | jq .
  CATEGORY_ID=$(echo "$CATEGORY_RESPONSE" | jq -r .id)
else
  echo -e "${RED}カテゴリマスターが見つかりません${NC}"
  CATEGORY_ID=""
fi

# 9. トランザクション作成
if [ -n "$ACCOUNT_ID" ] && [ -n "$CATEGORY_ID" ] && [ "$ACCOUNT_ID" != "null" ] && [ "$CATEGORY_ID" != "null" ]; then
  echo -e "\n${YELLOW}9. トランザクション作成${NC}"
  curl -s -X POST -H "X-Dev-User-ID: ${DEV_USER_ID}" \
    -H "Content-Type: application/json" \
    -d "{
      \"account_id\": \"$ACCOUNT_ID\",
      \"category_id\": \"$CATEGORY_ID\",
      \"transaction_type\": \"expense\",
      \"amount\": 1200,
      \"description\": \"カフェでランチ\",
      \"transaction_date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }" \
    ${API_URL}/transactions | jq .
fi

# 10. 月次サマリー
echo -e "\n${YELLOW}10. 月次サマリー${NC}"
YEAR=$(date +%Y)
MONTH=$(date +%m | sed 's/^0*//')  # 先頭の0を除去
curl -s -H "X-Dev-User-ID: ${DEV_USER_ID}" \
  "${API_URL}/transactions/summary/monthly?year=$YEAR&month=$MONTH" | jq .

echo -e "\n${GREEN}テスト完了！${NC}"
echo ""
echo "次のステップ:"
echo "1. Postmanで詳細なテストを実行"
echo "2. 実際のAuth0トークンでテスト"
echo "3. フロントエンドとの統合テスト"