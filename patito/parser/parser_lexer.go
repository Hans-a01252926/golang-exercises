package parser

/* Adaptador para usar el lexer con el parser generado por goyacc. El parser espera un
método Lex que retorne un token y un método Error para manejar errores de análisis.*/

import (
	"fmt"
	"patito/lexer"
	"patito/token"
)

type PatitoLexer struct {
	l      *lexer.Lexer
	Errors []string
}

func NewPatitoLexer(input string) *PatitoLexer {
	return &PatitoLexer{
		l: lexer.New(input),
	}
}

func (p *PatitoLexer) Lex(lval *yySymType) int {
	tok := p.l.NextToken()

	switch tok.Type {
	case token.EOF:
		return 0

	case token.ID:
		lval.lit = tok.Literal
		return ID

	case token.CTE_ENT:
		lval.lit = tok.Literal
		return CTE_ENT

	case token.CTE_FLOT:
		lval.lit = tok.Literal
		return CTE_FLOT

	case token.LETRERO:
		lval.lit = tok.Literal
		return LETRERO

	case token.PROGRAMA:
		return PROGRAMA
	case token.VARS:
		return VARS
	case token.INICIO:
		return INICIO
	case token.FIN:
		return FIN
	case token.ENTERO:
		return ENTERO
	case token.FLOTANTE:
		return FLOTANTE
	case token.NULA:
		return NULA
	case token.SI:
		return SI
	case token.SINO:
		return SINO
	case token.MIENTRAS:
		return MIENTRAS
	case token.HAZ:
		return HAZ
	case token.ESCRIBE:
		return ESCRIBE

	case token.ASIGNA:
		return ASIGNA
	case token.MAS:
		return MAS
	case token.MENOS:
		return MENOS
	case token.MULT:
		return MULT
	case token.DIVIDE:
		return DIVIDE
	case token.MAYOR:
		return MAYOR
	case token.MENOR:
		return MENOR
	case token.DIF:
		return DIF
	case token.IGUAL:
		return IGUAL

	case token.PUNTOCOMA:
		return PUNTOCOMA
	case token.DOSPUNTOS:
		return DOSPUNTOS
	case token.COMA:
		return COMA
	case token.PAR_IZQ:
		return PAR_IZQ
	case token.PAR_DER:
		return PAR_DER
	case token.LLAVE_IZQ:
		return LLAVE_IZQ
	case token.LLAVE_DER:
		return LLAVE_DER
	}

	p.Errors = append(p.Errors, fmt.Sprintf("token ilegal: %s", tok.Literal))
	return 0
}

func (p *PatitoLexer) Error(e string) {
	p.Errors = append(p.Errors, e)
}
