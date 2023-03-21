package memory

import "errors"

var (
	addressInvalid = errors.New("address invalid")
	memoryOverlap  = errors.New("memory overlap")
	readOnlyMemory = errors.New("read only memory")
)

func IsAddressInvalid(err error) bool { return err == addressInvalid }

func IsMemoryOverlap(err error) bool { return err == memoryOverlap }

func IsReadOnlyMemory(err error) bool { return err == readOnlyMemory }
