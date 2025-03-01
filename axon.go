package axon

type Axon interface {
	DbalDriver
	LoggerProvider
	ConfigProvider
}
