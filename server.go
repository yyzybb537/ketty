package ketty

type Server interface {
	AopListI

	RegisterMethod(handle ServiceHandle, implement interface{}) error

	Serve() error
}
