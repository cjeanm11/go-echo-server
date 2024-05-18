package future

type TaskCommand[T any] interface {
	execute() T
}

type Future[T any] struct {
	execute func() T
	result  chan T
}

func NewFuture[T any](execute func() T) Future[T] {
	return Future[T]{
		execute: execute,
		result:  make(chan T, 1), // Initialize as buffered channel
	}
}

func (f *Future[T]) Submit() {
	go func() {
		result := f.execute()
		f.result <- result
	}()
}

func (f *Future[T]) Get() T {
	return <-f.result
}
