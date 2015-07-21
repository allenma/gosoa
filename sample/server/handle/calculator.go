// calculator
package handle

import (
	"github.com/allenma/gosoa/sample/shared"
	"github.com/allenma/gosoa/sample/tutorial"
)

type CalculatorImpl struct {
}

func (c *CalculatorImpl) Ping() (err error) {
	err = nil
	return
}

func (c *CalculatorImpl) Add(num1 int32, num2 int32) (r int32, err error) {
	r = num1 + num2
	err = nil
	return
}

func (c *CalculatorImpl) Calculate(logid int32, w *tutorial.Work) (r int32, err error) {
	return 0, nil
}

func (c *CalculatorImpl) Zip() (err error) {
	return nil
}

func (c *CalculatorImpl) GetStruct(key int32) (r *shared.SharedStruct, err error) {
	r = &shared.SharedStruct{1, "value1"}
	err = nil
	return
}
