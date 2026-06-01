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

	localIntCount    int
	localFloatCount  int
	localStringCount int

	tempIntCount   int
	tempFloatCount int
	tempBoolCount  int
}

type ParamEntry struct {
	Name    string
	Type    Type
	Address int
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

func (fd *FunctionDirectory) AddVarToFunction(funcName string, varName string, varType Type, address int) error {
	fn, ok := fd.Functions[funcName]
	if !ok {
		return fmt.Errorf("función no encontrada: %s", funcName)
	}

	return fn.Vars.AddVar(varName, varType, funcName, address)
}

func (fd *FunctionDirectory) LookupVar(currentFunc string, varName string) (*VarEntry, bool) {
	if fn, ok := fd.Functions[currentFunc]; ok {
		if v, found := fn.Vars.GetVar(varName); found {
			return v, true
		}
	}

	if global, ok := fd.Functions["global"]; ok {
		if v, found := global.Vars.GetVar(varName); found {
			return v, true
		}
	}

	return nil, false
}

func (fd *FunctionDirectory) AddParam(funcName string, paramName string, paramType Type, address int) error {
	fn, ok := fd.Functions[funcName]
	if !ok {
		return fmt.Errorf("función no encontrada: %s", funcName)
	}

	if _, exists := fn.Vars.GetVar(paramName); exists {
		return fmt.Errorf("parámetro doblemente declarado: %s", paramName)
	}

	param := ParamEntry{
		Name:    paramName,
		Type:    paramType,
		Address: address,
	}

	fn.Params = append(fn.Params, param)

	return fn.Vars.AddVar(paramName, paramType, funcName, address)
}
