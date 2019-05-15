package docker

import "fmt"

type Provider struct {
	Image string `hcl:"image"`
}

func (e *Provider) Execute() {
	fmt.Printf("Execute called from docker using %s image\n", e.Image)
}
