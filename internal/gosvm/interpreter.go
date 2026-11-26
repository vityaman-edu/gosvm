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

type Interpreter struct {
	Block    BlockEntry
	Path     z3.Bool
	Z3       *gosmt.Z3
	Stack    []Frame
	PrintSMT bool
}

type InterpreterArgs struct {
	Entry    *ssa.BasicBlock
	Z3       *gosmt.Z3
	PrintSMT bool
}

type BlockEntry struct {
	Pred *ssa.BasicBlock
	This *ssa.BasicBlock
}

type Frame struct {
	Locals map[ssa.Value]z3.Value
}

func NewInterpreter(args InterpreterArgs) Interpreter {
	return Interpreter{
		Block:    BlockEntry{Pred: nil, This: args.Entry},
		Path:     args.Z3.Bool(true),
		Z3:       args.Z3,
		Stack:    []Frame{{Locals: make(map[ssa.Value]z3.Value)}},
		PrintSMT: args.PrintSMT,
	}
}

func (interp Interpreter) Frame() *Frame {
	return &interp.Stack[len(interp.Stack)-1]
}

func (interp Interpreter) Clone() Interpreter {
	dst := interp

	dst.Stack = slices.Clone(interp.Stack)
	for i := range dst.Stack {
		dst.Stack[i].Locals = maps.Clone(interp.Stack[i].Locals)
	}

	return dst
}

func (interp Interpreter) Interpret() ([]Interpreter, error) {
	block := interp.Block.This
	edgeIndex := slices.Index(block.Preds, interp.Block.Pred)

	last := len(block.Instrs) - 1
	body := block.Instrs[:last]
	terminator := block.Instrs[last]

	interp.logf("SSA: enter %s with path %s", block.String(), interp.Path.String())

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

func (interp Interpreter) interpretPhi(instruction *ssa.Phi, edgeIndex int) error {
	target := instruction
	source := instruction.Edges[edgeIndex]

	frame := interp.Frame()

	value, err := interp.Z3.TranslateValue(source, frame.Locals)
	if err != nil {
		return fmt.Errorf("translate phi operand: %w", err)
	}

	interp.store(target, value)

	return nil
}

//nolint:cyclop,gocyclo,funlen // Each SSA instruction type dispatches to its own handler.
func (interp Interpreter) interpretBody(instruction ssa.Instruction) error {
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

func (interp Interpreter) interpretTerminator(instruction ssa.Instruction) ([]Interpreter, error) {
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

func (interp Interpreter) interpretAlloc(instruction *ssa.Alloc) error {
	return gossa.Errorf("Alloc is not supported", instruction.Pos())
}

func (interp Interpreter) interpretBinOp(instruction *ssa.BinOp) error {
	return interp.interpretValue(instruction)
}

func (interp Interpreter) interpretCall(instruction *ssa.Call) error {
	return gossa.Errorf("Call is not supported", instruction.Pos())
}

func (interp Interpreter) interpretChangeInterface(instruction *ssa.ChangeInterface) error {
	return gossa.Errorf("ChangeInterface is not supported", instruction.Pos())
}

func (interp Interpreter) interpretChangeType(instruction *ssa.ChangeType) error {
	return gossa.Errorf("ChangeType is not supported", instruction.Pos())
}

func (interp Interpreter) interpretConvert(instruction *ssa.Convert) error {
	return gossa.Errorf("Convert is not supported", instruction.Pos())
}

func (interp Interpreter) interpretDebugRef(instruction *ssa.DebugRef) error {
	return gossa.Errorf("DebugRef is not supported", instruction.Pos())
}

func (interp Interpreter) interpretDefer(instruction *ssa.Defer) error {
	return gossa.Errorf("Defer is not supported", instruction.Pos())
}

func (interp Interpreter) interpretExtract(instruction *ssa.Extract) error {
	return gossa.Errorf("Extract is not supported", instruction.Pos())
}

func (interp Interpreter) interpretField(instruction *ssa.Field) error {
	return gossa.Errorf("Field is not supported", instruction.Pos())
}

func (interp Interpreter) interpretFieldAddr(instruction *ssa.FieldAddr) error {
	return gossa.Errorf("FieldAddr is not supported", instruction.Pos())
}

func (interp Interpreter) interpretGo(instruction *ssa.Go) error {
	return gossa.Errorf("Go is not supported", instruction.Pos())
}

func (interp Interpreter) interpretIndex(instruction *ssa.Index) error {
	return gossa.Errorf("Index is not supported", instruction.Pos())
}

func (interp Interpreter) interpretIndexAddr(instruction *ssa.IndexAddr) error {
	return gossa.Errorf("IndexAddr is not supported", instruction.Pos())
}

func (interp Interpreter) interpretLookup(instruction *ssa.Lookup) error {
	return gossa.Errorf("Lookup is not supported", instruction.Pos())
}

func (interp Interpreter) interpretMakeChan(instruction *ssa.MakeChan) error {
	return gossa.Errorf("MakeChan is not supported", instruction.Pos())
}

func (interp Interpreter) interpretMakeClosure(instruction *ssa.MakeClosure) error {
	return gossa.Errorf("MakeClosure is not supported", instruction.Pos())
}

func (interp Interpreter) interpretMakeInterface(instruction *ssa.MakeInterface) error {
	return gossa.Errorf("MakeInterface is not supported", instruction.Pos())
}

func (interp Interpreter) interpretMakeMap(instruction *ssa.MakeMap) error {
	return gossa.Errorf("MakeMap is not supported", instruction.Pos())
}

func (interp Interpreter) interpretMakeSlice(instruction *ssa.MakeSlice) error {
	return gossa.Errorf("MakeSlice is not supported", instruction.Pos())
}

func (interp Interpreter) interpretMapUpdate(instruction *ssa.MapUpdate) error {
	return gossa.Errorf("MapUpdate is not supported", instruction.Pos())
}

func (interp Interpreter) interpretMultiConvert(instruction *ssa.MultiConvert) error {
	return gossa.Errorf("MultiConvert is not supported", instruction.Pos())
}

func (interp Interpreter) interpretNext(instruction *ssa.Next) error {
	return gossa.Errorf("Next is not supported", instruction.Pos())
}

func (interp Interpreter) interpretRange(instruction *ssa.Range) error {
	return gossa.Errorf("Range is not supported", instruction.Pos())
}

func (interp Interpreter) interpretRunDefers(instruction *ssa.RunDefers) error {
	return gossa.Errorf("RunDefers is not supported", instruction.Pos())
}

func (interp Interpreter) interpretSelect(instruction *ssa.Select) error {
	return gossa.Errorf("Select is not supported", instruction.Pos())
}

func (interp Interpreter) interpretSend(instruction *ssa.Send) error {
	return gossa.Errorf("Send is not supported", instruction.Pos())
}

func (interp Interpreter) interpretSlice(instruction *ssa.Slice) error {
	return gossa.Errorf("Slice is not supported", instruction.Pos())
}

func (interp Interpreter) interpretSliceToArrayPointer(instruction *ssa.SliceToArrayPointer) error {
	return gossa.Errorf("SliceToArrayPointer is not supported", instruction.Pos())
}

func (interp Interpreter) interpretStore(instruction *ssa.Store) error {
	return gossa.Errorf("Store is not supported", instruction.Pos())
}

func (interp Interpreter) interpretTypeAssert(instruction *ssa.TypeAssert) error {
	return gossa.Errorf("TypeAssert is not supported", instruction.Pos())
}

func (interp Interpreter) interpretUnOp(instruction *ssa.UnOp) error {
	return interp.interpretValue(instruction)
}

func (interp Interpreter) interpretValue(instruction ssa.Value) error {
	value, err := interp.translateValue(instruction)
	if err != nil {
		return err
	}

	interp.store(instruction, value)

	return nil
}

func (interp Interpreter) translateValue(instruction ssa.Value) (z3.Value, error) {
	locals := interp.Frame().Locals

	value, err := interp.Z3.TranslateValue(instruction, locals)
	if err != nil {
		return nil, err //nolint:wrapcheck // gosmt already provides a source-aware diagnostic.
	}

	return value, nil
}

func (interp Interpreter) store(ssaV ssa.Value, z3V z3.Value) {
	locals := interp.Frame().Locals
	z3V = interp.Z3.Simplify(z3V)
	locals[ssaV] = z3V
	interp.logf("SMT: %s <- %s", ssaV.Name(), z3V.String())
}

func (interp Interpreter) interpretIf(instruction *ssa.If) ([]Interpreter, error) {
	condV, err := interp.translateValue(instruction.Cond)
	if err != nil {
		return nil, err
	}

	cond, okk := condV.(z3.Bool)
	if !okk {
		return nil, gossa.Errorf("expected Bool If condition, got %s", instruction.Pos(), cond.Sort())
	}

	succs := interp.Block.This.Succs

	brThen := interp.Clone()

	brThen.Path, okk = interp.Z3.Simplify(brThen.Path.And(cond)).(z3.Bool)
	if !okk {
		return nil, gossa.Errorf("simplification changed a type", instruction.Pos(), cond.Sort())
	}

	brThen.Block = BlockEntry{Pred: brThen.Block.This, This: succs[0]}

	brElse := interp.Clone()

	brElse.Path, okk = interp.Z3.Simplify(brElse.Path.And(cond.Not())).(z3.Bool)
	if !okk {
		return nil, gossa.Errorf("simplification changed a type", instruction.Pos(), cond.Sort())
	}

	brElse.Block = BlockEntry{Pred: brElse.Block.This, This: succs[1]}

	return []Interpreter{brThen, brElse}, nil
}

func (interp Interpreter) interpretJump(_ *ssa.Jump) ([]Interpreter, error) {
	succs := interp.Block.This.Succs

	then := interp.Clone()
	then.Block = BlockEntry{Pred: then.Block.This, This: succs[0]}

	return []Interpreter{then}, nil
}

func (interp Interpreter) interpretPanic(instruction *ssa.Panic) ([]Interpreter, error) {
	return nil, gossa.Errorf("Panic is not supported", instruction.Pos())
}

func (interp Interpreter) interpretReturn(_ *ssa.Return) ([]Interpreter, error) {
	return []Interpreter{}, nil
}

func (interp Interpreter) logf(format string, a ...any) {
	if !interp.PrintSMT {
		return
	}

	gossa.Logf(interp.Block.This.String()+": "+format, a...)
}
