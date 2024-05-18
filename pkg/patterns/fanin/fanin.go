package fanin

import (
    "sync"
)

type FanIn[T any] struct {
    inputChannels []<-chan T
    mergedChannel chan T
    wg            sync.WaitGroup
    waitOnce      sync.Once
}

func NewFanIn() *FanIn[any] {
    return &FanIn[any]{
        inputChannels: make([]<-chan any, 0),
        mergedChannel: make(chan any),
    }
}

func (f *FanIn[T]) AddInputChannel(ch <-chan T) {
    f.inputChannels = append(f.inputChannels, ch)
}

func (f *FanIn[T]) AddInputChannels(channels...<-chan T) {
    f.inputChannels = append(f.inputChannels, channels...)
}

func (f *FanIn[T]) Start() {
    f.waitOnce.Do(func() {
        f.wg.Add(len(f.inputChannels))
        go func() {
            defer close(f.mergedChannel)
            for _, ch := range f.inputChannels {
                go func(input <-chan T) {
                    defer f.wg.Done()
                    for val := range input {
                        f.mergedChannel <- val
                    }
                }(ch)
            }
            f.wg.Wait()
        }()
    })
}

func (f *FanIn[T]) MergedChannel() <-chan T {
    f.Start()
    return f.mergedChannel
}

func (f *FanIn[T]) Close() {
    f.waitOnce = sync.Once{}
}
