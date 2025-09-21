package repository

import (
	"context"

	"github.com/google/uuid"
)

// Pagination ページネーションのパラメータ
type Pagination struct {
	Page     int `json:"page"`      // ページ番号（1から開始）
	PageSize int `json:"page_size"` // 1ページあたりのアイテム数
	Offset   int `json:"offset"`    // オフセット（自動計算）
}

// NewPagination 新しいPaginationインスタンスを作成
func NewPagination(page, pageSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // デフォルトページサイズ
	}
	if pageSize > 100 {
		pageSize = 100 // 最大ページサイズ
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

// Sort ソートのパラメータ
type Sort struct {
	Field string `json:"field"` // ソートフィールド
	Order string `json:"order"` // ソート順（asc/desc）
}

// NewSort 新しいSortインスタンスを作成
func NewSort(field, order string) *Sort {
	if order != "asc" && order != "desc" {
		order = "asc" // デフォルトは昇順
	}

	return &Sort{
		Field: field,
		Order: order,
	}
}

// IsAscending 昇順かどうかを判定
func (s Sort) IsAscending() bool {
	return s.Order == "asc"
}

// IsDescending 降順かどうかを判定
func (s Sort) IsDescending() bool {
	return s.Order == "desc"
}

// PagedResult ページング結果を表現する構造体
type PagedResult[T any] struct {
	Items      []T   `json:"items"`       // アイテムのリスト
	TotalCount int64 `json:"total_count"` // 全アイテム数
	Page       int   `json:"page"`        // 現在のページ番号
	PageSize   int   `json:"page_size"`   // ページサイズ
	TotalPages int   `json:"total_pages"` // 総ページ数
	HasNext    bool  `json:"has_next"`    // 次のページがあるか
	HasPrev    bool  `json:"has_prev"`    // 前のページがあるか
}

// NewPagedResult 新しいPagedResultインスタンスを作成
func NewPagedResult[T any](items []T, totalCount int64, pagination *Pagination) *PagedResult[T] {
	totalPages := int((totalCount + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))

	return &PagedResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
		HasNext:    pagination.Page < totalPages,
		HasPrev:    pagination.Page > 1,
	}
}

// BaseRepository 全てのリポジトリが継承する基本インターフェース
type BaseRepository[T any] interface {
	// Save エンティティを保存（新規作成または更新）
	Save(ctx context.Context, entity T) error

	// FindByID IDでエンティティを検索
	FindByID(ctx context.Context, id uuid.UUID) (T, error)

	// Delete エンティティを削除
	Delete(ctx context.Context, id uuid.UUID) error

	// Exists エンティティが存在するかチェック
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

// UserScopedRepository ユーザースコープ付きのリポジトリインターフェース
// ユーザーに属するエンティティ用（Account, Transaction, Category, Budget等）
type UserScopedRepository[T any] interface {
	BaseRepository[T]

	// FindByUserID ユーザーIDで関連エンティティを検索
	FindByUserID(ctx context.Context, userID uuid.UUID, pagination *Pagination, sorts ...*Sort) (*PagedResult[T], error)

	// FindByUserIDAndID ユーザーIDとエンティティIDで検索
	FindByUserIDAndID(ctx context.Context, userID, entityID uuid.UUID) (T, error)

	// DeleteByUserIDAndID ユーザーIDとエンティティIDで削除
	DeleteByUserIDAndID(ctx context.Context, userID, entityID uuid.UUID) error

	// CountByUserID ユーザーのエンティティ数をカウント
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

// Transactional トランザクション処理のインターフェース
type Transactional interface {
	// WithTx トランザクション内で処理を実行
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// SearchFilter 検索フィルターの基本インターフェース
type SearchFilter interface {
	// ToMap フィルター条件をマップ形式で取得
	ToMap() map[string]any

	// IsEmpty フィルターが空かどうかを判定
	IsEmpty() bool
}

// DateRangeFilter 日付範囲フィルター
type DateRangeFilter struct {
	StartDate *string `json:"start_date,omitempty"` // YYYY-MM-DD形式
	EndDate   *string `json:"end_date,omitempty"`   // YYYY-MM-DD形式
}

// ToMap フィルター条件をマップ形式で取得
func (f DateRangeFilter) ToMap() map[string]any {
	result := make(map[string]any)
	if f.StartDate != nil {
		result["start_date"] = *f.StartDate
	}
	if f.EndDate != nil {
		result["end_date"] = *f.EndDate
	}
	return result
}

// IsEmpty フィルターが空かどうかを判定
func (f DateRangeFilter) IsEmpty() bool {
	return f.StartDate == nil && f.EndDate == nil
}
