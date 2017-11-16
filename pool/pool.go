package pool

import "github.com/partyzanex/golang-test-task/models"

type Worker func(
	jobs <-chan interface{},
	results chan<- interface{},
	errors chan<- interface{},
)

type Channels struct {
	jobs    chan interface{}
	results chan interface{}
	errors  chan interface{}
}

type Pool struct {
	*Channels
	MaxWorkers int
	Workers    []Worker
	Errors     []error
	Results    []models.UrlInfo
}

func (p *Pool) AddWorker(w Worker) {
	p.Workers = append(p.Workers, w)
}

func (p *Pool) Run() {
	p.jobs = make(chan interface{})
	p.results = make(chan interface{})
	p.errors = make(chan interface{})

	for _, worker := range p.Workers {
		go worker(p.jobs, p.results, p.errors)
	}
}

func (p Pool) GetResult() interface{} {
	count := 0
	total := len(p.Workers)

	var errors []interface{}
	//
	//p.Errors = errs
	//close(p.errors)

	e := 0
	results := make([]interface{}, total)

	for error := range p.errors {
		for result := range p.results {
			results = append(results, result)
			count++
			if count == total {
				close(p.results)
			}
		}
		errors = append(errors, error)
		e = len(errors)
	}






	return results
}

func NewPool(maxWorkers int) *Pool{
	return &Pool{
		Channels: &Channels{},
		MaxWorkers: maxWorkers,
	}
}
