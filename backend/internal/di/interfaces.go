package di

// DIContainer はDIコンテナのインターフェース
type DIContainer interface {
	// GetContainer はコンテナインスタンスを返す
	GetContainer() *Container
	// Close は全てのリソースをクローズする
	Close() error
}
