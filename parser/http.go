package parser

import (
	"github.com/valyala/fasthttp"
	"github.com/partyzanex/golang-test-task/pool"
	"github.com/partyzanex/golang-test-task/log"
	"fmt"
	"net/http"
	//"encoding/json"
	//"reflect"
	"github.com/partyzanex/golang-test-task/conf"
	"github.com/partyzanex/golang-test-task/models"
)

type Http struct {
	*pool.Pool
	*log.Logger
	Host     string
	Port     string
	Compress bool
}

func (h Http) GetAddr() string {
	return h.Host + ":" + h.Port
}

func (h Http) Serve() error {
	handler := h.HandleRequest
	if h.Compress {
		handler = fasthttp.CompressHandler(handler)
	}

	h.Logger.Write("Listen on " + h.Port)
	return fasthttp.ListenAndServe(h.GetAddr(), handler)
}

func (h *Http) HandleRequest(ctx *fasthttp.RequestCtx) {
	var resp []byte
	statusCode := fasthttp.StatusOK

	urls := models.Urls{}
	err := urls.SetFromBody(ctx.PostBody())
	if err != nil {
		statusCode = fasthttp.StatusBadRequest
		h.Logger.WriteError(err)
	} else {
		h.CreateWorkers(urls)
		h.Run()
		result := h.GetResult()
		fmt.Println(result)
		//
		//resp, err = json.Marshal(result)
		//if err != nil {
		//	statusCode = fasthttp.StatusInternalServerError
		//	h.Logger.WriteError(err)
		//}
	}

	ctx.Response.Header.Set("content-type", "application/json; charset=utf-8")
	ctx.Response.SetStatusCode(statusCode)
	fmt.Fprint(ctx, string(resp))
}

func (h *Http) CreateWorkers(urls models.Urls) {
	h.Pool.Workers = []pool.Worker{}
	for _, url := range urls {
		h.AddWorker(h.CreateWorker(url))
	}
}

func (h *Http) CreateWorker(url string) pool.Worker {
	return func(
		jobs <-chan interface{},
		results chan<- interface{},
		errors chan<- interface{},
	) {
		urlInfo := models.UrlInfo{
			Url: url,
		}

		response, err := http.Get(url)
		if err != nil {
			errors <- err
			results <- urlInfo
		}

		urlInfo.Meta = models.Meta{
			Status:        response.StatusCode,
			ContentType:   response.Header.Get("content-type"),
			ContentLength: response.ContentLength,
		}
		urlInfo.SetElements(response.Body)
		results <- urlInfo
	}
}

func NewHttp(config *conf.Config, logger *log.Logger) *Http {
	return &Http{
		Pool:     pool.NewPool(config.Get("max_workers").Int()),
		Logger:   logger,
		Host:     config.Get("host").String(),
		Port:     config.Get("port").String(),
		Compress: config.Get("compress").Bool(),
	}
}
