package tunnel

type Config struct {
	Addr string
	Keys []string
	Auth map[string]string
}
