package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	PROGRAMA  = "PROGRAMA"
	VARS      = "VARS"
	INICIO    = "INICIO"
	FIN       = "FIN"
	ENTERO    = "ENTERO"
	FLOTANTE  = "FLOTANTE"
	NULA      = "NULA"
	SI        = "SI"
	SINO      = "SINO"
	MIENTRAS  = "MIENTRAS"
	HAZ       = "HAZ"
	ESCRIBE   = "ESCRIBE"
	ID        = "ID"
	CTE_ENT   = "CTE_ENT"
	CTE_FLOT  = "CTE_FLOT"
	LETRERO   = "LETRERO"
	ASIGNA    = "="
	MAS       = "+"
	MENOS     = "-"
	MULTI     = "*"
	DIVIDE    = "/"
	MAYOR     = ">"
	MENOR     = "<"
	DIF       = "!="
	IGUAL     = "=="
	PUNTOCOMA = ";"
	DOSPUNTOS = ":"
	COMA      = ","
	PAR_IZQ   = "("
	PAR_DER   = ")"
	LLAVE_IZQ = "{"
	LLAVE_DER = "}"
	CORCH_IZQ = "["
	CORCH_DER = "]"
)
