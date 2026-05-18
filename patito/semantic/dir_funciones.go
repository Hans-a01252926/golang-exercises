package semantic

import "fmt"

type FunctionDirectory struct {
	Functions map[string]*FunctionEntry
}

type FunctionEntry struct {
	Name       string
	ReturnType Type
	Params     []ParamEntry
	Vars       *VarTable
}

type ParamEntry struct {
	Name string
	Type Type
}

func NewFunctionDirectory() *FunctionDirectory {
	fd := &FunctionDirectory{
		Functions: make(map[string]*FunctionEntry),
	}

	fd.Functions["global"] = &FunctionEntry{
		Name:       "global",
		ReturnType: TypeNula,
		Params:     []ParamEntry{},
		Vars:       NewVarTable(),
	}

	return fd
}

func (fd *FunctionDirectory) AddFunction(name string, returnType Type) error {
	if _, exists := fd.Functions[name]; exists {
		return fmt.Errorf("función doblemente declarada: %s", name)
	}

	fd.Functions[name] = &FunctionEntry{
		Name:       name,
		ReturnType: returnType,
		Params:     []ParamEntry{},
		Vars:       NewVarTable(),
	}

	return nil
}

func (fd *FunctionDirectory) GetFunction(name string) (*FunctionEntry, bool) {
	fn, ok := fd.Functions[name]
	return fn, ok
}
