package config

type contextual interface {
	Context() (*Context, error)
}
