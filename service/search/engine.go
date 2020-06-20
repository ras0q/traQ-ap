package search

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traQ/utils/optional"
)

// ErrServiceUnavailable エラー 現在検索サービスが利用できません
var ErrServiceUnavailable = errors.New("search service is unavailable")

// Engine 検索エンジンインターフェイス
type Engine interface {
	// Do 与えられたクエリで検索を実行します
	Do(q *Query) (Result, error)
	// Available 検索サービスが利用可能かどうかを返します
	Available() bool
	// Close 検索サービスを終了します
	Close() error
}

// Query 検索クエリ TODO
type Query struct {
	// Word 検索ワード (仮置き)
	Word        optional.String `query:"word"`   // 検索ワード 空白区切り(複数)をうまく扱ってくれる
	After       optional.Time   `query:"after"`  // 以降(投稿日時)
	Before      optional.Time   `query:"before"` // 以前(投稿日時)
	To          optional.UUID   `query:"to"`     // メンション先
	From        optional.UUID   `query:"from"`   // 投稿者
	Cite        optional.UUID   `query:"cite"`   // 引用しているメッセージ
	IsEdited    optional.Bool   `query:"isEdited"`
	IsCited     optional.Bool   `query:"isCited"`
	IsPinned    optional.Bool   `query:"isPinned"`
	HasURL      optional.Bool   `query:"hasURL"`
	HasEmbedded optional.Bool   `query:"hasEmbedded"`
	HasImage    optional.Bool   `query:"hasImage"`
	HasMovie    optional.Bool   `query:"hasMovie"`
	HasAudio    optional.Bool   `query:"hasAudio"`
}

func (q Query) Validate() error {
	//return vd.ValidateStruct(&q,
	//	vd.Field(&q.Word, validator.SearchWordRule...),
	return nil
}

// Result 検索結果インターフェイス TODO
type Result interface {
	// Get 仮置き
	Get() map[uuid.UUID]string
	// GetMessages() (ms []*model.Message, more bool)
}
