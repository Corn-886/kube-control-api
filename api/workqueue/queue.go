package workqueue

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/workqueue"
	"kube-control-api/api/common/set"
	"kube-control-api/api/common/wait"
	"sync"
	"time"
)

var Evaluator = NewQuotaEvaluator()

type quotaEvaluator struct {
	queue      *workqueue.Type
	work       map[string][]*Task
	dirtyWork  map[string][]*Task
	workLock   sync.Mutex
	workers    int
	init       sync.Once
	stopCh     chan struct{}
	inProgress set.String
}

func (e *quotaEvaluator) Stop() {
	close(e.stopCh)
}

func NewQuotaEvaluator() *quotaEvaluator {
	e := &quotaEvaluator{
		queue:      workqueue.New(),
		workers:    1, // 这里可以根据实际需求设置
		work:       map[string][]*Task{},
		dirtyWork:  map[string][]*Task{},
		inProgress: set.String{},
		stopCh:     make(chan struct{}),
	}
	e.init.Do(func() {
		e.run()
	}) // 初始化workqueue
	return e
}

// 持续运行work queue, 监听任务
func (e *quotaEvaluator) run() {
	logrus.Info("启动配额评估器 %d workers", e.workers)
	for i := 0; i < e.workers; i++ {
		go wait.Until(e.doWork, time.Second, e.stopCh)
	}
	logrus.Info("启动配额评估器完成")

}

func (e *quotaEvaluator) doWork() {

	logrus.Info("do work 开始")
	workFunc := func() bool {
		ts, work := e.getWork()
		defer e.completeWork(ts)

		if len(work) == 0 {
			return false
		}
		// 执行任务
		for _, task := range work {
			task.Execute(task)
		}
		return false
	}

	for {
		if quit := workFunc(); quit {
			logrus.Info("quota evaluator worker shutdown")
			return
		}
	}
}

// 添加任务到workqueue, 任务会被分配到worker执行
func (e *quotaEvaluator) AddWork(task *Task) {
	logrus.Info("添加任务 %s to 工作队列", task.Name)
	e.workLock.Lock()
	defer e.workLock.Unlock()

	e.queue.Add(task)

	if e.inProgress.Has(task.Name) { // 如果任务正在执行中，将任务放入dirtyWork
		e.dirtyWork[task.Name] = append(e.dirtyWork[task.Name], task)
		return
	}

	e.work[task.Name] = append(e.work[task.Name], task)
	logrus.Info("队列长度", len(e.work))
}

func (e *quotaEvaluator) getWork() (string, []*Task) {
	logrus.Info("开始获取任务")
	uncastTask, _ := e.queue.Get() // 从workqueue中获取任务，如果不存在，则等待
	//uncastTask 强制转为Task
	ts := uncastTask.(*Task).Name

	e.workLock.Lock()
	defer e.workLock.Unlock()

	work := e.work[ts]
	delete(e.work, ts)
	delete(e.dirtyWork, ts)
	e.inProgress.Insert(ts)

	return ts, work
}

func (e *quotaEvaluator) completeWork(ts string) {
	logrus.Info("完成任务", ts)
	e.workLock.Lock()
	defer e.workLock.Unlock()

	e.queue.Done(ts)
	e.work[ts] = e.dirtyWork[ts]
	delete(e.dirtyWork, ts)
	e.inProgress.Delete(ts)
}
