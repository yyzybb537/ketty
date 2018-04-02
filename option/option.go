package option

type OptionI interface{}

// 各协议自定义的option必须继承于Option
type Option struct {
	OptionI

	TimeoutMilliseconds int64
}

func DefaultOption() *Option {
	return &Option{}
}
