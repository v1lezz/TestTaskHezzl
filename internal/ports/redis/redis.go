package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"testTaskHezzl/internal/good"
	"testTaskHezzl/internal/meta"
	"time"
)

type RedisConn struct {
	conn *redis.Client
}

func NewRedisConn(ctx context.Context) *RedisConn {
	client := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "",
		DB:       0,
	})
	return &RedisConn{
		conn: client,
	}
}

func (rc *RedisConn) SaveOnKey(ctx context.Context, g good.Good) error {
	b, err := json.Marshal(g)
	if err != nil {
		return fmt.Errorf("error marshal json: %w", err)
	}
	err = rc.conn.Set(ctx, strconv.Itoa(g.ID), string(b), time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error save cache: %w", err)
	}
	return nil
}

func (rc *RedisConn) GetOnKeyWithLimitAndOffset(ctx context.Context, limit, offset int) (meta.Meta, []good.Good, error) {
	m := meta.Meta{
		Limit:  limit,
		Offset: offset,
	}
	result := make([]good.Good, 0, limit)
	for i := 0; i < limit; i++ {
		res, err := rc.conn.Get(ctx, strconv.Itoa(offset+limit+1)).Result()
		if err != nil {
			return meta.Meta{}, nil, fmt.Errorf("error get from cache: %w", err)
		}
		g := good.Good{}
		err = json.Unmarshal([]byte(res), &g)
		if err != nil {
			return meta.Meta{}, nil, fmt.Errorf("error unmarshal getted good from cache: %w", err)
		}
		result = append(result, g)
		m.Total++
		if g.Removed {
			m.Removed++
		}
	}
	return m, result, nil
}

func (rc *RedisConn) GetOnId(ctx context.Context, id int) (good.Good, error) {
	res, err := rc.conn.Get(ctx, strconv.Itoa(id)).Result()
	if err != nil {
		return good.Good{}, fmt.Errorf("error get from cache: %w", err)
	}
	g := good.Good{}
	err = json.Unmarshal([]byte(res), &g)
	if err != nil {
		return good.Good{}, fmt.Errorf("error unmarshal getted good from cache: %w", err)
	}
	return g, nil
}
