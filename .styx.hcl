variable "my-template" {
  default = "${var.app_name}-${var.version}-${var.release}"
}

task "go-get" {
  script = <<EOF
# do some bash stuff here
EOF
}


pipeline "build" {
  task "npm-install" {
    script = "npm install"
  }
}

pipeline "test" {
  task "npm-test" {
    script = "npm test"
  }
}

workflow "development" {
  image = "node:latest"

  driver "local" {
    home = "/tmp"
  }

  task "cool-test" {
    script = "hello"
  }

  pipeline "build" {}
  pipeline "test" {}

  pipeline "deploy" {
    task "cool-test" {}
    task "go-get" {}
    task "go-get" {
      script = "this is a pipeline local task"
    }
  }
}
