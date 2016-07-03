# Pipelines (07/02/2016)

### ___!!! This Proposal is a WIP !!!___

Changes:
  - (07/02/2016): Created

Pipelines allow users to define a set of [Tasks](./30-defined-tasks.md) that
will be executed as a group for things such as building or testing an app. They
extremely flexible and can even be triggered automatically by another pipeline
in the same workflow based on whether that pipeline passed or failed.

## Example

```hcl
workflow "release" {
  pipeline "setup" {
    task "go-get" {}
  }

  pipeline "build" {
    triggers {
      "pipeline:setup": "passed"
    }

    task "go-build" {}
  }
}
```

## Tasks
A pipeline may have any number of tasks and those tasks will be executed in
order. Any task that fails will stop execution of the pipeline and mark it as
failed triggering any pipelines that may need to be executed based on that 
failure.

## Pipeline States

A pipeline can be in the following states throughout its lifecycle:

  - PENDING
  - RUNNING
  - FAILED
  - PASSED

All of these states can be triggers for other pipelines which allows parrallel 
execution scenarios.
