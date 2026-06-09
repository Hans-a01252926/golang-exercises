package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"patito/maquinavirtual"
	"patito/objectfile"
	"patito/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso:")
		fmt.Println("  go run . archivo.patito")
		return
	}

	inputPath := os.Args[1]

	content, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error leyendo archivo:", err)
		return
	}

	l := parser.NewPatitoLexer(string(content))
	result := parser.Parse(l)

	if result != 0 || len(l.Errors) > 0 || len(l.Sem.Errors) > 0 || len(l.Gen.Errors) > 0 {
		fmt.Println("Programa Patito inválido")

		for _, e := range l.Errors {
			fmt.Println("Error sintáctico/léxico:", e)
		}

		for _, e := range l.Sem.Errors {
			fmt.Println("Error semántico:", e)
		}

		for _, e := range l.Gen.Errors {
			fmt.Println("Error de cuádruplos:", e)
		}

		return
	}

	fmt.Println("Programa Patito válido")

	fmt.Println("\nCuádruplos generados:")
	for i, q := range l.Gen.Quadruples {
		fmt.Printf("%d: %s\n", i, q.String())
	}

	objPath := makeObjPath(inputPath)

	// tabla de constantes a obj
	constants := []objectfile.ObjectConstant{}

	for _, c := range l.Sem.Constants.Constants {
		constants = append(constants, objectfile.ObjectConstant{
			Literal: c.Literal,
			Type:    string(c.Type),
			Address: int(c.Address),
		})
	}

	// crear objeto
	obj := objectfile.ObjectFile{
		Quadruples: l.Gen.Quadruples,
		Constants:  constants,
	}

	err = objectfile.Save(objPath, obj)
	if err != nil {
		fmt.Println("Error generando archivo objeto:", err)
		return
	}

	fmt.Println("\nArchivo objeto generado:", objPath)

	// cargar .obj a vm
	loadedObj, err := objectfile.Load(objPath)
	if err != nil {
		fmt.Println("Error cargando archivo objeto:", err)
		return
	}

	// crear VM usando los cuádruplos cargados desde el .obj
	vm := maquinavirtual.NewVirtualMachine(loadedObj.Quadruples)

	// cargar constantes desde el .obj a la memoria constante de la VM
	for _, c := range loadedObj.Constants {
		vm.LoadConstant(c.Address, c.Literal, c.Type)
	}

	fmt.Println("\nEjecutando programa Patito...")

	if err := vm.Run(); err != nil {
		fmt.Println("Error en ejecución:", err)
	}
}

func makeObjPath(inputPath string) string {
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)

	return base + ".obj"
}
