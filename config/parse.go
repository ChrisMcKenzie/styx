package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type hclContext struct {
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
	script string
	Fields []string `hcl:",decodedFields"`
}

func LoadHclContext(file []byte) (*hclContext, error) {
	t, err := hcl.ParseBytes(file)
	if err != nil {
		return nil, err
	}

	result := &hclContext{
		Bytes: file,
		Root:  t,
	}

	return result, nil
}

func (c *hclContext) Context() (*Context, error) {

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

	ctx := new(Context)
	ctx.Name = "Root"
	if len(rawContext.Variable) > 0 {
		ctx.Variables = make([]*Variable, 0, len(rawContext.Variable))

		for k, v := range rawContext.Variable {
			newVariable := &Variable{
				Name:    k,
				Default: v.Default,
			}
			ctx.Variables = append(ctx.Variables, newVariable)
		}
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
		ctx.Pipelines, err = loadPipelines(pipelines)
		if err != nil {
			return nil, err
		}
	}

	if workflows := list.Filter("workflow"); len(workflows.Items) > 0 {
		var err error
		ctx.Workflows, err = loadWorkflow(workflows)
		if err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

func loadPipelines(list *ast.ObjectList) ([]*Pipeline, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result []*Pipeline

	var rawPipeline struct {
		Task map[string]*hclTask
	}

	for _, item := range list.Items {
		k := item.Keys[0].Token.Value().(string)

		var pipeline = new(Pipeline)
		if err := hcl.DecodeObject(&rawPipeline, item.Val); err != nil {
			return nil, fmt.Errorf(
				"Error reading config for %s: %s",
				k,
				err)
		}

		if tasks := list.Filter("task"); len(tasks.Items) > 0 {
			var err error
			pipeline.Tasks, err = loadTasks(tasks)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, pipeline)
	}

	return result, nil
}

func loadTasks(list *ast.ObjectList) ([]*Task, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result []*Task

	for _, item := range list.Items {
		k := item.Keys[0].Token.Value().(string)
		var rawTask *hclTask

		var task = new(Task)
		if err := hcl.DecodeObject(&rawTask, item.Val); err != nil {
			return nil, fmt.Errorf(
				"Error reading config for %s: %s",
				k,
				err)
		}

		task.Name = k
		task.Script = &rawTask.script

		result = append(result, task)
	}

	return result, nil
}

func loadWorkflow(list *ast.ObjectList) ([]*Context, error) {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil, nil
	}

	var result []*Context

	for _, item := range list.Items {
		k := item.Keys[0].Token.Value().(string)
		i := item.Val.(*ast.ObjectType).List

		var rawContext struct {
			Variable map[string]*hclVariable
			Task     map[string]*hclTask
		}
		if err := hcl.DecodeObject(&rawContext, item); err != nil {
			return nil, err
		}

		ctx := new(Context)
		ctx.Name = k
		if len(rawContext.Variable) > 0 {
			ctx.Variables = make([]*Variable, 0, len(rawContext.Variable))
			for k, v := range rawContext.Variable {
				newVariable := &Variable{
					Name:    k,
					Default: v.Default,
				}
				ctx.Variables = append(ctx.Variables, newVariable)
			}
		}

		if pipelines := i.Filter("pipeline"); len(pipelines.Items) > 0 {
			var err error
			ctx.Pipelines, err = loadPipelines(pipelines)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, ctx)
	}

	return result, nil
}
