package fake_interface

import (
	"testing"
	"context"
)

type A struct {
}

func (this *A) Run(ctx context.Context, req string, resp string) context.Context {
	return nil
}

type B struct {
}

func (this *B) Run(ctx context.Context, req string, resp string) context.Context {
	return nil
}

func Test_Fake(t *testing.T) {
	var err error
	a := NewFakeInterface()	
	err = a.Add("Run", 3)
	if err != nil {
		panic(err)
    }
	err = a.Realize(new(A))
	if err != nil {
		panic(err)
    }

	b := NewFakeInterface()
	err = b.Add("Run", 3)
	if err != nil {
		panic(err)
    }
	err = b.Realize(new(B))
	if err != nil {
		panic(err)
    }

	println(a.RealizedTypeName())
	println(a.LookLike(b))
}
