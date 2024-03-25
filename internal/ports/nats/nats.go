package nats_port

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"testTaskHezzl/internal/good"
	"testTaskHezzl/internal/logger"
	"time"
)

type NatsConn struct {
	Conn *nats.Conn
	Repo logger.LoggerRepository
}

func NewNatsConn(ctx context.Context, repo logger.LoggerRepository) (*NatsConn, error) {
	connection, err := nats.Connect("nats://nats:4222")
	if err != nil {
		return nil, err
	}
	nc := &NatsConn{
		Conn: connection,
		Repo: repo,
	}
	//_, err = nc.Conn.Subscribe("log", func(m *nats.Msg) {
	//
	//})
	ch := make(chan *nats.Msg, 64)
	sub, err := nc.Conn.ChanSubscribe("log", ch)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = sub.Unsubscribe()
		close(ch)
	}()
	for i := 0; i < 5; i++ {
		go func() {
			for msg := range ch {
				log.Println(string(msg.Data))
				var li logger.LogInfo
				if err = json.Unmarshal(msg.Data, &li); err != nil {
					return
				}
				err = nc.Repo.AddLog(ctx, li)
				return
			}
		}()
	}
	return nc, nil
}

func (nc *NatsConn) Log(ctx context.Context, g good.Good, eventTime time.Time) error {
	b, err := json.Marshal(logger.NewLogInfo(g, eventTime))
	if err != nil {
		return err
	}
	if err = nc.Conn.Publish("log", b); err != nil {
		return err
	}
	return nil
}
