package hcl

import (
	"fmt"

	"github.com/ChrisMcKenzie/Styx/styx"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type contextProvider struct {
	Bytes []byte
	Root  *ast.File
}

type hclVariable struct {
	Default      interface{}
	Description  string
	DeclaredType string   `hcl:"type"`
	Fields       []string `hcl:",decodedFields"`
}

type hclTask struct {
	Script string
	Fields []string `hcl:",decodedFields"`
}

func NewContextProvider(file []byte) (*contextProvider, error) {
	t, err := hcl.ParseBytes(file)
	if err != nil {
		return nil, err
	}

	result := &contextProvider{
		Bytes: file,
		Root:  t,
	}

	return result, nil
}

func (c *contextProvider) Context() (*styx.Context, error) {

	var rawContext struct {
		Variable map[string]*hclVariable
		Task     map[string]*hclTask
	}

	// Top-level item should be the object list
	list, ok := c.Root.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	if err := hcl.DecodeObject(&rawContext, list); err != nil {
		return nil, err
	}

	ctx := new(styx.Context)
	ctx.Name = "Root"
	if len(rawContext.Variable) > 0 {
		ctx.Variables = make([]*styx.Variable, 0, len(rawContext.Variable))

		// TODO(ChrisMcKenzie): Variables should be interpolated
		for k, v := range rawContext.Variable {
			newVariable := &styx.Variable{
				Name:    k,
				Default: v.Default,
			}
			ctx.Variables = append(ctx.Variables, newVariable)
		}
	}

	if includes := list.Filter("include"); len(includes.Items) > 0 {
		fmt.Println(includes)
	}

	if tasks := list.Filter("task"); len(tasks.Items) > 0 {
		var err error
		ctx.Tasks, err = loadTasks(tasks)
		if err != nil {
			return nil, err
		}
	}

	if pipelines := list.Filter("pipeline"); len(pipelines.Items) > 0 {
		var err error
		ctx.Pipelines, err = loadRootPipelines(pipelines, ctx)
		if err != nil {
			return nil, err
		}
	}

	if workflows := list.Filter("workflow"); len(workflows.Items) > 0 {
		var err error
		ctx.Workflows, err = loadWorkflow(workflows, ctx)
		if err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

func loadPipelines(list *ast.ObjectList, wf *styx.Workflow, ctx *styx.Context) ([]*styx.Pipeline, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result []*styx.Pipeline

	var rawPipeline struct {
		Task map[string]*hclTask
	}

	for _, item := range list.Items {
		k := item.Keys[0].Token.Value().(string)
		i := item.Val.(*ast.ObjectType).List

		var pipeline = new(styx.Pipeline)
		if p, ok := ctx.Pipelines[k]; ok {
			pipeline = p
		} else {
			pipeline.Name = k
			if err := hcl.DecodeObject(&rawPipeline, item.Val); err != nil {
				return nil, fmt.Errorf(
					"Error reading styx for %s: %s",
					k,
					err)
			}

			if tasks := i.Filter("task"); len(tasks.Items) > 0 {
				var err error
				pipeline.Tasks, err = loadPipelineTasks(tasks, wf, ctx)
				if err != nil {
					return nil, err
				}
			}
		}

		result = append(result, pipeline)
	}

	return result, nil
}

func loadRootPipelines(list *ast.ObjectList, ctx *styx.Context) (map[string]*styx.Pipeline, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result = make(map[string]*styx.Pipeline)

	var rawPipeline struct {
		Task map[string]*hclTask
	}

	for _, item := range list.Items {
		k := item.Keys[0].Token.Value().(string)
		i := item.Val.(*ast.ObjectType).List

		var pipeline = new(styx.Pipeline)
		pipeline.Name = k
		if err := hcl.DecodeObject(&rawPipeline, item.Val); err != nil {
			return nil, fmt.Errorf(
				"Error reading styx for %s: %s",
				k,
				err)
		}

		if tasks := i.Filter("task"); len(tasks.Items) > 0 {
			var err error
			pipeline.Tasks, err = loadPipelineTasks(tasks, new(styx.Workflow), ctx)
			if err != nil {
				return nil, err
			}
		}

		result[k] = pipeline
	}

	return result, nil
}

func loadPipelineTasks(list *ast.ObjectList, wf *styx.Workflow, ctx *styx.Context) ([]*styx.Task, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result []*styx.Task

	for _, item := range list.Items {
		k := item.Keys[0].Token.Value().(string)

		var task = new(styx.Task)
		var rawTask *hclTask

		if err := hcl.DecodeObject(&rawTask, item.Val); err != nil {
			return nil, fmt.Errorf(
				"Error reading styx for %s: %s",
				k,
				err)
		}

		t, globalTaskExists := ctx.Tasks[k]
		wfTask, wfTaskExists := wf.Tasks[k]

		switch {
		case rawTask.Script != "":
			task.Name = k
			task.Script = rawTask.Script
		case wfTaskExists:
			task = wfTask
		case globalTaskExists:
			task = t
		}

		result = append(result, task)
	}

	return result, nil
}

func loadTasks(list *ast.ObjectList) (map[string]*styx.Task, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result = make(map[string]*styx.Task)

	for _, item := range list.Items {
		k := item.Keys[0].Token.Value().(string)
		var rawTask *hclTask

		var task = new(styx.Task)
		if err := hcl.DecodeObject(&rawTask, item.Val); err != nil {
			return nil, fmt.Errorf(
				"Error reading styx for %s: %s",
				k,
				err)
		}

		task.Name = k
		task.Script = rawTask.Script

		result[k] = task
	}

	return result, nil
}

func loadWorkflow(list *ast.ObjectList, c *styx.Context) (map[string]*styx.Workflow, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result = make(map[string]*styx.Workflow)

	for _, item := range list.Items {
		k := item.Keys[0].Token.Value().(string)
		i := item.Val.(*ast.ObjectType).List

		var rawContext struct {
			Image    string
			Variable map[string]*hclVariable
			Task     map[string]*hclTask
		}
		if err := hcl.DecodeObject(&rawContext, item); err != nil {
			return nil, err
		}

		ctx := new(styx.Workflow)
		ctx.Name = k
		ctx.Image = rawContext.Image
		if len(rawContext.Variable) > 0 {
			ctx.Variables = make(map[string]*styx.Variable)
			for k, v := range rawContext.Variable {
				newVariable := &styx.Variable{
					Name:    k,
					Default: v.Default,
				}
				ctx.Variables[k] = newVariable
			}
		}

		if tasks := i.Filter("task"); len(tasks.Items) > 0 {
			var err error
			ctx.Tasks, err = loadTasks(tasks)
			if err != nil {
				return nil, err
			}
		}

		if pipelines := i.Filter("pipeline"); len(pipelines.Items) > 0 {
			var err error
			ctx.Pipelines, err = loadPipelines(pipelines, ctx, c)
			if err != nil {
				return nil, err
			}
		}

		result[k] = ctx
	}

	return result, nil
}
