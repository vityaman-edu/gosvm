package gosvm

import (
	"fmt"
	"maps"
	"slices"

	"github.com/ebukreev/go-z3/z3"
	"github.com/vityaman-edu/gosvm/internal/gosmt"
	"github.com/vityaman-edu/gosvm/internal/gossa"
	"golang.org/x/tools/go/ssa"
)

type interpreter struct {
	block    blockEntry
	path     z3.Bool
	z3       *gosmt.Z3
	stack    []frame
	printSMT bool
}

type interpreterArgs struct {
	entry    *ssa.BasicBlock
	z3       *gosmt.Z3
	printSMT bool
}

type blockEntry struct {
	pred *ssa.BasicBlock
	this *ssa.BasicBlock
}

type frame struct {
	locals map[ssa.Value]z3.Value
}

func newInterpreter(args interpreterArgs) interpreter {
	return interpreter{
		block:    blockEntry{pred: nil, this: args.entry},
		path:     args.z3.Bool(true),
		z3:       args.z3,
		stack:    []frame{{locals: make(map[ssa.Value]z3.Value)}},
		printSMT: args.printSMT,
	}
}

func (interp interpreter) frame() *frame {
	return &interp.stack[len(interp.stack)-1]
}

func (interp interpreter) clone() interpreter {
	dst := interp

	dst.stack = slices.Clone(interp.stack)
	for i := range dst.stack {
		dst.stack[i].locals = maps.Clone(interp.stack[i].locals)
	}

	return dst
}

func (interp interpreter) interpret() ([]interpreter, error) {
	block := interp.block.this
	edgeIndex := slices.Index(block.Preds, interp.block.pred)

	last := len(block.Instrs) - 1
	body := block.Instrs[:last]
	terminator := block.Instrs[last]

	for _, instr := range body {
		if value, ok := instr.(ssa.Value); ok {
			interp.logf("SSA: %s = %s", value.Name(), instr.String())
		} else {
			interp.logf("SSA: %s", instr.String())
		}

		switch instr := instr.(type) {
		case *ssa.Phi:
			if edgeIndex == -1 {
				return nil, gossa.Errorf("expected edgeIndex for Phi, got none", instr.Pos())
			}

			err := interp.interpretPhi(instr, edgeIndex)
			if err != nil {
				return nil, err
			}

		default:
			err := interp.interpretBody(instr)
			if err != nil {
				return nil, err
			}
		}
	}

	return interp.interpretTerminator(terminator)
}

func (interp interpreter) interpretPhi(instruction *ssa.Phi, edgeIndex int) error {
	target := instruction
	source := instruction.Edges[edgeIndex]

	frame := interp.frame()

	value, err := interp.z3.TranslateValue(source, frame.locals)
	if err != nil {
		return fmt.Errorf("translate phi operand: %w", err)
	}

	frame.locals[target] = value

	return nil
}

//nolint:cyclop,gocyclo,funlen // Each SSA instruction type dispatches to its own handler.
func (interp interpreter) interpretBody(instruction ssa.Instruction) error {
	switch instruction := instruction.(type) {
	case *ssa.Alloc:
		return interp.interpretAlloc(instruction)
	case *ssa.BinOp:
		return interp.interpretBinOp(instruction)
	case *ssa.Call:
		return interp.interpretCall(instruction)
	case *ssa.ChangeInterface:
		return interp.interpretChangeInterface(instruction)
	case *ssa.ChangeType:
		return interp.interpretChangeType(instruction)
	case *ssa.Convert:
		return interp.interpretConvert(instruction)
	case *ssa.DebugRef:
		return interp.interpretDebugRef(instruction)
	case *ssa.Defer:
		return interp.interpretDefer(instruction)
	case *ssa.Extract:
		return interp.interpretExtract(instruction)
	case *ssa.Field:
		return interp.interpretField(instruction)
	case *ssa.FieldAddr:
		return interp.interpretFieldAddr(instruction)
	case *ssa.Go:
		return interp.interpretGo(instruction)
	case *ssa.Index:
		return interp.interpretIndex(instruction)
	case *ssa.IndexAddr:
		return interp.interpretIndexAddr(instruction)
	case *ssa.Lookup:
		return interp.interpretLookup(instruction)
	case *ssa.MakeChan:
		return interp.interpretMakeChan(instruction)
	case *ssa.MakeClosure:
		return interp.interpretMakeClosure(instruction)
	case *ssa.MakeInterface:
		return interp.interpretMakeInterface(instruction)
	case *ssa.MakeMap:
		return interp.interpretMakeMap(instruction)
	case *ssa.MakeSlice:
		return interp.interpretMakeSlice(instruction)
	case *ssa.MapUpdate:
		return interp.interpretMapUpdate(instruction)
	case *ssa.MultiConvert:
		return interp.interpretMultiConvert(instruction)
	case *ssa.Next:
		return interp.interpretNext(instruction)
	case *ssa.Range:
		return interp.interpretRange(instruction)
	case *ssa.RunDefers:
		return interp.interpretRunDefers(instruction)
	case *ssa.Select:
		return interp.interpretSelect(instruction)
	case *ssa.Send:
		return interp.interpretSend(instruction)
	case *ssa.Slice:
		return interp.interpretSlice(instruction)
	case *ssa.SliceToArrayPointer:
		return interp.interpretSliceToArrayPointer(instruction)
	case *ssa.Store:
		return interp.interpretStore(instruction)
	case *ssa.TypeAssert:
		return interp.interpretTypeAssert(instruction)
	case *ssa.UnOp:
		return interp.interpretUnOp(instruction)
	default:
		return gossa.Errorf("unsupported SSA instruction type %T", instruction.Pos(), instruction)
	}
}

func (interp interpreter) interpretTerminator(instruction ssa.Instruction) ([]interpreter, error) {
	interp.logf("SSA: %s", instruction.String())

	switch instruction := instruction.(type) {
	case *ssa.If:
		return interp.interpretIf(instruction)
	case *ssa.Jump:
		return interp.interpretJump(instruction)
	case *ssa.Panic:
		return interp.interpretPanic(instruction)
	case *ssa.Return:
		return interp.interpretReturn(instruction)
	default:
		return nil, gossa.Errorf("unsupported SSA terminator type %T", instruction.Pos(), instruction)
	}
}

func (interp interpreter) interpretAlloc(instruction *ssa.Alloc) error {
	return gossa.Errorf("Alloc is not supported", instruction.Pos())
}

func (interp interpreter) interpretBinOp(instruction *ssa.BinOp) error {
	return interp.interpretValue(instruction)
}

func (interp interpreter) interpretCall(instruction *ssa.Call) error {
	return gossa.Errorf("Call is not supported", instruction.Pos())
}

func (interp interpreter) interpretChangeInterface(instruction *ssa.ChangeInterface) error {
	return gossa.Errorf("ChangeInterface is not supported", instruction.Pos())
}

func (interp interpreter) interpretChangeType(instruction *ssa.ChangeType) error {
	return gossa.Errorf("ChangeType is not supported", instruction.Pos())
}

func (interp interpreter) interpretConvert(instruction *ssa.Convert) error {
	return gossa.Errorf("Convert is not supported", instruction.Pos())
}

func (interp interpreter) interpretDebugRef(instruction *ssa.DebugRef) error {
	return gossa.Errorf("DebugRef is not supported", instruction.Pos())
}

func (interp interpreter) interpretDefer(instruction *ssa.Defer) error {
	return gossa.Errorf("Defer is not supported", instruction.Pos())
}

func (interp interpreter) interpretExtract(instruction *ssa.Extract) error {
	return gossa.Errorf("Extract is not supported", instruction.Pos())
}

func (interp interpreter) interpretField(instruction *ssa.Field) error {
	return gossa.Errorf("Field is not supported", instruction.Pos())
}

func (interp interpreter) interpretFieldAddr(instruction *ssa.FieldAddr) error {
	return gossa.Errorf("FieldAddr is not supported", instruction.Pos())
}

func (interp interpreter) interpretGo(instruction *ssa.Go) error {
	return gossa.Errorf("Go is not supported", instruction.Pos())
}

func (interp interpreter) interpretIndex(instruction *ssa.Index) error {
	return gossa.Errorf("Index is not supported", instruction.Pos())
}

func (interp interpreter) interpretIndexAddr(instruction *ssa.IndexAddr) error {
	return gossa.Errorf("IndexAddr is not supported", instruction.Pos())
}

func (interp interpreter) interpretLookup(instruction *ssa.Lookup) error {
	return gossa.Errorf("Lookup is not supported", instruction.Pos())
}

func (interp interpreter) interpretMakeChan(instruction *ssa.MakeChan) error {
	return gossa.Errorf("MakeChan is not supported", instruction.Pos())
}

func (interp interpreter) interpretMakeClosure(instruction *ssa.MakeClosure) error {
	return gossa.Errorf("MakeClosure is not supported", instruction.Pos())
}

func (interp interpreter) interpretMakeInterface(instruction *ssa.MakeInterface) error {
	return gossa.Errorf("MakeInterface is not supported", instruction.Pos())
}

func (interp interpreter) interpretMakeMap(instruction *ssa.MakeMap) error {
	return gossa.Errorf("MakeMap is not supported", instruction.Pos())
}

func (interp interpreter) interpretMakeSlice(instruction *ssa.MakeSlice) error {
	return gossa.Errorf("MakeSlice is not supported", instruction.Pos())
}

func (interp interpreter) interpretMapUpdate(instruction *ssa.MapUpdate) error {
	return gossa.Errorf("MapUpdate is not supported", instruction.Pos())
}

func (interp interpreter) interpretMultiConvert(instruction *ssa.MultiConvert) error {
	return gossa.Errorf("MultiConvert is not supported", instruction.Pos())
}

func (interp interpreter) interpretNext(instruction *ssa.Next) error {
	return gossa.Errorf("Next is not supported", instruction.Pos())
}

func (interp interpreter) interpretRange(instruction *ssa.Range) error {
	return gossa.Errorf("Range is not supported", instruction.Pos())
}

func (interp interpreter) interpretRunDefers(instruction *ssa.RunDefers) error {
	return gossa.Errorf("RunDefers is not supported", instruction.Pos())
}

func (interp interpreter) interpretSelect(instruction *ssa.Select) error {
	return gossa.Errorf("Select is not supported", instruction.Pos())
}

func (interp interpreter) interpretSend(instruction *ssa.Send) error {
	return gossa.Errorf("Send is not supported", instruction.Pos())
}

func (interp interpreter) interpretSlice(instruction *ssa.Slice) error {
	return gossa.Errorf("Slice is not supported", instruction.Pos())
}

func (interp interpreter) interpretSliceToArrayPointer(instruction *ssa.SliceToArrayPointer) error {
	return gossa.Errorf("SliceToArrayPointer is not supported", instruction.Pos())
}

func (interp interpreter) interpretStore(instruction *ssa.Store) error {
	return gossa.Errorf("Store is not supported", instruction.Pos())
}

func (interp interpreter) interpretTypeAssert(instruction *ssa.TypeAssert) error {
	return gossa.Errorf("TypeAssert is not supported", instruction.Pos())
}

func (interp interpreter) interpretUnOp(instruction *ssa.UnOp) error {
	return interp.interpretValue(instruction)
}

func (interp interpreter) interpretValue(instruction ssa.Value) error {
	value, err := interp.translateValue(instruction)
	if err != nil {
		return err
	}

	locals := interp.frame().locals
	locals[instruction] = value

	interp.logf("SMT: %s <- %s", instruction.Name(), value.String())

	return nil
}

func (interp interpreter) translateValue(instruction ssa.Value) (z3.Value, error) {
	locals := interp.frame().locals

	value, err := interp.z3.TranslateValue(instruction, locals)
	if err != nil {
		return nil, err //nolint:wrapcheck // gosmt already provides a source-aware diagnostic.
	}

	return value, nil
}

func (interp interpreter) interpretIf(instruction *ssa.If) ([]interpreter, error) {
	condV, err := interp.translateValue(instruction.Cond)
	if err != nil {
		return nil, err
	}

	cond, ok := condV.(z3.Bool)
	if !ok {
		return nil, gossa.Errorf("expected Bool If condition, got %s", instruction.Pos(), cond.Sort())
	}

	succs := interp.block.this.Succs

	brThen := interp.clone()
	brThen.path = brThen.path.And(cond)
	brThen.block = blockEntry{pred: brThen.block.this, this: succs[0]}

	brElse := interp.clone()
	brElse.path = brElse.path.And(cond.Not())
	brElse.block = blockEntry{pred: brElse.block.this, this: succs[1]}

	return []interpreter{brThen, brElse}, nil
}

func (interp interpreter) interpretJump(instruction *ssa.Jump) ([]interpreter, error) {
	return nil, gossa.Errorf("Jump is not supported", instruction.Pos())
}

func (interp interpreter) interpretPanic(instruction *ssa.Panic) ([]interpreter, error) {
	return nil, gossa.Errorf("Panic is not supported", instruction.Pos())
}

func (interp interpreter) interpretReturn(_ *ssa.Return) ([]interpreter, error) {
	return []interpreter{}, nil
}

func (interp interpreter) logf(format string, a ...any) {
	if !interp.printSMT {
		return
	}

	gossa.Logf(interp.block.this.String()+": "+format, a...)
}
