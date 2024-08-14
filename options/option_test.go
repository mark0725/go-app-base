package options

import (
	"fmt"
	"testing"
)

func TestMapOptions(t *testing.T) {
	fmt.Println(NewMapOptions(
		With("a", 1),
		With("b", 2),
		With("c", 3),
	))
}

type MyOptions struct {
	A int
	B string
}

func TestOptions(t *testing.T) {
	options := MyOptions{}
	err := NewOptions(&options,
		With("A", 1),
		With("B", "HELLO"),
	)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", options)

}
