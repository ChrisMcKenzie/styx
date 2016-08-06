# Worker (07/02/2016)

Changes:
  - (07/02/2016): Created

A Worker node is a machine responsible for executing pipelines for workflows 
and then reports back to a central realtime database the status and logs of 
those pipelines.

Workers will report their current load to the scheduler so that it can then 
schedule work on appropriate workers. Workers will also report back pipeline 
logs and status. Workers must also be able to stop a pipeline if the user
requests it.

## Local Execution (Development Scenario)

A user must also be able to execute a workflow pipeline on a users local 
development machine. this means that the cli tool must also have the same 
execution logic.
