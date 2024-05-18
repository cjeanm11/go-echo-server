package pipeline

type Pipeline struct {
	stages []Stage
}

type Stage func(in <-chan any) <-chan any

func NewStage(process func(any) any) Stage {
	return func(in <-chan any) <-chan any {
		out := make(chan any)
		go func() {
			defer close(out)
			for val := range in {
				out <- process(val)
			}
		}()
		return out
	}
}

func (pb *Pipeline) AddStage(process func(any) any) {
	stage := NewStage(process)
	pb.stages = append(pb.stages, stage)
}

func NewPipelineBuilder() *Pipeline {
	return &Pipeline{}
}

func (pb *Pipeline) Build() Stage {
	return func(in <-chan any) <-chan any {
		out := in
		for _, stage := range pb.stages {
			out = stage(out)
		}
		return out
	}
}

func ExecutePipeline(in <-chan any, pipeline Stage) <-chan any {
	return pipeline(in)
}

