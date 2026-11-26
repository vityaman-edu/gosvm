package gosvm

import (
	"container/list"

	"github.com/vityaman-edu/gosvm/internal/gosmt"
	"github.com/vityaman-edu/gosvm/internal/gossa"
)

type loopInput struct {
	Function gossa.Function
	PrintSMT bool
}

func loop(input loopInput) error {
	const limit = 1000

	start := newInterpreter(interpreterArgs{
		entry:    input.Function.Entry(),
		z3:       gosmt.NewZ3(),
		printSMT: input.PrintSMT,
	})

	queue := list.New()
	queue.PushBack(start)

	for range limit {
		if queue.Len() == 0 {
			break
		}

		element := queue.Front()

		pred, ok := element.Value.(interpreter)
		if !ok {
			panic("interpreter queue contains an unexpected value")
		}

		queue.Remove(element)

		next, err := pred.interpret()
		if err != nil {
			return err
		}

		for _, interpreter := range next {
			queue.PushBack(interpreter)
		}
	}

	return nil
}
