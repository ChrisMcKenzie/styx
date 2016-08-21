package styx

import "fmt"

type Context struct {
	Name      string
	Workflows map[string]*Workflow
	Tasks     map[string]*Task
	Pipelines map[string]*Pipeline
	Variables []*Variable
	Workflow  *Workflow
}

func (c *Context) SetWorkflow(name string) error {
	if ctx, ok := c.Workflows[name]; ok {
		c.Workflow = ctx
		return nil
	}

	return fmt.Errorf("No workflow exists by that name: %s", name)
}

type ContextProvider interface {
	Context() (*Context, error)
}
