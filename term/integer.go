package term

import "fmt"
import "math/big"

// Integer represents an unbounded, signed integer value
type Integer big.Int

// NewInt parses an integer's string representation to create a new
// integer value. Panics if the string's is not a valid integer
func NewInt(text string) Number {
	if len(text) == 0 {
		panic("Empty string is not a valid integer")
	}

	// see §6.4.4 for syntax details
	if text[0] == '0' && len(text) >= 3 {
		switch text[1] {
		case '\'':
			return parseEscape(text[2:])
		case 'b':
			return parseInteger("%b", text[2:])
		case 'o':
			return parseInteger("%o", text[2:])
		case 'x':
			return parseInteger("%x", text[2:])
		default:
			return parseInteger("%d", text)
		}
	}
	return parseInteger("%d", text)
}

// helper for when an int64 is already available
func NewInt64(i int64) Number {
	return (*Integer)(big.NewInt(i))
}

// helper for when a big.Int is already available
func NewBigInt(val *big.Int) Number {
	return (*Integer)(val)
}

// NewCode returns an integer whose value is the character code if
// the given rune.
func NewCode(c rune) *Integer {
	return (*Integer)(big.NewInt(int64(c)))
}

func parseInteger(format, text string) (*Integer, error) {
	i := new(big.Int)
	n, err := fmt.Sscanf(text, format, i)
	//maybePanic(err)
	if err != nil {
		return (*Integer)(i), err
	}
	if n == 0 {
		panic("Parsed no integers")
	}

	return (*Integer)(i)
}

// see "single quoted character" - §6.4.2.1
func parseEscape(text string) (*Integer, error) {
	var r rune
	if text[0] == '\\' {
		if len(text) < 2 {
			err := fmt.Errorf("Invalid integer character constant: %s", text)
			//panic(msg)
			return 0, err
		}
		switch text[1] {
		// "meta escape sequence" - §6.4.2.1
		case '\\':
			r = '\\'
		case '\'':
			r = '\''
		case '"':
			r = '"'
		case '`':
			r = '`'

		// "control escape char" - §6.4.2.1
		case 'a':
			r = '\a'
		case 'b':
			r = '\b'
		case 'f':
			r = '\f'
		case 'n':
			r = '\n'
		case 'r':
			r = '\r'
		case 's':
			r = ' ' // SWI-Prolog extension
		case 't':
			r = '\t'
		case 'v':
			r = '\v'

		// "hex escape char" - §6.4.2.1
		case 'x':
			return parseInteger("%x", text[2:len(text)-1])

		// "octal escape char" - §6.4.2.1
		case '0', '1', '2', '3', '4', '5', '6', '7':
			return parseInteger("%o", text[1:len(text)-1])

		// unexpected escape sequence
		default:
			err := fmt.Errorf("Invalid character escape sequence: %s", text)
			//panic(msg)
			return 0, err
		}
	} else {
		// "non quote char" - §6.4.2.1
		runes := []rune(text)
		r = runes[0]
	}
	code := int64(r)
	return (*Integer)(big.NewInt(code)), nil
}

func (self *Integer) Value() *big.Int {
	return (*big.Int)(self)
}

// treat this integer as a character code. should be a method on
// a Code interface someday
func (self *Integer) Code() rune {
	i := (*big.Int)(self)
	return rune(i.Int64())
}

func (self *Integer) String() string {
	return self.Value().String()
}

func (self *Integer) Type() int {
	return IntegerType
}

func (self *Integer) Indicator() string {
	return self.String()
}

func (a *Integer) Unify(e Bindings, b Term) (Bindings, error) {
	if IsVariable(b) {
		return b.Unify(e, a)
	}
	if IsInteger(b) {
		if a.Value().Cmp(b.(*Integer).Value()) == 0 {
			return e, nil
		}
	}

	return e, CantUnify
}

func (self *Integer) ReplaceVariables(env Bindings) Term {
	return self
}

// implement Number interface
func (self *Integer) Float64() float64 {
	return float64(self.Value().Int64())
}

func (self *Integer) LosslessInt() (*big.Int, bool) {
	return self.Value(), true
}

func (self *Integer) LosslessRat() (*big.Rat, bool) {
	r := new(big.Rat).SetFrac(self.Value(), big.NewInt(1))
	return r, true
}
