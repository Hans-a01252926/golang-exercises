package maquinavirtual

import "fmt"

type ExecutionMemory struct {
	GlobalMemory *Memory
	TempMemory   *Memory
	ConstMemory  *Memory

	CallStack  []*ActivationRecord
	PendingERA *ActivationRecord
}

func NewExecutionMemory() *ExecutionMemory {
	return &ExecutionMemory{
		GlobalMemory: NewMemory(),
		TempMemory:   NewMemory(),
		ConstMemory:  NewMemory(),
		CallStack:    []*ActivationRecord{},
	}
}

func (em *ExecutionMemory) CurrentLocal() *Memory {
	if len(em.CallStack) == 0 {
		return nil
	}
	return em.CallStack[len(em.CallStack)-1].LocalMemory
}

func (em *ExecutionMemory) Set(addr int, value interface{}) error {
	switch {
	case addr >= 1000 && addr < 5000:
		em.GlobalMemory.Set(addr, value)
	case addr >= 5000 && addr < 9000:
		local := em.CurrentLocal()
		if local == nil {
			return fmt.Errorf("no hay memoria local activa para dirección %d", addr)
		}
		local.Set(addr, value)
	case addr >= 9000 && addr < 13000:
		em.TempMemory.Set(addr, value)
	case addr >= 13000 && addr < 17000:
		em.ConstMemory.Set(addr, value)
	default:
		return fmt.Errorf("dirección fuera de rango: %d", addr)
	}
	return nil
}

func (em *ExecutionMemory) Get(addr int) (interface{}, error) {
	switch {
	case addr >= 1000 && addr < 5000:
		return em.GlobalMemory.Get(addr)
	case addr >= 5000 && addr < 9000:
		local := em.CurrentLocal()
		if local == nil {
			return nil, fmt.Errorf("no hay memoria local activa para dirección %d", addr)
		}
		return local.Get(addr)
	case addr >= 9000 && addr < 13000:
		return em.TempMemory.Get(addr)
	case addr >= 13000 && addr < 17000:
		return em.ConstMemory.Get(addr)
	default:
		return nil, fmt.Errorf("dirección fuera de rango: %d", addr)
	}
}
