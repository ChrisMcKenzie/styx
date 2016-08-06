package exec

import (
	"fmt"
	"io"
	"os"
	"sync"

	"golang.org/x/net/context"

	"github.com/ChrisMcKenzie/Styx/config"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

var cli *client.Client

func init() {
	var err error
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err = client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}
}

type Responder struct {
	*sync.Mutex
	out io.Writer
	ctx *config.Context
}

func NewResponder(ctx *config.Context, out io.Writer) *Responder {
	return &Responder{&sync.Mutex{}, out, ctx}
}

func (r *Responder) Exec() {
	fmt.Println("Create Working Environment")
	ctx := context.Background()
	err := createEnvironment(ctx, r.ctx.Workflow)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer destroyEnvironment(ctx, r.ctx.Workflow)

	for _, pl := range r.ctx.Workflow.Pipelines {
		fmt.Println("Processing", pl.Name)
		for _, task := range pl.Tasks {
			fmt.Println("    Running", task.Name)
			fmt.Println("       ", task.Script)

		}
	}
}

func createEnvironment(ctx context.Context, wf *config.Workflow) error {
	out, err := cli.ImagePull(ctx, wf.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		AttachStdout: true,
		Image:        wf.Image,
	}, nil, nil, "")
	if err != nil {
		return err
	}

	wf.ContainerId = resp.ID

	fmt.Println("Created container:", wf.ContainerId)
	return nil
}

func destroyEnvironment(ctx context.Context, wf *config.Workflow) error {
	fmt.Println("Destroying Container:", wf.ContainerId)
	err := cli.ContainerRemove(ctx, wf.ContainerId, types.ContainerRemoveOptions{})
	return err
}
