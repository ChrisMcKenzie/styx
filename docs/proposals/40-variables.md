# Variable (07/02/2016)

### ___!!! This Proposal is a WIP !!!___

Changes:
  - (07/02/2016): Created

Variables are pretty explainatory for the most part with a few exceptions. Styx
Variables allow users to define norm variable which have defaults and can be
overridden, they can also contain secrets by defining a variable that is linked
to a key in a secure store (eg. Vault) which can then be retrieved only when 
needed.

## Example

```
# prompt for value
var "version" {}

# variable with default
var "release" {
  default = "alpha"
}

# variable with rendered default
var "binary-template" {
  default = "${var.app_name}-${var.version}-${var.release}"
}

# pull secret from vault
var "secret" {
  from = "vault"
  key = "secrets/my-secret"
}
```

