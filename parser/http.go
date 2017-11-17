package parser

import (
	"fmt"
	"encoding/json"
	"mime"

	"github.com/valyala/fasthttp"

	"github.com/partyzanex/golang-test-task/pool"
	"github.com/partyzanex/golang-test-task/log"
	"github.com/partyzanex/golang-test-task/conf"
	"github.com/partyzanex/golang-test-task/models"
)

// Handler
type Http struct {
	*pool.Pool
	*log.Logger
	*fasthttp.Client
	Host         string
	Port         string
	Compress     bool
	MaxRedirects int
}

func (h Http) GetAddr() string {
	return h.Host + ":" + h.Port
}

// run http-server
func (h Http) Serve() error {
	handler := h.HandleRequest
	if h.Compress {
		handler = fasthttp.CompressHandler(handler)
	}

	h.Logger.Write("Listen on " + h.Port)
	return fasthttp.ListenAndServe(h.GetAddr(), handler)
}

// request hendler
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

		resp, err = json.Marshal(result)
		if err != nil {
			statusCode = fasthttp.StatusInternalServerError
			h.Logger.WriteError(err)
		}
	}

	ctx.Response.Header.Set("content-type", "application/json; charset=utf-8")
	ctx.Response.SetStatusCode(statusCode)
	fmt.Fprint(ctx, string(resp))
}

// adding workers in Pool
func (h *Http) CreateWorkers(urls models.Urls) {
	h.Pool.Workers = []pool.Worker{}
	for _, url := range urls {
		h.Pool.AddWorker(h.CreateWorker(url))
	}
}

// creating function-worker
func (h *Http) CreateWorker(url string) pool.Worker {
	return func(results chan interface{}, next func()) {
		urlInfo := models.UrlInfo{
			Url: url,
		}

		response, err := h.DoRequest(url)
		if err != nil {
			urlInfo.SetError(err)

			results <- urlInfo
			next()
			return
		}

		mimeType, _, err := mime.ParseMediaType(string(response.Header.ContentType()))
		if err != nil {
			urlInfo.SetError(err)
		}

		urlInfo.Meta = models.Meta{
			Status:        response.StatusCode(),
			ContentType:   mimeType,
			ContentLength: response.Header.ContentLength(),
		}
		urlInfo.SetElements(response.Body())

		results <- urlInfo
		next()
	}
}

// open url and return response
func (h Http) DoRequest(url string) (*fasthttp.Response, error) {
	response := fasthttp.AcquireResponse()

	request := fasthttp.AcquireRequest()
	request.SetRequestURI(url)

	err := h.Client.Do(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// constructor
func NewHttp(config *conf.Config, logger *log.Logger) *Http {
	return &Http{
		Pool:     pool.NewPool(config.Get("max_workers").Int()),
		Logger:   logger,
		Client:   &fasthttp.Client{},
		Host:     config.Get("host").String(),
		Port:     config.Get("port").String(),
		Compress: config.Get("compress").Bool(),
	}
}
