package common

type WorkerPool struct {
	workerlen   int
	JobQueue    chan Job
	WorkerQueue chan chan Job
	PoolRes     chan PoolResult
}

type Job interface {
	Do()
	GetResult() PoolResult
}

type PoolResult struct {
	Status  int
	Msg     string
	JsonRes string
}

func NewWorkerPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,
		JobQueue:    make(chan Job),
		WorkerQueue: make(chan chan Job, workerlen),
		PoolRes:     make(chan PoolResult, workerlen),
	}
}
func (wp *WorkerPool) Run() {
	//初始化worker
	for i := 0; i < wp.workerlen; i++ {
		worker := NewWorker()
		worker.Run(wp.WorkerQueue, wp.PoolRes)
	}
	// 循环获取可用的worker,往worker中写job
	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				worker := <-wp.WorkerQueue
				worker <- job
			}
		}
	}()
}
