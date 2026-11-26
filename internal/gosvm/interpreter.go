package gosvm

import (
	"slices"

	"github.com/ebukreev/go-z3/z3"
	"github.com/vityaman-edu/gosvm/internal/gosmt"
	"github.com/vityaman-edu/gosvm/internal/gossa"
	"golang.org/x/tools/go/ssa"
)

type interpreter struct {
	entry    blockEntry
	stack    []frame
	path     z3.Value
	z3       *gosmt.Z3
	printSMT bool
}

type interpreterArgs struct {
	entry    *ssa.BasicBlock
	z3       *gosmt.Z3
	printSMT bool
}

type blockEntry struct {
	pred *ssa.BasicBlock
	next *ssa.BasicBlock
}

type frame struct {
	locals map[ssa.Value]z3.Value
}

func newInterpreter(args interpreterArgs) interpreter {
	return interpreter{
		entry:    blockEntry{pred: nil, next: args.entry},
		stack:    []frame{{locals: make(map[ssa.Value]z3.Value)}},
		path:     args.z3.Bool(true),
		z3:       args.z3,
		printSMT: args.printSMT,
	}
}

func (interpreter interpreter) frame() *frame {
	return &interpreter.stack[len(interpreter.stack)-1]
}

func (interpreter interpreter) interpret() ([]interpreter, error) {
	block := interpreter.entry.next
	edgeIndex := slices.Index(block.Preds, interpreter.entry.pred)

	last := len(block.Instrs) - 1
	body := block.Instrs[:last]
	terminator := block.Instrs[last]

	for _, instr := range body {
		switch instr := instr.(type) {
		case *ssa.Phi:
			if edgeIndex == -1 {
				return nil, gossa.Errorf("expected edgeIndex for Phi, got none", instr.Pos())
			}

			err := interpreter.interpretPhi(instr, edgeIndex)
			if err != nil {
				return nil, err
			}

		default:
			err := interpreter.interpretBody(instr)
			if err != nil {
				return nil, err
			}
		}
	}

	return interpreter.interpretTerminator(terminator)
}

func (interpreter interpreter) interpretPhi(instruction *ssa.Phi, edgeIndex int) error {
	target := instruction
	source := instruction.Edges[edgeIndex]

	z3 := interpreter.z3
	frame := interpreter.frame()

	value, err := z3.TranslateValue(source, frame.locals)
	if err != nil {
		return err
	}

	frame.locals[target] = value

	return nil
}

func (interpreter interpreter) interpretBody(instruction ssa.Instruction) error {
	switch instruction := instruction.(type) {
	case *ssa.Alloc:
		return interpreter.interpretAlloc(instruction)
	case *ssa.BinOp:
		return interpreter.interpretBinOp(instruction)
	case *ssa.Call:
		return interpreter.interpretCall(instruction)
	case *ssa.ChangeInterface:
		return interpreter.interpretChangeInterface(instruction)
	case *ssa.ChangeType:
		return interpreter.interpretChangeType(instruction)
	case *ssa.Convert:
		return interpreter.interpretConvert(instruction)
	case *ssa.DebugRef:
		return interpreter.interpretDebugRef(instruction)
	case *ssa.Defer:
		return interpreter.interpretDefer(instruction)
	case *ssa.Extract:
		return interpreter.interpretExtract(instruction)
	case *ssa.Field:
		return interpreter.interpretField(instruction)
	case *ssa.FieldAddr:
		return interpreter.interpretFieldAddr(instruction)
	case *ssa.Go:
		return interpreter.interpretGo(instruction)
	case *ssa.Index:
		return interpreter.interpretIndex(instruction)
	case *ssa.IndexAddr:
		return interpreter.interpretIndexAddr(instruction)
	case *ssa.Lookup:
		return interpreter.interpretLookup(instruction)
	case *ssa.MakeChan:
		return interpreter.interpretMakeChan(instruction)
	case *ssa.MakeClosure:
		return interpreter.interpretMakeClosure(instruction)
	case *ssa.MakeInterface:
		return interpreter.interpretMakeInterface(instruction)
	case *ssa.MakeMap:
		return interpreter.interpretMakeMap(instruction)
	case *ssa.MakeSlice:
		return interpreter.interpretMakeSlice(instruction)
	case *ssa.MapUpdate:
		return interpreter.interpretMapUpdate(instruction)
	case *ssa.MultiConvert:
		return interpreter.interpretMultiConvert(instruction)
	case *ssa.Next:
		return interpreter.interpretNext(instruction)
	case *ssa.Range:
		return interpreter.interpretRange(instruction)
	case *ssa.RunDefers:
		return interpreter.interpretRunDefers(instruction)
	case *ssa.Select:
		return interpreter.interpretSelect(instruction)
	case *ssa.Send:
		return interpreter.interpretSend(instruction)
	case *ssa.Slice:
		return interpreter.interpretSlice(instruction)
	case *ssa.SliceToArrayPointer:
		return interpreter.interpretSliceToArrayPointer(instruction)
	case *ssa.Store:
		return interpreter.interpretStore(instruction)
	case *ssa.TypeAssert:
		return interpreter.interpretTypeAssert(instruction)
	case *ssa.UnOp:
		return interpreter.interpretUnOp(instruction)
	default:
		return gossa.Errorf("unsupported SSA instruction type %T", instruction.Pos(), instruction)
	}
}

func (interpreter interpreter) interpretTerminator(instruction ssa.Instruction) ([]interpreter, error) {
	switch instruction := instruction.(type) {
	case *ssa.If:
		return interpreter.interpretIf(instruction)
	case *ssa.Jump:
		return interpreter.interpretJump(instruction)
	case *ssa.Panic:
		return interpreter.interpretPanic(instruction)
	case *ssa.Return:
		return interpreter.interpretReturn(instruction)
	default:
		return nil, gossa.Errorf("unsupported SSA terminator type %T", instruction.Pos(), instruction)
	}
}

func (interpreter interpreter) interpretAlloc(instruction *ssa.Alloc) error {
	return gossa.Errorf("Alloc is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretBinOp(instruction *ssa.BinOp) error {
	return gossa.Errorf("BinOp is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretCall(instruction *ssa.Call) error {
	return gossa.Errorf("Call is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretChangeInterface(instruction *ssa.ChangeInterface) error {
	return gossa.Errorf("ChangeInterface is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretChangeType(instruction *ssa.ChangeType) error {
	return gossa.Errorf("ChangeType is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretConvert(instruction *ssa.Convert) error {
	return gossa.Errorf("Convert is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretDebugRef(instruction *ssa.DebugRef) error {
	return gossa.Errorf("DebugRef is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretDefer(instruction *ssa.Defer) error {
	return gossa.Errorf("Defer is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretExtract(instruction *ssa.Extract) error {
	return gossa.Errorf("Extract is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretField(instruction *ssa.Field) error {
	return gossa.Errorf("Field is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretFieldAddr(instruction *ssa.FieldAddr) error {
	return gossa.Errorf("FieldAddr is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretGo(instruction *ssa.Go) error {
	return gossa.Errorf("Go is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretIndex(instruction *ssa.Index) error {
	return gossa.Errorf("Index is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretIndexAddr(instruction *ssa.IndexAddr) error {
	return gossa.Errorf("IndexAddr is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretLookup(instruction *ssa.Lookup) error {
	return gossa.Errorf("Lookup is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretMakeChan(instruction *ssa.MakeChan) error {
	return gossa.Errorf("MakeChan is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretMakeClosure(instruction *ssa.MakeClosure) error {
	return gossa.Errorf("MakeClosure is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretMakeInterface(instruction *ssa.MakeInterface) error {
	return gossa.Errorf("MakeInterface is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretMakeMap(instruction *ssa.MakeMap) error {
	return gossa.Errorf("MakeMap is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretMakeSlice(instruction *ssa.MakeSlice) error {
	return gossa.Errorf("MakeSlice is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretMapUpdate(instruction *ssa.MapUpdate) error {
	return gossa.Errorf("MapUpdate is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretMultiConvert(instruction *ssa.MultiConvert) error {
	return gossa.Errorf("MultiConvert is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretNext(instruction *ssa.Next) error {
	return gossa.Errorf("Next is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretRange(instruction *ssa.Range) error {
	return gossa.Errorf("Range is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretRunDefers(instruction *ssa.RunDefers) error {
	return gossa.Errorf("RunDefers is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretSelect(instruction *ssa.Select) error {
	return gossa.Errorf("Select is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretSend(instruction *ssa.Send) error {
	return gossa.Errorf("Send is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretSlice(instruction *ssa.Slice) error {
	return gossa.Errorf("Slice is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretSliceToArrayPointer(instruction *ssa.SliceToArrayPointer) error {
	return gossa.Errorf("SliceToArrayPointer is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretStore(instruction *ssa.Store) error {
	return gossa.Errorf("Store is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretTypeAssert(instruction *ssa.TypeAssert) error {
	return gossa.Errorf("TypeAssert is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretUnOp(instruction *ssa.UnOp) error {
	return gossa.Errorf("UnOp is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretIf(instruction *ssa.If) ([]interpreter, error) {
	return nil, gossa.Errorf("If is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretJump(instruction *ssa.Jump) ([]interpreter, error) {
	return nil, gossa.Errorf("Jump is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretPanic(instruction *ssa.Panic) ([]interpreter, error) {
	return nil, gossa.Errorf("Panic is not supported", instruction.Pos())
}

func (interpreter interpreter) interpretReturn(instruction *ssa.Return) ([]interpreter, error) {
	return nil, gossa.Errorf("Return is not supported", instruction.Pos())
}
