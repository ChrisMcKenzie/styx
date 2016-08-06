# Scheduler (07/02/2016)

Changes:
  - (07/02/2016): Created

In order to execute Workflows in parallel we will need to distribute work across
multiple machine, but we would like to do it in a smart instead of just round
robin distributing the work. We would like to distribute jobs based on the "load"
of each node ensure that Workflows are performed as quickly as possible.
