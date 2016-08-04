package config

type Config struct {
	Workspaces []*Context
	Tasks      []*Task
}

type Context struct {
	Name      string
	Workflows []*Context
	Pipelines []*Pipeline
	Tasks     []*Task
	Variables []*Variable
}

type Variable struct {
	Name    string
	Default interface{}
}

type Pipeline struct {
	Name     string
	Triggers []*Event
	Tasks    []*Task
}

type Task struct {
	Name     string
	Script   *string
	FilePath *string
}
