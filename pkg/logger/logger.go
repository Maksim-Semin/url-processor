package logger

import (
	"fmt"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var err error

func init() {
	Logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := Logger.Sync(); err != nil {
			fmt.Println(err)
		}
	}()
}
