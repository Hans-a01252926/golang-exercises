package lexer

import (
	"patito/token"
)

var keywords = map[string]token.TokenType{
	"programa": token.PROGRAMA,
	"vars":     token.VARS,
	"inicio":   token.INICIO,
	"fin":      token.FIN,
	"entero":   token.ENTERO,
	"flotante": token.FLOTANTE,
	"nula":     token.NULA,
	"si":       token.SI,
	"sino":     token.SINO,
	"mientras": token.MIENTRAS,
	"haz":      token.HAZ,
	"escribe":  token.ESCRIBE,
}

type Lexer struct {
	input        string
	position     int  // apunta al caracter actual
	readPosition int  // apunta a la posición de lectura actual en el input (después del caracter actual)
	ch           byte // caracter actual
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() { // para avanzar un caracter en el input
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.readPosition < len(l.input) && l.input[l.readPosition] == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.IGUAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASIGNA, l.ch)
		}
	case '+':
		tok = newToken(token.MAS, l.ch)
	case '-':
		tok = newToken(token.MENOS, l.ch)
	case '*':
		tok = newToken(token.MULT, l.ch)
	case '/':
		tok = newToken(token.DIVIDE, l.ch)
	case '>':
		tok = newToken(token.MAYOR, l.ch)
	case '<':
		tok = newToken(token.MENOR, l.ch)
	case '!':
		if l.readPosition < len(l.input) && l.input[l.readPosition] == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DIF, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case ';':
		tok = newToken(token.PUNTOCOMA, l.ch)
	case ':':
		tok = newToken(token.DOSPUNTOS, l.ch)
	case ',':
		tok = newToken(token.COMA, l.ch)
	case '(':
		tok = newToken(token.PAR_IZQ, l.ch)
	case ')':
		tok = newToken(token.PAR_DER, l.ch)
	case '{':
		tok = newToken(token.LLAVE_IZQ, l.ch)
	case '}':
		tok = newToken(token.LLAVE_DER, l.ch)
	case '[':
		tok = newToken(token.CORCH_IZQ, l.ch)
	case ']':
		tok = newToken(token.CORCH_DER, l.ch)
	case '"':
		tok.Type = token.LETRERO
		tok.Literal = l.readString()
		return tok // Retorna aquí para evitar avanzar el caracter después de leer el string
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok.Type = lookupIdent(literal)
			tok.Literal = literal
			return tok
		} else if isDigit(l.ch) {
			literal, tokType := l.readNumber()
			tok.Type = tokType
			tok.Literal = literal
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// Auxiliar Numeros
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	var tokType token.TokenType = token.CTE_ENT

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		tokType = token.CTE_FLOT
		l.readChar()

		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position], tokType
}

// Auxiliar Identificadores
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func lookupIdent(ident string) token.TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return token.ID
}

func (l *Lexer) readString() string {
	start := l.position // Fix para error de comillas

	// Avanza al primer carácter dentro del string
	l.readChar()

	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}

	// Si llegó al final del input sin encontrar comilla de cierre
	if l.ch == 0 {
		return l.input[start:l.position]
	}

	// l.position está en la comilla final.
	// Incluimos la comilla final con +1.
	literal := l.input[start : l.position+1]

	// Avanza después de la comilla final
	l.readChar()

	return literal
}

func (l *Lexer) peekChar() byte { // Casos con  ==, !=, que tienen dos caracteres
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
