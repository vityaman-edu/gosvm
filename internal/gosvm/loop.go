package gosvm

import (
	"github.com/vityaman-edu/gosvm/internal/gosched"
	"github.com/vityaman-edu/gosvm/internal/gosmt"
	"github.com/vityaman-edu/gosvm/internal/gossa"
)

type loopInput struct {
	Function gossa.Function
	PrintSMT bool
}

func loop(input loopInput) error {
	const limit = 128

	var sched gosched.Scheduler[*Interpreter]

	sched = gosched.NewEisenhowerScheduler(gosched.NewBfsEisenhower[*Interpreter]())
	sched = gosched.NewLimitedScheduler(limit, sched)

	start := NewInterpreter(InterpreterArgs{
		Entry:    input.Function.Entry(),
		Z3:       gosmt.NewZ3(),
		PrintSMT: input.PrintSMT,
	})

	sched.Submit(&start)

	for {
		pred, ok := sched.Poll()
		if !ok {
			break
		}

		next, err := pred.Interpret()
		if err != nil {
			return err
		}

		for _, interpreter := range next {
			sched.Submit(&interpreter)
		}
	}

	return nil
}
