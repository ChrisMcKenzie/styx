package exec

import (
	"fmt"
	"io"
	"sync"

	"github.com/ChrisMcKenzie/Styx/config"
)

type Responder struct {
	*sync.Mutex
	out io.Writer
	ctx *config.Context
}

func NewResponder(ctx *config.Context, out io.Writer) *Responder {
	return &Responder{&sync.Mutex{}, out, ctx}
}

func (r *Responder) Listen() {
	for _, pl := range r.ctx.Workflow.Pipelines {
		fmt.Println("Processing", pl.Name)
		for _, task := range pl.Tasks {
			fmt.Println("    Running", task.Name)
			fmt.Println("       ", task.Script)
		}
	}
}
