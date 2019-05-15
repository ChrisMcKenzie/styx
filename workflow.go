package styx

type Workflow struct {
	Name      string
	Driver    ExecService
	Tasks     map[string]*Task
	Variables map[string]*Variable
	Pipelines []*Pipeline
}

func (wf *Workflow) Run() error {
	return nil
}
