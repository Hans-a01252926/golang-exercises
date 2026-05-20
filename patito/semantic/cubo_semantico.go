package semantic

type Operator string

const (
	OpSuma   Operator = "+"
	OpResta  Operator = "-"
	OpMult   Operator = "*"
	OpDiv    Operator = "/"
	OpMayor  Operator = ">"
	OpMenor  Operator = "<"
	OpDif    Operator = "!="
	OpIgual  Operator = "=="
	OpAsigna Operator = "="
)

type CubeKey struct {
	Left  Type
	Op    Operator
	Right Type
}

type SemanticCube map[CubeKey]Type

func NewSemanticCube() SemanticCube {
	cube := SemanticCube{}

	// Aritméticos
	cube[CubeKey{TypeEntero, OpSuma, TypeEntero}] = TypeEntero
	cube[CubeKey{TypeEntero, OpResta, TypeEntero}] = TypeEntero
	cube[CubeKey{TypeEntero, OpMult, TypeEntero}] = TypeEntero
	cube[CubeKey{TypeEntero, OpDiv, TypeEntero}] = TypeFlotante

	cube[CubeKey{TypeEntero, OpSuma, TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeEntero, OpResta, TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeEntero, OpMult, TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeEntero, OpDiv, TypeFlotante}] = TypeFlotante

	cube[CubeKey{TypeFlotante, OpSuma, TypeEntero}] = TypeFlotante
	cube[CubeKey{TypeFlotante, OpResta, TypeEntero}] = TypeFlotante
	cube[CubeKey{TypeFlotante, OpMult, TypeEntero}] = TypeFlotante
	cube[CubeKey{TypeFlotante, OpDiv, TypeEntero}] = TypeFlotante

	cube[CubeKey{TypeFlotante, OpSuma, TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeFlotante, OpResta, TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeFlotante, OpMult, TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeFlotante, OpDiv, TypeFlotante}] = TypeFlotante

	// Relacionales
	for _, op := range []Operator{OpMayor, OpMenor, OpDif, OpIgual} {
		cube[CubeKey{TypeEntero, op, TypeEntero}] = TypeBool
		cube[CubeKey{TypeEntero, op, TypeFlotante}] = TypeBool
		cube[CubeKey{TypeFlotante, op, TypeEntero}] = TypeBool
		cube[CubeKey{TypeFlotante, op, TypeFlotante}] = TypeBool
	}

	// Asignación
	cube[CubeKey{TypeEntero, OpAsigna, TypeEntero}] = TypeEntero
	cube[CubeKey{TypeFlotante, OpAsigna, TypeEntero}] = TypeFlotante
	cube[CubeKey{TypeFlotante, OpAsigna, TypeFlotante}] = TypeFlotante

	return cube
}

func (c SemanticCube) Result(left Type, op Operator, right Type) Type {
	if result, ok := c[CubeKey{left, op, right}]; ok {
		return result
	}
	return TypeError
}
