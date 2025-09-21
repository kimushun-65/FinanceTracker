FROM golang:1.23-alpine

# 必要なツールのインストール
RUN apk update && apk add --no-cache \
    git \
    curl \
    postgresql-client \
    ca-certificates \
    tzdata \
    make \
    && rm -rf /var/cache/apk/*

# Atlas CLIのインストール
RUN curl -sSf https://atlasgo.sh | sh

# golangci-lintのインストール (v1.64.8)
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.8

# airのインストール (ホットリロード用) - v1.61.1を指定（Go 1.23対応）
RUN go install github.com/air-verse/air@v1.61.1

# タイムゾーン設定
ENV TZ=Asia/Tokyo
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# 作業ディレクトリ設定
WORKDIR /app

# Go modulesの依存関係を先にコピー（キャッシュ効率化）
COPY go.mod go.sum ./
RUN go mod download

# ソースコード全体をコピー
COPY . .

# 開発用エントリーポイント（airでホットリロード）
CMD ["air", "-c", ".air.toml"]