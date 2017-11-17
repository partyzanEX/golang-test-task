package pool

// interface of worker function
type Worker func(jobs chan interface{}, results chan interface{})

// Workers pool
type Pool struct {
	jobs   chan interface{}
	results    chan interface{}
	MaxWorkers int
	Workers    []Worker
}

// adding worker to workers array
func (p *Pool) AddWorker(w Worker) {
	p.Workers = append(p.Workers, w)
}

// run goroutines
func (p *Pool) Run() {
	var max int
	if p.MaxWorkers > 0 {
		max = p.MaxWorkers
	} else {
		max = len(p.Workers)
	}

	p.results = make(chan interface{})
	p.jobs = make(chan interface{}, max)

	for _, worker := range p.Workers {
		p.jobs <- 1
		go worker(p.jobs, p.results)
	}
}

// load results from channels
func (p Pool) GetResult() interface{} {
	count := 0
	total := len(p.Workers)

	results := make([]interface{}, total)
	for result := range p.results {
		results[count] = result
		count++

		if count == total {
			close(p.results)
		}
	}

	return results
}

// constructor
func NewPool(maxWorkers int) *Pool {
	return &Pool{
		MaxWorkers: maxWorkers,
	}
}
