package memoria

import (
	"fmt"
)

type Address int

const (
	GlobalIntStart    Address = 1000
	GlobalFloatStart  Address = 2000
	GlobalStringStart Address = 3000
	GlobalBoolStart   Address = 4000

	LocalIntStart    Address = 5000
	LocalFloatStart  Address = 6000
	LocalStringStart Address = 7000
	LocalBoolStart   Address = 8000

	TempIntStart    Address = 9000
	TempFloatStart  Address = 10000
	TempStringStart Address = 11000
	TempBoolStart   Address = 12000

	ConstIntStart    Address = 13000
	ConstFloatStart  Address = 14000
	ConstStringStart Address = 15000
	ConstBoolStart   Address = 16000
)

type AddressAllocator struct {
	nextGlobalInt    Address
	nextGlobalFloat  Address
	nextGlobalString Address
	nextGlobalBool   Address

	nextLocalInt    Address
	nextLocalFloat  Address
	nextLocalString Address
	nextLocalBool   Address

	nextTempInt    Address
	nextTempFloat  Address
	nextTempString Address
	nextTempBool   Address

	nextConstInt    Address
	nextConstFloat  Address
	nextConstString Address
	nextConstBool   Address
}

func NewAddressAllocator() *AddressAllocator {
	return &AddressAllocator{
		nextGlobalInt:    GlobalIntStart,
		nextGlobalFloat:  GlobalFloatStart,
		nextGlobalString: GlobalStringStart,
		nextGlobalBool:   GlobalBoolStart,

		nextLocalInt:    LocalIntStart,
		nextLocalFloat:  LocalFloatStart,
		nextLocalString: LocalStringStart,
		nextLocalBool:   LocalBoolStart,

		nextTempInt:    TempIntStart,
		nextTempFloat:  TempFloatStart,
		nextTempString: TempStringStart,
		nextTempBool:   TempBoolStart,

		nextConstInt:    ConstIntStart,
		nextConstFloat:  ConstFloatStart,
		nextConstString: ConstStringStart,
		nextConstBool:   ConstBoolStart,
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
	case "bool":
		addr := a.nextGlobalBool
		a.nextGlobalBool++
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
	case "bool":
		addr := a.nextLocalBool
		a.nextLocalBool++
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
	case "string":
		addr := a.nextTempString
		a.nextTempString++
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
	case "bool":
		addr := a.nextConstBool
		a.nextConstBool++
		return addr, nil
	default:
		return 0, fmt.Errorf("tipo constante no soportado: %s", t)
	}
}

func (a *AddressAllocator) ResetLocal() {
	a.nextLocalInt = LocalIntStart
	a.nextLocalFloat = LocalFloatStart
	a.nextLocalString = LocalStringStart
	a.nextLocalBool = LocalBoolStart
}

func (a *AddressAllocator) ResetTemps() {
	a.nextTempInt = TempIntStart
	a.nextTempFloat = TempFloatStart
	a.nextTempString = TempStringStart
	a.nextTempBool = TempBoolStart
}
