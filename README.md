# Styx

A Powerful Workflow Oriented CI/CD Platform.

Styx is a Continuous Integration and Delivery platform aimed at allowing users
to create different workflows for each environment that they manage for their
services.

See [Design Proposal v1](./docs/proposals/README.md)

## How does a workflow work

A workflow is a fundamentally a set of large steps that need to be taken to 
accomplish some work. We call the child tasks of a workflow "Pipelines" and 
they contain the small tasks needed to complete the pipeline, for example you
may have a "production" workflow which contains three pipelines "build", "test",
"deploy"; thos pipeline can then have the "tasks" needed to complete each pipeline.
This allows users to design workflows that contain "pipelines" that depend on
the outcome of other "pipelines" which is the fundamentals of CI/CD. The other 
powerful feature is that pipelines can depend on another to start allow parallel
pipeline execution.
