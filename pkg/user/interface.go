package user

type DeviceInterface interface {
	LoadConfig(file string) error
	Headers() map[string]string
	QueryParams() map[string]string
}
