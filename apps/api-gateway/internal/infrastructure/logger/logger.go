package logger

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Tags map[string]interface{}

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

func Info(ctx context.Context, msg string, tags map[string]interface{}) {
	log.SetLevel(log.InfoLevel)
	log.WithContext(ctx).WithFields(log.Fields(tags)).Info(msg)
}

func Error(ctx context.Context, err error, msg string, tags map[string]interface{}) {
	if err != nil {
		msg = fmt.Sprintf("%s : %v", msg, err)
	}
	log.SetLevel(log.ErrorLevel)
	log.WithContext(ctx).WithFields(log.Fields(tags)).Error(msg)
}

func Warn(ctx context.Context, msg string, tags map[string]interface{}) {
	log.SetLevel(log.WarnLevel)
	log.WithContext(ctx).WithFields(log.Fields(tags)).Warn(msg)
}

func Debug(ctx context.Context, msg string, tags map[string]interface{}) {
	log.SetLevel(log.DebugLevel)
	log.WithContext(ctx).WithFields(log.Fields(tags)).Debug(msg)
}
