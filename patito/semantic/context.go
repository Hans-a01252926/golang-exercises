package semantic

import "fmt"

type SemanticContext struct {
	DirFunc     *FunctionDirectory
	Cube        SemanticCube
	CurrentFunc string
	CurrentType Type
	Errors      []string
}

func NewSemanticContext() *SemanticContext {
	return &SemanticContext{
		DirFunc:     NewFunctionDirectory(),
		Cube:        NewSemanticCube(),
		CurrentFunc: "global",
		CurrentType: TypeNula,
		Errors:      []string{},
	}
}

func (s *SemanticContext) AddError(msg string) {
	s.Errors = append(s.Errors, msg)
}

func (s *SemanticContext) SetCurrentType(t Type) {
	s.CurrentType = t
}

func (s *SemanticContext) AddVar(name string, varType Type) {
	err := s.DirFunc.AddVarToFunction(s.CurrentFunc, name, varType)
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
