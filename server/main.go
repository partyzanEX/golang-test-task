package main

import (
	"github.com/partyzanex/golang-test-task/parser"
	"github.com/partyzanex/golang-test-task/log"
	"github.com/partyzanex/golang-test-task/conf"
)

var (
	config, errConf = conf.NewConfig("./config.json")
	logger, errLog  = log.NewLogger("./task.log", !config.Get("log").Bool())
)

func main() {
	handleError(errLog)
	handleError(errConf)

	s := parser.NewHttp(config, logger)
	s.MaxRedirects = config.Get("max_redirects").Int()

	if err := s.Serve(); err != nil {
		logger.WriteError(err)
		panic(err)
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
