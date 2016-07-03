# Tasks (07/02/2016)

### ___!!! This Proposal is a WIP !!!___

Changes:
  - (07/02/2016): Created

Tasks allow users to easily script complex operations that can be reused in other
parts of their styx configurations. Tasks can behave much like pipelines
in a workflow in that if they are empty Styx will try to find a global task by
the given name otherwise the task is considered local and will be parsed
accordingly.

## Example

```hcl
task "go-setup-workspace" {
  file = "./workspace-setup.sh"
}

task "go-get" {
  script = <<EOF
  #!/bin/bash
  go get ./..
EOF
}

pipeline "setup" {
  task "go-setup-workspace" {}
  task "go-get" {}
}
```

This defines two tasks one of which will use a file in the repo to execute the
other is an inline script, below that you can see their usage in a global 
pipeline.

## Execution
Tasks are executed in the docker context that is created by the workflow or 
pipeline, and the Executor will read the exit code of the script and will flag
the task as fail if the code is > 0. The Executor will also read the stdio of 
the script and output that to the user for any debugging.
