package term

import "fmt"
import "strconv"
import "math/big"

type Float float64

func NewFloat(text string) (Number,error) {
	r, ok := NewRational(text)
	if ok {
		return r
	}
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return (*Float)(&f),err
	}
	return (*Float)(&f)
}

func NewFloat64(f float64) (*Float ,error) {
	return (*Float)(&f) , nil
}

func (self *Float) Value() float64 {
	return float64(*self)
}

func (self *Float) String() string {
	return fmt.Sprintf("%g", self.Value())
}
func (self *Float) Type() int {
	return FloatType
}
func (self *Float) Indicator() string {
	return self.String()
}

func (a *Float) Unify(e Bindings, b Term) (Bindings, error) {
	if IsVariable(b) {
		return b.Unify(e, a)
	}
	if IsFloat(b) {
		if a.Value() == b.(*Float).Value() {
			return e, nil
		}
	}

	return e, CantUnify
}

func (self *Float) ReplaceVariables(env Bindings) Term {
	return self
}

// implement Number interface
func (self *Float) Float64() float64 {
	return self.Value()
}

func (self *Float) LosslessInt() (*big.Int, bool) {
	return nil, false
}

func (self *Float) LosslessRat() (*big.Rat, bool) {
	return nil, false
}
