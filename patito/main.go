package main

import (
	"fmt"
	"os"

	"patito/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run . archivo.patito")
		return
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error leyendo archivo:", err)
		return
	}

	l := parser.NewPatitoLexer(string(content)) // Este es tu adaptador lexer -> parser
	result := parser.Parse(l)                   // Esta es la función Parse generada por goyacc

	if result == 0 && len(l.Errors) == 0 && len(l.Sem.Errors) == 0 {
		fmt.Println("Programa Patito válido")
		return
	}

	fmt.Println("Programa Patito inválido")
	for _, e := range l.Errors {
		fmt.Println("Error:", e)
	}

	for _, e := range l.Sem.Errors {
		fmt.Println("Error semántico:", e)
	}
}
