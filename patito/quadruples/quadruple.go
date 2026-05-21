package quadruples

import "fmt"

type Quadruple struct {
	Operator string
	Left     string
	Right    string
	Result   string
}

func (q Quadruple) String() string {
	return fmt.Sprintf("(%s, %s, %s, %s)", q.Operator, q.Left, q.Right, q.Result)
}
