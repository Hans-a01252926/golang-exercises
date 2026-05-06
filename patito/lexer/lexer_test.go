package lexer

import (
	"patito/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `programa prueba;

vars x, y : entero;
vars z : flotante;

inicio {
	x = 5;
	y = 10;
	z = 3.14;

	escribe("resultado", x + y);

	si (x < y) {
		escribe("x es menor");
	} sino {
		escribe("x no es menor");
	};

	mientras (x != y) haz {
		x = x + 1;
	};
} fin
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PROGRAMA, "programa"},
		{token.ID, "prueba"},
		{token.PUNTOCOMA, ";"},

		{token.VARS, "vars"},
		{token.ID, "x"},
		{token.COMA, ","},
		{token.ID, "y"},
		{token.DOSPUNTOS, ":"},
		{token.ENTERO, "entero"},
		{token.PUNTOCOMA, ";"},

		{token.VARS, "vars"},
		{token.ID, "z"},
		{token.DOSPUNTOS, ":"},
		{token.FLOTANTE, "flotante"},
		{token.PUNTOCOMA, ";"},

		{token.INICIO, "inicio"},
		{token.LLAVE_IZQ, "{"},

		{token.ID, "x"},
		{token.ASIGNA, "="},
		{token.CTE_ENT, "5"},
		{token.PUNTOCOMA, ";"},

		{token.ID, "y"},
		{token.ASIGNA, "="},
		{token.CTE_ENT, "10"},
		{token.PUNTOCOMA, ";"},

		{token.ID, "z"},
		{token.ASIGNA, "="},
		{token.CTE_FLOT, "3.14"},
		{token.PUNTOCOMA, ";"},

		{token.ESCRIBE, "escribe"},
		{token.PAR_IZQ, "("},
		{token.LETRERO, `"resultado"`},
		{token.COMA, ","},
		{token.ID, "x"},
		{token.MAS, "+"},
		{token.ID, "y"},
		{token.PAR_DER, ")"},
		{token.PUNTOCOMA, ";"},

		{token.SI, "si"},
		{token.PAR_IZQ, "("},
		{token.ID, "x"},
		{token.MENOR, "<"},
		{token.ID, "y"},
		{token.PAR_DER, ")"},
		{token.LLAVE_IZQ, "{"},

		{token.ESCRIBE, "escribe"},
		{token.PAR_IZQ, "("},
		{token.LETRERO, `"x es menor"`},
		{token.PAR_DER, ")"},
		{token.PUNTOCOMA, ";"},

		{token.LLAVE_DER, "}"},
		{token.SINO, "sino"},
		{token.LLAVE_IZQ, "{"},

		{token.ESCRIBE, "escribe"},
		{token.PAR_IZQ, "("},
		{token.LETRERO, `"x no es menor"`},
		{token.PAR_DER, ")"},
		{token.PUNTOCOMA, ";"},

		{token.LLAVE_DER, "}"},
		{token.PUNTOCOMA, ";"},

		{token.MIENTRAS, "mientras"},
		{token.PAR_IZQ, "("},
		{token.ID, "x"},
		{token.DIF, "!="},
		{token.ID, "y"},
		{token.PAR_DER, ")"},
		{token.HAZ, "haz"},
		{token.LLAVE_IZQ, "{"},

		{token.ID, "x"},
		{token.ASIGNA, "="},
		{token.ID, "x"},
		{token.MAS, "+"},
		{token.CTE_ENT, "1"},
		{token.PUNTOCOMA, ";"},

		{token.LLAVE_DER, "}"},
		{token.PUNTOCOMA, ";"},

		{token.LLAVE_DER, "}"},
		{token.FIN, "fin"},

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
