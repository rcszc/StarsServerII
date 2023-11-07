package ServerPool

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"strconv"
	"sync"
)

// GoPool Go 协程池.
type GoPool struct {
	WorkQueue chan func()
	Workers   uint32
	WorksNum  uint32
	WorkWait  sync.WaitGroup
	WorkMutex sync.Mutex
}

// CreateGoPool 创建 Go 协程池.
func CreateGoPool(Workers uint32) *GoPool {
	TempSrc := strconv.FormatUint(uint64(Workers), 10)
	ServerLog.StarLogFmt("INFO", "Create GoPool Workers: "+TempSrc, "Server", "Core", "Pool", "Task")
	return &GoPool{
		WorkQueue: make(chan func(), Workers),
		Workers:   Workers,
	}
}

// Push 添加任务.
func (WorkPool *GoPool) Push(WorkTask func()) {
	WorkPool.WorkQueue <- WorkTask
}

// TasksNumber 当前并行任务数量.
func (WorkPool *GoPool) TasksNumber() (RetNum uint32) {
	WorkPool.WorkMutex.Lock()
	ReturnTasksNumber := WorkPool.WorksNum
	WorkPool.WorkMutex.Unlock()
	return ReturnTasksNumber
}

// Start 启动 Go 协程池.
func (WorkPool *GoPool) Start() {
	ServerLog.StarLogFmt("INFO", "Start GoPool.", "Server", "Core", "Pool", "Task")
	WorkPool.WorkWait.Add(int(WorkPool.Workers))
	for i := uint32(0); i < WorkPool.Workers; i++ {
		go WorkPool.worker()
	}
}

// Close 关闭 Go 协程池
func (WorkPool *GoPool) Close() {
	ServerLog.StarLogFmt("INFO", "Close GoPool.", "Server", "Core", "Pool", "Task")
	close(WorkPool.WorkQueue)
	WorkPool.WorkWait.Wait()
}

func (WorkPool *GoPool) worker() {
	defer WorkPool.WorkWait.Done()
	for TaskFunction := range WorkPool.WorkQueue {
		// 当前并行任务计数.
		WorkPool.WorkMutex.Lock()
		WorkPool.WorksNum++
		WorkPool.WorkMutex.Unlock()

		// execute task function.
		TaskFunction()

		WorkPool.WorkMutex.Lock()
		WorkPool.WorksNum--
		WorkPool.WorkMutex.Unlock()
	}
}
