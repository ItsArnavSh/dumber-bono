package main

import "go.uber.org/zap"

func getLogger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	sugar := logger.Sugar()
	sugar.Info("Setting up the logger")
	return sugar
}
