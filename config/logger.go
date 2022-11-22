package config

import "go.uber.org/zap"

var (
	logger *zap.SugaredLogger
)

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
}
