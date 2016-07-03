# Workflows (07/02/2016)

### ___!!! This Proposal is a WIP !!!___

Changes:
  - (07/02/2016): Created

Workflows allow a user to define a custom "plan" for how their project will be
built, tested and deployed. This paradigm allows users to not only define how
their software is integrated and delivered but also allows users to define how
their sofware is developed.

A workflow should allow the definition of [pipelines] that make up a workflow as
well as allowing users to include globally defined pipelines.

## Example
An example `Workflow` structure might look like the following:

```hcl
# globally defined pipelines
pipeline "build" {
  task "npm-build" {}
}

# workflow for development
workflow "development" {
  # Use image from docker hub
  image = "node:0.6.0"

  # bring in globally defined pipeline by declaring it empty
  pipeline "build" {}

  # pipeline specific to the "development workflow"
  pipeline "test" {
   task "npm-test" {}
  }
}
```

This defines one global [pipeline] "build" which is then used in the "development" 
workflow which also defines a child [pipeline] "test" and a Docker image to use
as the environment on which to execute the pipelines.

## Execution
A workflow will execute its pipelines in the order that they are written unless
the pipeline . if a
pipeline fails then the workflow would be considered "failed" and would stop 
execution.

## Workflows and Docker Images
A workflow can define an image that it will use as the base for pipeline 
execution this image can however be overriden by the pipeline if necessary.
Styx will use the same syntax as docker for retrieving images. 
