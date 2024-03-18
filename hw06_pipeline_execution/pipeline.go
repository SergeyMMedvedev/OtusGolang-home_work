package hw06pipelineexecution

import (
	"fmt"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func (s Stage) String() string {
	return fmt.Sprintf("%T", s)
}

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

// func ExecutePipeline(in In, done In, stages ...Stage) Out {
// 	// Place your code here.
// 	switch len(stages) {
// 	case 0:
// 		return in
// 	case 1:
// 		return ExecuteStage(in, done, stages[0])
// 	default:
// 		return ExecutePipeline(ExecuteStage(in, done, stages[0]), done, stages[1:]...)
// 	}
// }

//	func ExecutePipeline(in In, done In, stages ...Stage) Out {
//		// Place your code here.
//		for counter := 0; counter < len(stages); counter++ {
//			in = ExecuteStage(in, done, stages[counter])
//		}
//		return in
//	}
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	return stages[3](stages[2](stages[1](stages[0](in))))
}
