package maquinavirtual

type ActivationRecord struct {
	FunctionName string
	LocalMemory  *Memory
	ReturnIP     int
}

func NewActivationRecord(functionName string, returnIP int) *ActivationRecord {
	return &ActivationRecord{
		FunctionName: functionName,
		LocalMemory:  NewMemory(),
		ReturnIP:     returnIP,
	}
}
