package styx

import "fmt"

type Workflow struct {
	Name        string
	Image       string
	Driver      map[string]interface{}
	Tasks       map[string]*Task
	Variables   map[string]*Variable
	Pipelines   []*Pipeline
	ContainerId string
}

func (c *Context) SetWorkflow(name string) error {
	if ctx, ok := c.Workflows[name]; ok {
		c.Workflow = ctx
		return nil
	}

	return fmt.Errorf("No workflow exists by that name: %s", name)
}
