package framework

import "testing"

//

type A struct {
	a string
}

func TestAdd(t *testing.T) {

	o := &A{}
	if o == nil {
		t.Error("o 是nil")
	}

}
