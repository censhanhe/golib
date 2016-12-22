/*
* author : cenxinwei@yy.com
* date : 2016-12-17
*/
package pool

import (
    "runtime"
    console "log"
)

type ThreadPool struct{
    workerNum int
    workers map[int]*ThreadWorker   // 这里用map 是为了遍历时的随机化
}

type ThreadWorker struct {
    workerID int
    queue chan func()
    taskQueueSize int
}


func NewThreaedPool( poolSize int, taskQueueSize int) *ThreadPool{
    pool := &ThreadPool{
        workerNum : poolSize,
        workers : make(map[int]*ThreadWorker, poolSize),
    }
    
    for i := 0; i < poolSize; i++{
        worker := &ThreadWorker{}
        worker.init(i,taskQueueSize)
        pool.workers[i] = worker
    }
    return pool
}


func (worker *ThreadWorker)init(id int, taskQueueSize int){
    worker.workerID = id
    worker.taskQueueSize = taskQueueSize
    worker.queue = make(chan func(), taskQueueSize)

    go func(){
        for{
            worker.loop()
        }   
    }()
}

func (worker *ThreadWorker)loop(){
    defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
        if err:=recover();err!=nil{
            const size = 64 << 10
            buf := make([]byte, size)
            buf = buf[:runtime.Stack(buf, false)]            
            // 这里吧\n 替换掉，是位了在syslog里能完整输出
            console.Printf("task worker panic %s, %s", err, buf)
        }
    }() 

    for {
        select {
            case f := <- worker.queue:
                f()
        }
    }
}

// 发起异步调用
func (pool *ThreadPool)AsyncInvokeFunc(f func()) bool {
    // 循环遍历直到找到一个空闲的任务队列
    for _,worker := range pool.workers{
        select {
            case worker.queue <- f:
                return true
            default :
        }
    }
    return false
}