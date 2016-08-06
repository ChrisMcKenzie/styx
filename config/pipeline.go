package config

type Pipeline struct {
	Name     string
	Triggers []*Event
	Tasks    []*Task
}
