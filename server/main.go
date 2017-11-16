package main

import (
	"github.com/partyzanex/golang-test-task/parser"
	"github.com/partyzanex/golang-test-task/log"
	"github.com/partyzanex/golang-test-task/conf"
)

var (
	logger, errLog  = log.NewLogger("./task.log", false)
	config, errConf = conf.NewConfig("./config.json")
)

func main() {
	handleError(errLog)
	handleError(errConf)

	s := &parser.Http{
		Host:     config.Get("host").String(),
		Port:     config.Get("port").String(),
		Compress: config.Get("compress").Bool(),
	}

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
