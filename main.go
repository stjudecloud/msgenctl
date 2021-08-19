package main

import (
	"github.com/stjudecloud/msgenctl/cmd"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	cmd.Execute()
}
