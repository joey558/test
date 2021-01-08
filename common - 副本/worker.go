package common

type Worker struct {
	JobQueue chan Job
}

func NewWorker() Worker {
	return Worker{JobQueue: make(chan Job)}
}
func (w Worker) Run(wq chan chan Job, ps chan PoolResult) {
	go func() {
		for {
			wq <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				job.Do()
				ps <- job.GetResult()
			}
		}
	}()
}
