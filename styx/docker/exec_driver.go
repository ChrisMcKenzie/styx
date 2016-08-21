package docker

import "fmt"

type ExecService struct {
	Image string `hcl:"image"`
}

func (e *ExecService) Execute() {
	fmt.Printf("Execute called from docker using %s image\n", e.Image)
}
