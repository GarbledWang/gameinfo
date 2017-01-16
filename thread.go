package pachong

import (
	"fmt"
)

type GoroutStruct struct{
    Queue chan func() error
    Number int
    Total int
    result chan error
    finishCallback func()
}

func (self *GoroutStruct) Init(number, total int){
    self.Queue = make(chan func() error, total)
    self.Number = number;
    self.Total = total;
    self.result = make(chan error,total)
}

func (self *GoroutStruct) Start(){
    for i := 0; i < self.Number; i++ {
        go func ()  {
            task,ok := <-self.Queue
            if !ok {
                break
            }
            err := task()
            self.result <- err
        }
    }
    
    for j := 0; j < self.Total; j++ {
        res,ok := <-self.result
        if !ok {
            break
        }
        if res != nil {
            fmt.Println(res)
        }
    }

    if self.finishCallback != nil {
        self.finishCallback()
    }
}

func (self *GoroutStruct) Stop(){
    close(self.Queue)
    close(self.result)
}

func (self *GoroutStruct) AddTask(task func() error){
    self.Queue <- task
}

func (self *GoroutStruct) SetFinishCallback(callback func()){
    self.finishCallback = callback
}