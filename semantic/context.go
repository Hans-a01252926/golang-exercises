package semantic

import (
	"fmt"
	"patito/memoria"
)

type SemanticContext struct {
	DirFunc     *FunctionDirectory
	Cube        SemanticCube
	CurrentFunc string
	CurrentType Type
	Errors      []string

	Allocator   *memoria.AddressAllocator
	Constants   *memoria.ConstantsTable
	CurrentCall *CallContext

	LastCallReturnType    Type
	LastCallReturnAddress int
}

type CallContext struct {
	FunctionName string
	ArgIndex     int
}

func NewSemanticContext(allocator *memoria.AddressAllocator) *SemanticContext {
	return &SemanticContext{
		DirFunc:               NewFunctionDirectory(),
		Cube:                  NewSemanticCube(),
		CurrentFunc:           "global",
		CurrentType:           TypeNula,
		Errors:                []string{},
		Allocator:             allocator,
		Constants:             memoria.NewConstantsTable(allocator),
		LastCallReturnType:    TypeNula,
		LastCallReturnAddress: -1,
	}
}

func (s *SemanticContext) AddError(msg string) {
	s.Errors = append(s.Errors, msg)
}

func (s *SemanticContext) SetCurrentType(t Type) {
	s.CurrentType = t
}

func (s *SemanticContext) AddVar(name string, varType Type) {
	var addr memoria.Address
	var err error

	if s.CurrentFunc == "global" {
		addr, err = s.Allocator.AllocateGlobal(string(varType))
	} else {
		addr, err = s.Allocator.AllocateLocal(string(varType))
	}

	if err != nil {
		s.AddError(err.Error())
		return
	}

	err = s.DirFunc.AddVarToFunction(s.CurrentFunc, name, varType, int(addr))
	if err != nil {
		s.AddError(err.Error())
	}
}

func (s *SemanticContext) AddVars(names []string, varType Type) {
	for _, name := range names {
		s.AddVar(name, varType)
	}
}

func (s *SemanticContext) StartFunction(name string, returnType Type) {
	err := s.DirFunc.AddFunction(name, returnType)
	if err != nil {
		s.AddError(err.Error())
		return
	}

	if returnType != TypeNula {
		addr, err := s.Allocator.AllocateGlobal(string(returnType))
		if err != nil {
			s.AddError(err.Error())
			return
		}

		s.DirFunc.SetFunctionReturnAddress(name, int(addr))
	}

	s.CurrentFunc = name
}

func (s *SemanticContext) EndFunction() {
	s.CurrentFunc = "global"
}

func (s *SemanticContext) GetVarType(name string) Type {
	v, ok := s.DirFunc.LookupVar(s.CurrentFunc, name)
	if !ok {
		s.AddError(fmt.Sprintf("variable no declarada: %s", name))
		return TypeError
	}

	return v.Type
}

func (s *SemanticContext) CheckOperation(left Type, op Operator, right Type) Type {
	result := s.Cube.Result(left, op, right)

	if result == TypeError {
		s.AddError(fmt.Sprintf("operación inválida: %s %s %s", left, op, right))
	}

	return result
}

func (s *SemanticContext) CheckAssignment(varName string, exprType Type) {
	varType := s.GetVarType(varName)

	if varType == TypeError || exprType == TypeError {
		return
	}

	result := s.Cube.Result(varType, OpAsigna, exprType)

	if result == TypeError {
		s.AddError(fmt.Sprintf("asignación incompatible: no se puede asignar %s a %s", exprType, varType))
	}
}

func (s *SemanticContext) CheckCondition(exprType Type) {
	if exprType != TypeBool {
		s.AddError(fmt.Sprintf("condición inválida: se esperaba bool y se obtuvo %s", exprType))
	}
}

func (s *SemanticContext) AddConstant(literal string, t Type) string {
	addr, err := s.Constants.AddConstant(literal, string(t))
	if err != nil {
		s.AddError(err.Error())
		return "_"
	}
	return fmt.Sprintf("%d", addr)
}

func (s *SemanticContext) GetVarAddressAndType(name string) (string, Type) {
	v, ok := s.DirFunc.LookupVar(s.CurrentFunc, name)
	if !ok {
		s.AddError(fmt.Sprintf("variable no declarada: %s", name))
		return "_", TypeError
	}

	return fmt.Sprintf("%d", v.Address), v.Type
}

func (s *SemanticContext) StartCall(name string) {
	if _, ok := s.DirFunc.GetFunction(name); !ok {
		s.AddError(fmt.Sprintf("función no declarada: %s", name))
		return
	}

	s.CurrentCall = &CallContext{
		FunctionName: name,
		ArgIndex:     0,
	}
}

func (s *SemanticContext) GetCurrentParam() (ParamEntry, bool) {
	if s.CurrentCall == nil {
		s.AddError("no hay llamada activa")
		return ParamEntry{}, false
	}

	fn, ok := s.DirFunc.GetFunction(s.CurrentCall.FunctionName)
	if !ok {
		s.AddError(fmt.Sprintf("función no declarada: %s", s.CurrentCall.FunctionName))
		return ParamEntry{}, false
	}

	if s.CurrentCall.ArgIndex >= len(fn.Params) {
		s.AddError(fmt.Sprintf("demasiados argumentos para función %s", fn.Name))
		return ParamEntry{}, false
	}

	param := fn.Params[s.CurrentCall.ArgIndex]
	s.CurrentCall.ArgIndex++

	return param, true
}

func (s *SemanticContext) EndCall(name string) int {
	fn, ok := s.DirFunc.GetFunction(name)
	if !ok {
		s.LastCallReturnType = TypeError
		s.LastCallReturnAddress = -1
		return -1
	}

	if s.CurrentCall != nil && s.CurrentCall.ArgIndex != len(fn.Params) {
		s.AddError(fmt.Sprintf("cantidad incorrecta de argumentos en llamada a %s", name))
	}

	s.LastCallReturnType = fn.ReturnType
	s.LastCallReturnAddress = fn.ReturnAddress

	s.CurrentCall = nil
	return fn.StartQuad
}

func (s *SemanticContext) SetFunctionStart(name string, startQuad int) {
	err := s.DirFunc.SetFunctionStart(name, startQuad)
	if err != nil {
		s.AddError(err.Error())
	}
}

func (s *SemanticContext) AddParam(name string, paramType Type) {
	var addr memoria.Address
	var err error

	addr, err = s.Allocator.AllocateLocal(string(paramType))
	if err != nil {
		s.AddError(err.Error())
		return
	}

	err = s.DirFunc.AddParam(s.CurrentFunc, name, paramType, int(addr))
	if err != nil {
		s.AddError(err.Error())
	}
}

func (s *SemanticContext) CheckReturn(exprType Type) {
	fn, ok := s.DirFunc.GetFunction(s.CurrentFunc)
	if !ok {
		s.AddError("return fuera de función")
		return
	}

	if fn.ReturnType == TypeNula {
		s.AddError("una función nula no debe regresar valor")
		return
	}

	result := s.Cube.Result(fn.ReturnType, OpAsigna, exprType)
	if result == TypeError {
		s.AddError(fmt.Sprintf("return incompatible: se esperaba %s y se recibió %s", fn.ReturnType, exprType))
	}
}

func (s *SemanticContext) CheckMissingReturn() {
	fn, ok := s.DirFunc.GetFunction(s.CurrentFunc)
	if !ok {
		return
	}

	if fn.ReturnType != TypeNula {
		s.AddError(fmt.Sprintf("función %s debe tener return", fn.Name))
	}
}

func (s *SemanticContext) GetCurrentFunctionReturnAddress() int {
	fn, ok := s.DirFunc.GetFunction(s.CurrentFunc)
	if !ok {
		s.AddError("función actual no encontrada: " + s.CurrentFunc)
		return -1
	}

	return fn.ReturnAddress
}

func (s *SemanticContext) GetLastCallReturnType() Type {
	return s.LastCallReturnType
}

func (s *SemanticContext) GetLastCallReturnAddress() int {
	return s.LastCallReturnAddress
}
