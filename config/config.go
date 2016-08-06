package config

type Config struct {
	Workspaces []*Context
	Tasks      []*Task
}

type Variable struct {
	Name    string
	Default interface{}
}

type Task struct {
	Name     string
	Script   string
	FilePath string
}
