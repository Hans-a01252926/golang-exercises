package quadruples

import (
	"fmt"

	"patito/memoria"
	"patito/semantic"
)

type Generator struct {
	Operators    *Stack[string]
	Operands     *Stack[string]
	Types        *Stack[semantic.Type]
	PendingJumps *Stack[int]
	Quadruples   []Quadruple

	tempCounter int
	Cube        semantic.SemanticCube
	Allocator   *memoria.AddressAllocator
	Errors      []string
}

func NewGenerator(cube semantic.SemanticCube, allocator *memoria.AddressAllocator) *Generator {
	return &Generator{
		Operators:    NewStack[string](),
		Operands:     NewStack[string](),
		Types:        NewStack[semantic.Type](),
		PendingJumps: NewStack[int](),
		Quadruples:   []Quadruple{},
		tempCounter:  0,
		Cube:         cube,
		Allocator:    allocator,
		Errors:       []string{},
	}
}

func (g *Generator) NewTemp(t semantic.Type) string {
	addr, err := g.Allocator.AllocateTemp(string(t))
	if err != nil {
		g.AddError(err.Error())
		return "_"
	}
	return fmt.Sprintf("%d", addr)
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
		g.AddError("No hay suficientes elementos para generar operación")
		return
	}

	resultType := g.Cube.Result(leftType, semantic.Operator(operator), rightType)

	if resultType == semantic.TypeError {
		g.AddError(fmt.Sprintf("Operación inválida: %s %s %s", leftType, operator, rightType))
		return
	}

	temp := g.NewTemp(resultType)
	g.AddQuad(operator, leftOperand, rightOperand, temp)

	g.Operands.Push(temp)
	g.Types.Push(resultType)
}

func (g *Generator) GenerateAssignment(varName string, varType semantic.Type) {
	exprOperand, ok1 := g.Operands.Pop()
	exprType, ok2 := g.Types.Pop()

	if !ok1 || !ok2 {
		g.AddError("No hay expresión para asignar")
		return
	}

	resultType := g.Cube.Result(varType, semantic.OpAsigna, exprType)

	if resultType == semantic.TypeError {
		g.AddError(fmt.Sprintf("Asignación incompatible: No se puede asignar %s a %s", exprType, varType))
		return
	}

	g.AddQuad("=", exprOperand, "_", varName)
}

func (g *Generator) GeneratePrint(value string) {
	g.AddQuad("print", value, "_", "_")
}

func (g *Generator) PeekType() semantic.Type {
	t, ok := g.Types.Peek()
	if !ok {
		return semantic.TypeError
	}
	return t
}

func (g *Generator) PopOperandForPrint() string {
	value, ok1 := g.Operands.Pop()
	_, ok2 := g.Types.Pop()

	if !ok1 || !ok2 {
		g.AddError("No hay valor para imprimir")
		return "_"
	}

	return value
}

func (g *Generator) NextQuad() int {
	return len(g.Quadruples)
}

func (g *Generator) Fill(index int, destination int) {
	if index < 0 || index >= len(g.Quadruples) {
		g.AddError(fmt.Sprintf("Índice de cuádruplo inválido para fill: %d", index))
		return
	}

	g.Quadruples[index].Result = fmt.Sprintf("%d", destination)
}

func (g *Generator) StartIf() {
	expr, ok1 := g.Operands.Pop()
	exprType, ok2 := g.Types.Pop()

	if !ok1 || !ok2 {
		g.AddError("No hay expresión para condición")
		return
	}

	if exprType != semantic.TypeBool {
		g.AddError(fmt.Sprintf("Condición inválida: se esperaba bool y se obtuvo %s", exprType))
		return
	}

	g.AddQuad("GOTOF", expr, "_", "_")
	g.PendingJumps.Push(g.NextQuad() - 1)
}

func (g *Generator) EndIf() {
	end, ok := g.PendingJumps.Pop()
	if !ok {
		g.AddError("No hay salto pendiente para cerrar si")
		return
	}

	g.Fill(end, g.NextQuad())
}

func (g *Generator) ElseIf() {
	// Generar salto para brincar el bloque falso después del bloque verdadero
	g.AddQuad("GOTO", "_", "_", "_")

	falseJump, ok := g.PendingJumps.Pop()
	if !ok {
		g.AddError("No hay salto falso pendiente para sino")
		return
	}

	// Guardar el GOTO pendiente de rellenar al final del sino
	g.PendingJumps.Push(g.NextQuad() - 1)

	// El GOTOF debe brincar al inicio del bloque sino
	g.Fill(falseJump, g.NextQuad())
}

func (g *Generator) EndIfElse() {
	endJump, ok := g.PendingJumps.Pop()
	if !ok {
		g.AddError("No hay salto pendiente para cerrar sino")
		return
	}

	g.Fill(endJump, g.NextQuad())
}

func (g *Generator) StartWhile() {
	g.PendingJumps.Push(g.NextQuad())
}

func (g *Generator) WhileCondition() {
	expr, ok1 := g.Operands.Pop()
	exprType, ok2 := g.Types.Pop()

	if !ok1 || !ok2 {
		g.AddError("No hay expresión para condición de mientras")
		return
	}

	if exprType != semantic.TypeBool {
		g.AddError(fmt.Sprintf("Condición inválida en mientras: se esperaba bool y se obtuvo %s", exprType))
		return
	}

	g.AddQuad("GOTOF", expr, "_", "_")
	g.PendingJumps.Push(g.NextQuad() - 1)
}

func (g *Generator) EndWhile() {
	falseJump, ok1 := g.PendingJumps.Pop()
	returnJump, ok2 := g.PendingJumps.Pop()

	if !ok1 || !ok2 {
		g.AddError("Saltos pendientes insuficientes para cerrar mientras")
		return
	}

	g.AddQuad("GOTO", "_", "_", fmt.Sprintf("%d", returnJump))
	g.Fill(falseJump, g.NextQuad())
}

func (g *Generator) GenerateMainGoto() {
	g.AddQuad("GOTO", "_", "_", "_")
	g.PendingJumps.Push(g.NextQuad() - 1)
}

func (g *Generator) FillMainGoto() {
	jump, ok := g.PendingJumps.Pop()
	if !ok {
		g.AddError("no hay salto pendiente hacia main")
		return
	}
	g.Fill(jump, g.NextQuad())
}

func (g *Generator) GenerateEndFunc() {
	g.AddQuad("ENDFunc", "_", "_", "_")
}
