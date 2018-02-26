package ketty

type ServiceHandle interface {
	Implement() interface{}

	ServiceName() string
}
