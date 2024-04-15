package resource

import (
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	lg "main/pkg/logger"
)

type Config struct {
	DB         string
	DBUser     string
	DBPassword string
	DBName     string
}

var CFG Config

func init() {
	_, err := toml.DecodeFile("resource/config.toml", &CFG)
	if err != nil {
		lg.Logger.Error("error on decoding config", zap.Error(err))
	}
}
