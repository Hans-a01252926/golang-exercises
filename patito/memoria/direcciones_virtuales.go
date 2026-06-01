package memoria

import (
	"fmt"
)

type Address int

const (
	GlobalIntStart    Address = 1000
	GlobalFloatStart  Address = 2000
	GlobalStringStart Address = 3000

	LocalIntStart    Address = 5000
	LocalFloatStart  Address = 6000
	LocalStringStart Address = 7000

	TempIntStart   Address = 9000
	TempFloatStart Address = 10000
	TempBoolStart  Address = 11000

	ConstIntStart    Address = 13000
	ConstFloatStart  Address = 14000
	ConstStringStart Address = 15000
)

type AddressAllocator struct {
	nextGlobalInt    Address
	nextGlobalFloat  Address
	nextGlobalString Address

	nextLocalInt    Address
	nextLocalFloat  Address
	nextLocalString Address

	nextTempInt   Address
	nextTempFloat Address
	nextTempBool  Address

	nextConstInt    Address
	nextConstFloat  Address
	nextConstString Address
}

func NewAddressAllocator() *AddressAllocator {
	return &AddressAllocator{
		nextGlobalInt:    GlobalIntStart,
		nextGlobalFloat:  GlobalFloatStart,
		nextGlobalString: GlobalStringStart,

		nextLocalInt:    LocalIntStart,
		nextLocalFloat:  LocalFloatStart,
		nextLocalString: LocalStringStart,

		nextTempInt:   TempIntStart,
		nextTempFloat: TempFloatStart,
		nextTempBool:  TempBoolStart,

		nextConstInt:    ConstIntStart,
		nextConstFloat:  ConstFloatStart,
		nextConstString: ConstStringStart,
	}
}

func (a *AddressAllocator) AllocateGlobal(t string) (Address, error) {
	switch t {
	case "entero":
		addr := a.nextGlobalInt
		a.nextGlobalInt++
		return addr, nil
	case "flotante":
		addr := a.nextGlobalFloat
		a.nextGlobalFloat++
		return addr, nil
	case "string":
		addr := a.nextGlobalString
		a.nextGlobalString++
		return addr, nil
	default:
		return 0, fmt.Errorf("tipo global no soportado: %s", t)
	}
}

func (a *AddressAllocator) AllocateLocal(t string) (Address, error) {
	switch t {
	case "entero":
		addr := a.nextLocalInt
		a.nextLocalInt++
		return addr, nil
	case "flotante":
		addr := a.nextLocalFloat
		a.nextLocalFloat++
		return addr, nil
	case "string":
		addr := a.nextLocalString
		a.nextLocalString++
		return addr, nil
	default:
		return 0, fmt.Errorf("tipo local no soportado: %s", t)
	}
}

func (a *AddressAllocator) AllocateTemp(t string) (Address, error) {
	switch t {
	case "entero":
		addr := a.nextTempInt
		a.nextTempInt++
		return addr, nil
	case "flotante":
		addr := a.nextTempFloat
		a.nextTempFloat++
		return addr, nil
	case "bool":
		addr := a.nextTempBool
		a.nextTempBool++
		return addr, nil
	default:
		return 0, fmt.Errorf("tipo temporal no soportado: %s", t)
	}
}

func (a *AddressAllocator) AllocateConst(t string) (Address, error) {
	switch t {
	case "entero":
		addr := a.nextConstInt
		a.nextConstInt++
		return addr, nil
	case "flotante":
		addr := a.nextConstFloat
		a.nextConstFloat++
		return addr, nil
	case "string":
		addr := a.nextConstString
		a.nextConstString++
		return addr, nil
	default:
		return 0, fmt.Errorf("tipo constante no soportado: %s", t)
	}
}

func (a *AddressAllocator) ResetLocal() {
	a.nextLocalInt = LocalIntStart
	a.nextLocalFloat = LocalFloatStart
	a.nextLocalString = LocalStringStart
}

func (a *AddressAllocator) ResetTemps() {
	a.nextTempInt = TempIntStart
	a.nextTempFloat = TempFloatStart
	a.nextTempBool = TempBoolStart
}
