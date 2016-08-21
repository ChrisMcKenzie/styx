package styx

type Context struct {
	Name      string
	Workflows map[string]*Workflow
	Tasks     map[string]*Task
	Pipelines map[string]*Pipeline
	Variables []*Variable
	Workflow  *Workflow
}

type ContextProvider interface {
	Context() (*Context, error)
}
