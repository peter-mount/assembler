package common

import (
	"github.com/peter-mount/assembler/assembler/context"
	"github.com/peter-mount/assembler/assembler/lexer"
	"github.com/peter-mount/assembler/assembler/node"
	"github.com/peter-mount/assembler/memory"
	"reflect"
	"testing"
)

func TestGetNodeAddress(t *testing.T) {
	type args struct {
		n   *node.Node
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    memory.Address
		wantErr bool
	}{
		// TODO: Handle test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeAddress(tt.args.n, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNodeAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNodeInt(t *testing.T) {
	type args struct {
		n   *node.Node
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Handle test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeInt(tt.args.n, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNodeInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNodeInterface(t *testing.T) {
	type args struct {
		n   *node.Node
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Handle test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeInterface(tt.args.n, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNodeInterface() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToAddr(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    memory.Address
		wantErr bool
	}{
		// TODO: Handle test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToAddr(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToAddr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expected  int64
		wantError bool
	}{
		{name: "int64", value: int64(42), expected: 42},
		{name: "int", value: 32767, expected: 32767},
		{name: "Label", value: memory.Address(0x6543), expected: 0x6543},
		{name: "Line", value: &lexer.Line{Address: memory.Address(0x6543)}, expected: 0x6543},
		{name: "string", value: "1024", expected: 1024},
		{name: "NAN", value: "This is not a number", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt(tt.value)
			if err != nil {
				if !tt.wantError {
					t.Error(err)
				}
			} else if got != tt.expected {
				t.Errorf("Got %d expected %d", got, tt.expected)
			}
		})
	}
}
