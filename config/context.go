package config

import "fmt"

type Context struct {
	Name      string
	Workflows map[string]*Workflow
	Tasks     map[string]*Task
	Pipelines map[string]*Pipeline
	Variables []*Variable
	Workflow  *Workflow
}

type Workflow struct {
	Name      string
	Tasks     map[string]*Task
	Variables map[string]*Variable
	Pipelines []*Pipeline
}

func (c *Context) SetWorkflow(name string) error {
	if ctx, ok := c.Workflows[name]; ok {
		c.Workflow = ctx
		return nil
	}

	return fmt.Errorf("No workflow exists by that name: %s", name)
}
