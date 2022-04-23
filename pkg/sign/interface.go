package sign

type SignInterface interface {
	Sign(data interface{}) (map[string]string, error)
}
