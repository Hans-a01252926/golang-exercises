package quadruples

import (
	"fmt"

	"patito/semantic"
)

type Generator struct {
	Operators  *Stack[string]
	Operands   *Stack[string]
	Types      *Stack[semantic.Type]
	Quadruples []Quadruple

	tempCounter int
	Cube        semantic.SemanticCube
	Errors      []string
}

func NewGenerator(cube semantic.SemanticCube) *Generator {
	return &Generator{
		Operators:   NewStack[string](),
		Operands:    NewStack[string](),
		Types:       NewStack[semantic.Type](),
		Quadruples:  []Quadruple{},
		tempCounter: 0,
		Cube:        cube,
		Errors:      []string{},
	}
}

func (g *Generator) NewTemp() string {
	g.tempCounter++
	return fmt.Sprintf("t%d", g.tempCounter)
}

func (g *Generator) AddError(msg string) {
	g.Errors = append(g.Errors, msg)
}

func (g *Generator) AddQuad(op string, left string, right string, result string) {
	g.Quadruples = append(g.Quadruples, Quadruple{
		Operator: op,
		Left:     left,
		Right:    right,
		Result:   result,
	})
}

func (g *Generator) PushOperand(operand string, typ semantic.Type) {
	g.Operands.Push(operand)
	g.Types.Push(typ)
}

func (g *Generator) PushOperator(op string) {
	g.Operators.Push(op)
}

func (g *Generator) GenerateBinaryOperation() {
	rightOperand, ok1 := g.Operands.Pop()
	rightType, ok2 := g.Types.Pop()
	leftOperand, ok3 := g.Operands.Pop()
	leftType, ok4 := g.Types.Pop()
	operator, ok5 := g.Operators.Pop()

	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		g.AddError("no hay suficientes elementos para generar operación")
		return
	}

	resultType := g.Cube.Result(leftType, semantic.Operator(operator), rightType)

	if resultType == semantic.TypeError {
		g.AddError(fmt.Sprintf("operación inválida: %s %s %s", leftType, operator, rightType))
		return
	}

	temp := g.NewTemp()
	g.AddQuad(operator, leftOperand, rightOperand, temp)

	g.Operands.Push(temp)
	g.Types.Push(resultType)
}

func (g *Generator) GenerateAssignment(varName string, varType semantic.Type) {
	exprOperand, ok1 := g.Operands.Pop()
	exprType, ok2 := g.Types.Pop()

	if !ok1 || !ok2 {
		g.AddError("no hay expresión para asignar")
		return
	}

	resultType := g.Cube.Result(varType, semantic.OpAsigna, exprType)

	if resultType == semantic.TypeError {
		g.AddError(fmt.Sprintf("asignación incompatible: no se puede asignar %s a %s", exprType, varType))
		return
	}

	g.AddQuad("=", exprOperand, "_", varName)
}

func (g *Generator) GeneratePrint(value string) {
	g.AddQuad("print", value, "_", "_")
}
