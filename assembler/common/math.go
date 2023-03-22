package common

import (
	"assembler/assembler/errors"
	"assembler/memory"
	"assembler/util"
)

func ToInt(v interface{}) (int64, error) {
	if i, ok := v.(int64); ok {
		return i, nil
	}
	if i, ok := v.(int); ok {
		return int64(i), nil
	}
	if s, ok := v.(string); ok {
		return util.Atoi(s)
	}
	return 0, errors.IllegalArgument()
}

func ToAddr(v interface{}) (memory.Address, error) {
	i, err := ToInt(v)
	if err != nil {
		return 0, err
	}

	return memory.Address(i), nil
}
