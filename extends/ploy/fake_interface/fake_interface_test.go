package fake_interface

import (
	"testing"
	"context"
	"fmt"
)

type A struct {
}

func (this *A) Run(ctx context.Context, req string, resp string) context.Context {
	fmt.Printf("%s-%s\n", req, resp)
	return nil
}

type B struct {
}

func (this *B) AddPloy(ctx context.Context, A *FakeInterface) context.Context {
	ctx = A.Do("Run", ctx, "1", "1")
	if ctx != nil && ctx.Err() != nil{
		fmt.Println(ctx.Err())
		return ctx
	}
	return ctx
}

func Test_Fake(t *testing.T) {
	var err error
	Ploy := NewFakeInterface()	
	err = Ploy.Add("Run", 3)
	if err != nil {
		panic(err)
    }
	err = Ploy.Realize(new(A))
	if err != nil {
		panic(err)
    }

	Flow := NewFakeInterface()
	err = Flow.Add("AddPloy", 2)
	if err != nil {
		panic(err)
    }
	err = Flow.Realize(new(B))
	if err != nil {
		panic(err)
    }
	
	for i := 0;i < 5; i++{
		Flow.Do("AddPloy", context.Background(), Ploy)
    }
}
