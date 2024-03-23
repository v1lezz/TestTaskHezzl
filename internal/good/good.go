package good

import (
	"context"
	"testTaskHezzl/internal/meta"
	"time"
)

type Good struct {
	ID          int
	ProjectID   int
	Name        string
	Description string
	Priority    int
	Removed     bool
	CreatedAt   time.Time
}

type DBRepository interface {
	CreateGood(context.Context, int, string) (Good, error)              // projectId, name
	GetListGoods(context.Context, int, int) (meta.Meta, []Good, error)  //limit, offset
	UpdateGood(context.Context, int, int, string, string) (Good, error) // id, projectId, name, desc
	RemoveGood(context.Context, int, int) (Good, error)                 //id, projectId
	GoodIsExist(context.Context, int, int) (bool, error)                //id, projectId
	ReprioritiizeGood(context.Context, int, int, int) ([]Good, error)   //id, projectId, newPriority
}

type CacheRepository interface {
	SaveOnKey(ctx context.Context, g Good) error
	GetOnKeyWithLimitAndOffset(ctx context.Context, limit, offset int) (meta.Meta, []Good, error)
	GetOnId(ctx context.Context, id int) (Good, error)
}
