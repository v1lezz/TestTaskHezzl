package good

import (
	"context"
	"testTaskHezzl/internal/meta"
	"time"
)

type Good struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
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
	SaveOnKey(context.Context, Good) error
	GetOnKeyWithLimitAndOffset(context.Context, int, int) (meta.Meta, []Good, error)
	GetOnId(context.Context, int) (Good, error)
	GoodIsExist(context.Context, int) (bool, error)
}
