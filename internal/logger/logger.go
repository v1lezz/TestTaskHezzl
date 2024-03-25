package logger

import (
	"context"
	"testTaskHezzl/internal/good"
	"time"
)

type LogInfo struct {
	good.Good
	EventTime time.Time
}

func NewLogInfo(g good.Good, eventTime time.Time) LogInfo {
	return LogInfo{
		Good:      g,
		EventTime: eventTime,
	}
}

type Logger interface {
	Log(context.Context, good.Good, time.Time) error
}

type LoggerRepository interface {
	AddLog(context.Context, LogInfo) error
}
