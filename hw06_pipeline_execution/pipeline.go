package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecuteStage(in In, done In, stage Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for v := range stage(in) {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for counter := 0; counter < len(stages); counter++ {
		in = ExecuteStage(in, done, stages[counter])
	}
	return in
}
