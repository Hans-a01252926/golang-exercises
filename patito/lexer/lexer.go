package lexer

import "token/token"

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

	switch l.ch {
	case '=':
		tok = newToken(token.ASIGNA, l.ch)
	case '+':
		tok = newToken(token.MAS, l.ch)
	case '-':
		tok = newToken(token.MENOS, l.ch)
	case '*':
		tok = newToken(token.MULTI, l.ch)
	case '/':
		tok = newToken(token.DIVIDE, l.ch)
	case '>':
		tok = newToken(token.MAYOR, l.ch)
	case '<':
		tok = newToken(token.MENOR, l.ch)
	case '!':
		if l.peekChar() == '=' {
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
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}
