package maquinavirtual

import "fmt"

type Memory struct {
	values map[int]interface{}
}

func NewMemory() *Memory {
	return &Memory{values: make(map[int]interface{})}
}

func (m *Memory) Set(addr int, value interface{}) {
	m.values[addr] = value
}

func (m *Memory) Get(addr int) (interface{}, error) {
	v, ok := m.values[addr]
	if !ok {
		return nil, fmt.Errorf("dirección no inicializada: %d", addr)
	}
	return v, nil
}
