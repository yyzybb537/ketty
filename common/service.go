package common

type ServiceHandle interface {
	Implement() interface{}

	ServiceName() string
}
