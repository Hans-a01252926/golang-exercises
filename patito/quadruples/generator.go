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
