# 開発用 Dockerfile
FROM node:20-alpine

# 作業ディレクトリを設定
WORKDIR /app

# 開発環境を指定
ENV NODE_ENV=development
ENV NEXT_TELEMETRY_DISABLED=1

# pnpm をインストール（より高速）
RUN npm install -g pnpm

# package.json と package-lock.json をコピー
COPY package*.json ./
COPY pnpm-lock.yaml* ./

# 依存関係をインストール（devDependencies も含む）
RUN npm ci || pnpm install --frozen-lockfile || npm install

# ポートを解放
EXPOSE 3000

# 開発サーバーをホットリロード有効で起動
CMD ["npm", "run", "dev"]