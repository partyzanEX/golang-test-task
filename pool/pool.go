package pool

type Worker func(
	v interface{},
	jobs <-chan interface{},
	results chan<- interface{},
	errors chan<- error,
)

type Pool struct {
	MaxWorkers int
	Workers    []Worker
	jobs       chan interface{}
	results    chan interface{}
	errors     chan error
}

func (p *Pool) AddWorker(w func(v interface{})) {

}

func (p Pool) Run() {

}

func NewPool(maxWorkers int) *Pool {
	return &Pool{
		MaxWorkers: maxWorkers,
		jobs:       make(chan interface{}),
		results:    make(chan interface{}),
		errors:     make(chan error),
	}
}
