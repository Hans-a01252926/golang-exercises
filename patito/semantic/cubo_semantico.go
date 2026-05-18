package semantic

type Operator string

type CubeKey struct {
	Left  Type
	Op    Operator
	Right Type
}

type SemanticCube map[CubeKey]Type

func NewSemanticCube() SemanticCube {
	cube := SemanticCube{}

	// Aritméticos
	cube[CubeKey{TypeEntero, "+", TypeEntero}] = TypeEntero
	cube[CubeKey{TypeEntero, "-", TypeEntero}] = TypeEntero
	cube[CubeKey{TypeEntero, "*", TypeEntero}] = TypeEntero
	cube[CubeKey{TypeEntero, "/", TypeEntero}] = TypeFlotante

	cube[CubeKey{TypeEntero, "+", TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeEntero, "-", TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeEntero, "*", TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeEntero, "/", TypeFlotante}] = TypeFlotante

	cube[CubeKey{TypeFlotante, "+", TypeEntero}] = TypeFlotante
	cube[CubeKey{TypeFlotante, "-", TypeEntero}] = TypeFlotante
	cube[CubeKey{TypeFlotante, "*", TypeEntero}] = TypeFlotante
	cube[CubeKey{TypeFlotante, "/", TypeEntero}] = TypeFlotante

	cube[CubeKey{TypeFlotante, "+", TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeFlotante, "-", TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeFlotante, "*", TypeFlotante}] = TypeFlotante
	cube[CubeKey{TypeFlotante, "/", TypeFlotante}] = TypeFlotante

	// Relacionales
	for _, op := range []Operator{">", "<", "!=", "=="} {
		cube[CubeKey{TypeEntero, op, TypeEntero}] = TypeBool
		cube[CubeKey{TypeEntero, op, TypeFlotante}] = TypeBool
		cube[CubeKey{TypeFlotante, op, TypeEntero}] = TypeBool
		cube[CubeKey{TypeFlotante, op, TypeFlotante}] = TypeBool
	}

	return cube
}

func (c SemanticCube) Result(left Type, op Operator, right Type) Type {
	if result, ok := c[CubeKey{left, op, right}]; ok {
		return result
	}
	return TypeError
}
