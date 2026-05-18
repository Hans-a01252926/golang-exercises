package semantic

import "fmt"

type VarTable struct {
	Vars map[string]*VarEntry
}

type VarEntry struct {
	Name  string
	Type  Type
	Scope string
}

func NewVarTable() *VarTable {
	return &VarTable{
		Vars: make(map[string]*VarEntry),
	}
}

func (vt *VarTable) AddVar(name string, varType Type, scope string) error {
	if _, exists := vt.Vars[name]; exists {
		return fmt.Errorf("variable doblemente declarada: %s", name)
	}

	vt.Vars[name] = &VarEntry{
		Name:  name,
		Type:  varType,
		Scope: scope,
	}

	return nil
}

func (vt *VarTable) GetVar(name string) (*VarEntry, bool) {
	v, ok := vt.Vars[name]
	return v, ok
}
