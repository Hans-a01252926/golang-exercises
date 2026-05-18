package semantica

type Type string
type Operator string

const (
	TypeEntero   Type = "entero"
	TypeFlotante Type = "flotante"
	TypeBool     Type = "bool"
	TypeNula     Type = "nula"
	TypeError    Type = "error"
)

type CubeKey struct {
	Left  Type
	Op    Operator
	Right Type
}

type SemanticCube map[CubeKey]Type
