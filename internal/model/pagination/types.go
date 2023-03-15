package pagination

import (
	"context"
	"github.com/google/uuid"
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	suffixFront   = ":front"
	suffixBack    = ":back"
	defaultExpire = time.Minute * 5
)

var (
	defaultPageSize = int64(10)
)

type PaginationOptions struct {
	Limit     *int64
	Offset    *int64
	Backward  *bool
	LastToken *string
}

func (p *PaginationOptions) EnsureSafe() {
	if p.Backward == nil {
		p.Backward = new(bool)
	}
	if p.Limit == nil {
		p.Limit = &defaultPageSize
	}
}

type CachePaginator struct {
	sorter     any
	sorterType reflect.Type
	cache      cache.Cache
}

func NewCachePaginator(c cache.Cache, sorter any) *CachePaginator {
	t := reflect.TypeOf(sorter)
	for t.Kind() == reflect.Interface || t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return &CachePaginator{
		sorter:     sorter,
		sorterType: t,
		cache:      c,
	}
}
func (s *CachePaginator) GetSorter() any {
	return s.sorter
}

func (s *CachePaginator) LoadSorter(ctx context.Context, key string, backward bool) error {
	if backward {
		key += suffixFront
	} else {
		key += suffixBack
	}
	s.sorter = reflect.New(s.sorterType).Interface()
	err := s.cache.GetCtx(ctx, key, s.sorter)
	if err != nil {
		return err
	}
	return nil
}

func (s *CachePaginator) StoreSorter(ctx context.Context, prefix string, lastToken *string, first, last any) (*string, error) {
	if lastToken == nil {
		lastToken = new(string)
		*lastToken = uuid.New().String()
	}
	front := reflect.New(s.sorterType).Interface()
	err := copier.CopyWithOption(front, first, copier.Option{Converters: []copier.TypeConverter{{
		SrcType: primitive.ObjectID{},
		DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			return src.(primitive.ObjectID).Hex(), nil
		},
	}}})
	if err != nil {
		return nil, err
	}
	// TODO 假如第一次成功，第二次失败会发生什么
	err = s.cache.SetWithExpireCtx(ctx, prefix+*lastToken+suffixFront, front, defaultExpire)
	if err != nil {
		return nil, err
	}

	back := reflect.New(s.sorterType).Interface()
	err = copier.CopyWithOption(back, last, copier.Option{Converters: []copier.TypeConverter{{
		SrcType: primitive.ObjectID{},
		DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			return src.(primitive.ObjectID).Hex(), nil
		},
	}}})
	if err != nil {
		return nil, err
	}
	err = s.cache.SetWithExpireCtx(ctx, prefix+*lastToken+suffixBack, back, defaultExpire)
	if err != nil {
		return nil, err
	}
	return lastToken, nil
}
