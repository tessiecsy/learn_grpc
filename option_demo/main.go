package main

import "fmt"

const defaultName = "tessie"

// 内部字段小写，避免向外暴露，只能通过特定函数WithconfigOptionName来更改
type configOption struct {
	name string
	age  int
}

// 模仿构造函数，新建一个configOption实例
func NewConfigOption(age int, opts ...ConfigOption) *configOption {
	cfg := &configOption{
		name: defaultName,
		age:  age,
	}
	// opts是可选的，可传可不传
	for _, opt := range opts {
		// 每个opt执行apply方法
		opt.apply(cfg)
	}
	return cfg
}

// ConfigOption接口实现apply方法
type ConfigOption interface {
	apply(*configOption)
}

type funcOption struct {
	f func(*configOption)
}

func (f funcOption) apply(cfg *configOption) {
	f.f(cfg)
}

// 模拟构造函数
func NewfuncOption(f func(*configOption)) funcOption {
	return funcOption{f: f}
}

func WithconfigOptionName(name string) ConfigOption {
	return NewfuncOption(func(cfg *configOption) {
		cfg.name = name
	})
}

func main() {
	cfg := NewConfigOption(18)
	fmt.Printf("cfg:%#v\n", cfg)
	cfg2 := NewConfigOption(18, WithconfigOptionName("csy"))
	fmt.Printf("cfg:%#v\n", cfg2)

}

// f是一个函数，它的参数与apply方法相同，并且定义了apply方法，并在apply中调用自己
// 那么此函数f实现了接口ConfigOption，称为接口型函数
// type f1 func(*configOption)

// func (fx f1) apply(cfg *configOption) {
// 	fx(cfg)
// }
