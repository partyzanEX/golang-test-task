package parser

import (
	"github.com/valyala/fasthttp"
	"log"
)

type Http struct {
	Host     string
	Port     string
	Compress bool
	Logger   *log.Logger
}

func (h Http) GetAddr() string {
	return h.Host + ":" + h.Port
}

func (h Http) Serve() error {
	handler := h.HandleRequest
	if h.Compress {
		handler = fasthttp.CompressHandler(handler)
	}

	return fasthttp.ListenAndServe(h.GetAddr(), handler)
}

func (h Http) HandleRequest(ctx *fasthttp.RequestCtx) {

}

func createWorker() func() {
	return func() {

	}
}
