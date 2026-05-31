package memory

import (
	"fmt"

	"patito/semantic"
)

type ConstantEntry struct {
	Literal string
	Type    semantic.Type
	Address Address
}

type ConstantsTable struct {
	Constants map[string]ConstantEntry
	Allocator *AddressAllocator
}

func NewConstantsTable(allocator *AddressAllocator) *ConstantsTable {
	return &ConstantsTable{
		Constants: make(map[string]ConstantEntry),
		Allocator: allocator,
	}
}

func constantKey(literal string, t semantic.Type) string {
	return fmt.Sprintf("%s:%s", t, literal)
}

func (ct *ConstantsTable) AddConstant(literal string, t semantic.Type) (Address, error) {
	key := constantKey(literal, t)

	if entry, exists := ct.Constants[key]; exists {
		return entry.Address, nil
	}

	addr, err := ct.Allocator.AllocateConst(t)
	if err != nil {
		return 0, err
	}

	ct.Constants[key] = ConstantEntry{
		Literal: literal,
		Type:    t,
		Address: addr,
	}

	return addr, nil
}
