package lexer

import (
	"testing"
	"token/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){}[],;:*-/> < != ==`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASIGNA, "="},
		{token.MAS, "+"},
		{token.PAR_IZQ, "("},
		{token.PAR_DER, ")"},
		{token.LLAVE_IZQ, "{"},
		{token.LLAVE_DER, "}"},
		{token.CORCH_IZQ, "["},
		{token.CORCH_DER, "]"},
		{token.COMA, ","},
		{token.PUNTOCOMA, ";"},
		{token.DOSPUNTOS, ":"},
		{token.MULT, "*"},
		{token.MENOS, "-"},
		{token.DIVIDE, "/"},
		{token.MAYOR, ">"},
		{token.MENOR, "<"},
		{token.DIF, "!="},
		{token.IGUAL, "=="},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf(
				"tests[%d] - tokentype wrong. expected=%q, got=%q",
				i,
				tt.expectedType,
				tok.Type,
			)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] - literal wrong. expected=%q, got=%q",
				i,
				tt.expectedLiteral,
				tok.Literal,
			)
		}
	}
}
