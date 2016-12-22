/*
Author : cenxinwei@163.com
Date : 2016/12/22
*/
package pool

// Pool A Pool is safe for use by multiple goroutines simultaneously.
type Pool struct {
    ctor     func()interface{}
    maxIdle int                 // size of the local array
    freeCh   chan interface{}      
}

//NewPool create a object pool whith constructor and max Idle object size 
func New(ctor func()interface{}, maxIdle int) (p *Pool) {
    p = &Pool{
        ctor : ctor,
        maxIdle : maxIdle,
        freeCh : make(chan interface{}, maxIdle),
    }
    return p
}

// Put adds x to the pool.
func (p *Pool)Put(v interface{}) {
    select {
        case p.freeCh <- v :
        default:
    }
}

// Get selects an arbitrary item from the Pool, removes it from the
// Pool, and returns it to the caller.
func (p *Pool)Get() interface{}{
    select {
        case x := <- p.freeCh:
            return x
        default:
            return p.ctor()
    }
}