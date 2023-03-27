package context

import (
	"fmt"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/memory"
	"github.com/peter-mount/assembler/util"
	"reflect"
	"strconv"
)

// The types of a Value returned by Type()
const (
	_ = iota // Ignore 0 so use _ then if someone manually creates Value it's an unknown type
	VarNull
	VarBool
	VarInt
	VarLabel
	VarFloat
	VarComplex
	VarString
)

var (
	nullValue  = Value{varType: VarNull}
	falseValue = Value{varType: VarBool, boolVal: false}
	trueValue  = Value{varType: VarBool, boolVal: true}
	zeroValue  = Value{varType: VarInt, intVal: 0}
)

// Value an immutable Value of some kind.
type Value struct {
	varType    int
	boolVal    bool
	intVal     int64
	labelVal   *lexer.Line
	floatVal   float64
	complexVal complex128
	stringVal  string
}

func (v *Value) Same(b *Value) bool {
	if b == nil {
		return false
	}

	if v == b {
		return true
	}

	return v.varType == b.varType &&
		v.boolVal == b.boolVal &&
		v.intVal == b.intVal &&
		v.labelVal == b.labelVal &&
		v.floatVal == b.floatVal &&
		v.complexVal == b.complexVal &&
		v.stringVal == b.stringVal
}

func (v *Value) Equal(b *Value) bool {
	switch v.Type() {
	case VarBool:
		return v.Bool() == b.Bool()
	case VarInt:
		return v.Int() == b.Int()
	case VarLabel:
		return v.Label() == b.Label()
	case VarFloat:
		return v.Float() == b.Float()
	case VarComplex:
		return v.Complex() == b.Complex()
	case VarString:
		return v.String() == b.String()
	default:
		return false
	}
}

// NullValue returns the Value for Null/nil
func NullValue() *Value {
	return &nullValue
}

// Type of this value
func (v *Value) Type() int {
	return v.varType
}

// BoolValue returns a Value for a bool
func BoolValue(i bool) *Value {
	if i {
		return &trueValue
	}
	return &falseValue
}

// IntValue returns a Value for an int64
func IntValue(i int64) *Value {
	return &Value{varType: VarInt, intVal: i}
}

func LabelValue(a *lexer.Line) *Value {
	return &Value{varType: VarLabel, labelVal: a}
}

// FloatValue returns a Value for an float64
func FloatValue(i float64) *Value {
	return &Value{varType: VarFloat, floatVal: i}
}

// ComplexValue returns a Value for a complex128
func ComplexValue(i complex128) *Value {
	return &Value{varType: VarComplex, complexVal: i}
}

// StringValue returns a Value for an string
func StringValue(i string) *Value {
	return &Value{varType: VarString, stringVal: i}
}

// IsNull returns true if the Value is null
func (v *Value) IsNull() bool {
	return v.varType == VarNull
}

// IsZero returns true if the Value is null, false, 0, 0.0 or "" dependent on it's type
func (v *Value) IsZero() bool {
	switch v.varType {
	case VarNull:
		return true
	case VarBool:
		return !v.boolVal
	case VarInt:
		return v.intVal == 0
	case VarLabel:
		return v.labelVal == nil || v.labelVal.Address == 0
	case VarFloat:
		return v.floatVal == 0.0
	case VarComplex:
		return v.complexVal == 0+0i
	case VarString:
		return v.stringVal == ""
	default:
		return false
	}
}

// IsNumeric returns true if the Value is a number, i.e. int64 or float64
func (v *Value) IsNumeric() bool {
	return v.varType == VarInt || v.varType == VarFloat || v.varType == VarLabel
}

func (v *Value) IsComplex() bool {
	return v.varType == VarComplex
}

func (v *Value) IsNegative() bool {
	switch v.varType {
	case VarInt:
		return v.intVal < 0
	case VarFloat:
		return v.floatVal < 0.0
	default:
		return false
	}
}

// Bool returns the value as a bool.
// For non-booleans then 0=false, any=true.
// For strings then if the string starts with one of t, T, y, Y or 1 then true
func (v *Value) Bool() bool {
	switch v.varType {
	case VarBool:
		return v.boolVal
	case VarInt:
		if v.intVal != 0 {
			return false
		}
		return true
	case VarLabel:
		return v.labelVal == nil || v.labelVal.Address == 0
	case VarFloat:
		if v.floatVal != 0 {
			return false
		}
		return true
	case VarComplex:
		if v.complexVal != 0+0i {
			return false
		}
		return true
	case VarString:
		if v.stringVal == "" {
			return false
		}
		s := v.stringVal[0]
		return s == 't' || s == 'T' || s == 'y' || s == 'Y' || s == '1'
	default:
		return false
	}
}

// Int returns the value as an int64
func (v *Value) Int() int64 {
	switch v.varType {
	case VarBool:
		if v.boolVal {
			return 1
		}
		return 0
	case VarInt:
		return v.intVal
	case VarLabel:
		if v.labelVal == nil {
			return 0
		}
		return int64(v.labelVal.Address)
	case VarFloat:
		return int64(v.floatVal)
	case VarComplex:
		return int64(real(v.complexVal))
	case VarString:
		r, err := util.Atoi(v.stringVal)
		if err == nil {
			return r
		}
		// Panic instead?
		return 0
	default:
		return 0
	}
}

func (v *Value) Label() memory.Address {
	switch v.varType {
	case VarBool:
		if v.boolVal {
			return 1
		}
		return 0
	case VarInt:
		return memory.Address(v.intVal)
	case VarLabel:
		if v.labelVal == nil {
			return 0
		}
		return v.labelVal.Address
	case VarFloat:
		return memory.Address(int64(v.floatVal))
	case VarComplex:
		return memory.Address(int64(real(v.complexVal)))
	case VarString:
		r, err := util.Atoi(v.stringVal)
		if err == nil {
			return memory.Address(r)
		}
		// Panic instead?
		return 0
	default:
		return 0
	}
}

// Float returns the value as a float64
func (v *Value) Float() float64 {
	switch v.varType {
	case VarBool:
		if v.boolVal {
			return 1.0
		}
		return 0.0
	case VarInt:
		return float64(v.intVal)
	case VarLabel:
		if v.labelVal == nil {
			return 0.0
		}
		return float64(v.labelVal.Address)
	case VarFloat:
		return v.floatVal
	case VarComplex:
		return real(v.complexVal)
	case VarString:
		r, err := strconv.ParseFloat(v.stringVal, 64)
		if err == nil {
			return r
		}
		// Panic instead?
		return 0.0
	default:
		return 0.0
	}
}

// Complex returns the complex value of this Value.
// For real numbers this will return the complex value with 0 for the imaginary component.
func (v *Value) Complex() complex128 {
	if v.varType == VarComplex {
		return v.complexVal
	}
	return complex(v.Float(), 0)
}

// Real returns the real value of this Value.
// For real numbers this is the same as Float() but for complex numbers this
// returns the real component.
func (v *Value) Real() float64 {
	if v.varType == VarComplex {
		return real(v.complexVal)
	}
	return v.Float()
}

// Imaginary returns the imaginary value of this Value.
// For real numbers this returns 0.0 but for complex numbers this returns the
// imaginary component.
func (v *Value) Imaginary() float64 {
	if v.varType == VarComplex {
		return imag(v.complexVal)
	}
	return 0.0
}

// String returns the value as a string
func (v *Value) String() string {
	switch v.varType {
	case VarBool:
		if v.boolVal {
			return "true"
		}
		return "false"
	case VarInt:
		return strconv.FormatInt(v.intVal, 10)
	case VarFloat:
		return strconv.FormatFloat(v.floatVal, 'f', 10, 64)
	case VarLabel:
		if v.labelVal == nil {
			return ""
		}
		return "0x" + strconv.FormatInt(int64(v.labelVal.Address), 16)
	case VarComplex:
		return fmt.Sprintf("%v", v.complexVal)
	case VarString:
		return v.stringVal
	default:
		return ""
	}
}

func Of(v interface{}) *Value {
	if a, ok := v.(*Value); ok {
		return a
	}
	if a, ok := v.(Value); ok {
		return &a
	}
	if i, ok := v.(int); ok {
		return IntValue(int64(i))
	}
	if i, ok := v.(int64); ok {
		return IntValue(i)
	}
	if a, ok := v.(memory.Address); ok {
		return IntValue(int64(a))
	}
	if a, ok := v.(*lexer.Line); ok {
		return LabelValue(a)
	}
	if f, ok := v.(float64); ok {
		return FloatValue(f)
	}
	if b, ok := v.(bool); ok {
		return BoolValue(b)
	}
	if c, ok := v.(complex128); ok {
		return ComplexValue(c)
	}
	if s, ok := v.(string); ok {
		return StringValue(s)
	}
	// TODO how to fail here?
	//panic(errors.IllegalArgument())
	t := v.(reflect.Type)
	panic(fmt.Errorf("cannot handle %v (%v) for Value", t, v))
}

// OperationType returns the type of the suggested value when performing some
// operation like addition or multiplication to keep the precision of the result.
// For example, if a Value is an Integer but the passed value is Float then
// we should use float.
func (v *Value) OperationType(b *Value) int {
	t := v.Type()
	if v.Type() == VarString || b.Type() == VarString {
		t = VarString
	} else if v.Type() == VarComplex || b.Type() == VarComplex {
		t = VarComplex
	} else if v.Type() == VarFloat || b.Type() == VarFloat {
		t = VarFloat
	} else if v.Type() == VarInt || b.Type() == VarInt || v.Type() == VarLabel || b.Type() == VarLabel {
		t = VarInt
	}
	return t
}
