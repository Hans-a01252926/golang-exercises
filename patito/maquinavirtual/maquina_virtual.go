package maquinavirtual

import (
	"fmt"
	"strconv"

	"patito/quadruples"
)

type VirtualMachine struct {
	Quads  []quadruples.Quadruple
	Memory *ExecutionMemory
	IP     int
}

func NewVirtualMachine(quads []quadruples.Quadruple) *VirtualMachine {
	return &VirtualMachine{
		Quads:  quads,
		Memory: NewExecutionMemory(),
		IP:     0,
	}
}

func (vm *VirtualMachine) Run() error {
	for vm.IP < len(vm.Quads) {
		q := vm.Quads[vm.IP]

		switch q.Operator {
		case "+", "-", "*", "/":
			if err := vm.binaryOp(q); err != nil {
				return err
			}
		case "<", ">", "==", "!=":
			if err := vm.relOp(q); err != nil {
				return err
			}
		case "=":
			if err := vm.assign(q); err != nil {
				return err
			}
		case "print":
			if err := vm.print(q); err != nil {
				return err
			}
		case "GOTO":
			dest, _ := strconv.Atoi(q.Result)
			vm.IP = dest
			continue
		case "GOTOF":
			cond, err := vm.getBool(q.Left)
			if err != nil {
				return err
			}
			if !cond {
				dest, _ := strconv.Atoi(q.Result)
				vm.IP = dest
				continue
			}
		case "ERA":
			vm.era(q)
		case "PARAM":
			if err := vm.param(q); err != nil {
				return err
			}
		case "GOSUB":
			vm.gosub(q)
			continue
		case "ENDFunc":
			vm.endFunc()
			continue
		case "RETURN":
			if err := vm.returnFunc(q); err != nil {
				return err
			}
			continue
		default:
			return fmt.Errorf("operador no soportado: %s", q.Operator)
		}

		vm.IP++
	}
	return nil
}

func parseAddr(s string) int {
	addr, _ := strconv.Atoi(s)
	return addr
}

func (vm *VirtualMachine) assign(q quadruples.Quadruple) error {
	src := parseAddr(q.Left)
	dst := parseAddr(q.Result)

	value, err := vm.Memory.Get(src)
	if err != nil {
		return err
	}

	return vm.Memory.Set(dst, value)
}

func (vm *VirtualMachine) print(q quadruples.Quadruple) error {
	addr := parseAddr(q.Left)

	value, err := vm.Memory.Get(addr)
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

func (vm *VirtualMachine) getBool(addrStr string) (bool, error) {
	addr := parseAddr(addrStr)

	value, err := vm.Memory.Get(addr)
	if err != nil {
		return false, err
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("valor en dirección %d no es bool", addr)
	}

	return boolValue, nil
}

func (vm *VirtualMachine) binaryOp(q quadruples.Quadruple) error {
	leftAddr := parseAddr(q.Left)
	rightAddr := parseAddr(q.Right)
	resultAddr := parseAddr(q.Result)

	left, err := vm.Memory.Get(leftAddr)
	if err != nil {
		return err
	}

	right, err := vm.Memory.Get(rightAddr)
	if err != nil {
		return err
	}

	result, err := applyArithmetic(q.Operator, left, right)
	if err != nil {
		return err
	}

	return vm.Memory.Set(resultAddr, result)
}

func (vm *VirtualMachine) relOp(q quadruples.Quadruple) error {
	leftAddr := parseAddr(q.Left)
	rightAddr := parseAddr(q.Right)
	resultAddr := parseAddr(q.Result)

	left, err := vm.Memory.Get(leftAddr)
	if err != nil {
		return err
	}

	right, err := vm.Memory.Get(rightAddr)
	if err != nil {
		return err
	}

	result, err := applyRelational(q.Operator, left, right)
	if err != nil {
		return err
	}

	return vm.Memory.Set(resultAddr, result)
}

func toFloat64(v interface{}) (float64, error) {
	switch n := v.(type) {
	case int:
		return float64(n), nil
	case float64:
		return n, nil
	default:
		return 0, fmt.Errorf("valor no numérico: %v", v)
	}
}

func applyArithmetic(op string, left interface{}, right interface{}) (interface{}, error) {
	l, err := toFloat64(left)
	if err != nil {
		return nil, err
	}

	r, err := toFloat64(right)
	if err != nil {
		return nil, err
	}

	switch op {
	case "+":
		return l + r, nil
	case "-":
		return l - r, nil
	case "*":
		return l * r, nil
	case "/":
		return l / r, nil
	default:
		return nil, fmt.Errorf("operador aritmético no soportado: %s", op)
	}
}

func applyRelational(op string, left interface{}, right interface{}) (bool, error) {
	l, err := toFloat64(left)
	if err != nil {
		return false, err
	}

	r, err := toFloat64(right)
	if err != nil {
		return false, err
	}

	switch op {
	case "<":
		return l < r, nil
	case ">":
		return l > r, nil
	case "==":
		return l == r, nil
	case "!=":
		return l != r, nil
	default:
		return false, fmt.Errorf("operador relacional no soportado: %s", op)
	}
}

func (vm *VirtualMachine) era(q quadruples.Quadruple) {
	vm.Memory.PendingERA = NewActivationRecord(q.Left, vm.IP+1)
}

func (vm *VirtualMachine) param(q quadruples.Quadruple) error {
	argAddr, _ := strconv.Atoi(q.Left)
	paramAddr, _ := strconv.Atoi(q.Result)

	value, err := vm.Memory.Get(argAddr)
	if err != nil {
		return err
	}

	vm.Memory.PendingERA.LocalMemory.Set(paramAddr, value)
	return nil
}

func (vm *VirtualMachine) gosub(q quadruples.Quadruple) {
	dest, _ := strconv.Atoi(q.Result)

	vm.Memory.PendingERA.ReturnIP = vm.IP + 1
	vm.Memory.CallStack = append(vm.Memory.CallStack, vm.Memory.PendingERA)
	vm.Memory.PendingERA = nil

	vm.IP = dest
}

func (vm *VirtualMachine) endFunc() {
	top := vm.Memory.CallStack[len(vm.Memory.CallStack)-1]
	vm.Memory.CallStack = vm.Memory.CallStack[:len(vm.Memory.CallStack)-1]
	vm.IP = top.ReturnIP
}

func (vm *VirtualMachine) returnFunc(q quadruples.Quadruple) error {
	valueAddr, err := strconv.Atoi(q.Left)
	if err != nil {
		return err
	}

	value, err := vm.Memory.Get(valueAddr)
	if err != nil {
		return err
	}

	fmt.Println("RETURN:", value)

	vm.endFunc()
	return nil
}
