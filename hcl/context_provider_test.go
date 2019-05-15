package hcl

import (
	"reflect"
	"testing"

	"github.com/ChrisMcKenzie/Styx/styx"
	"github.com/ChrisMcKenzie/Styx/styx/docker"
	"github.com/ChrisMcKenzie/Styx/styx/local"
)

func TestDriverParsing(t *testing.T) {
	tests := []struct {
		name   string
		text   []byte
		result styx.ExecService
	}{
		{
			name: "test local driver",
			text: []byte(`
			workflow "test" {
				driver "local" {}
			}
			`),
			result: &local.ExecService{},
		},
		{
			name: "test docker driver",
			text: []byte(`
			workflow "test" {
				driver "docker" {
					image = "ubuntu:latest"
				}
			}
			`),
			result: &docker.ExecService{
				Image: "ubuntu:latest",
			},
		},
	}

	for _, test := range tests {
		c, err := NewContextProvider(test.text)
		if err != nil {
			t.Error(err)
		}

		ctx, err := c.Context()
		if err != nil {
			t.Error(err)
		}

		if reflect.TypeOf(ctx.Workflows["test"].Driver) != reflect.TypeOf(test.result) {
			t.Errorf("%s: Expected driver to be %v got %v",
				test.name,
				reflect.TypeOf(test.result),
				reflect.TypeOf(ctx.Workflows["test"].Driver))
		}

		if !reflect.DeepEqual(ctx.Workflows["test"].Driver, test.result) {
			t.Errorf("%s: Expected driver config to be %v got %v",
				test.name,
				test.result,
				ctx.Workflows["test"].Driver)
		}
	}
}
