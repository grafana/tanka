package export

import (
	"github.com/grafana/tanka/pkg/tanka"
)

func parse(jobs []job, n int) ([]tanka.LoadResult, error) {
	jobCh := make(chan job, len(jobs))
	resCh := make(chan res, len(jobs))

	for w := 0; w <= n; w++ {
		go worker(jobCh, resCh)
	}

	var err error
	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)

	results := make([]tanka.LoadResult, 0, len(jobs))
	for range jobs {
		r := <-resCh
		if r.err != nil {
			err = r.err
		}
		if r.data != nil {
			results = append(results, *r.data)
		}
	}

	return results, err
}

func worker(jobs <-chan job, results chan<- res) {
	for j := range jobs {
		l, err := tanka.Load(j.path, j.opts)
		results <- res{data: l, err: err}
	}
}

type job struct {
	path string
	opts tanka.Opts
}

type res struct {
	data *tanka.LoadResult
	err  error
}
