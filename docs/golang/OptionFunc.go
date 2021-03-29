package main

import "fmt"

type OptionFunc func(op *OptionMenu)

type OptionMenu struct {
	op1 string
	op2 string
	op3 int
	op4 int
	//......可能还会有跟多的属性加入，这给实例化带来了巨大的问题
}

func InitOptions(optionFuncs ...OptionFunc) OptionMenu {
	option := OptionMenu{}
	for _, op := range optionFuncs {
		op(&option)
	}
	return option
}

func WithOp1(op1 string) OptionFunc {
	return func(op *OptionMenu) {
		op.op1 = op1
	}
}

func WithOp2(op2 string) OptionFunc {
	return func(op *OptionMenu) {
		op.op2 = op2
	}
}

func WithOp3(op3 int) OptionFunc {
	return func(op *OptionMenu) {
		op.op3 = op3
	}
}

func WithOp4(op4 int) OptionFunc {
	return func(op *OptionMenu) {
		op.op4 = op4
	}
}

func main() {
	op := InitOptions(
		WithOp1("op1"),
		WithOp3(3),
		WithOp4(4),
	)
	fmt.Printf("%#v\n", op)
}
