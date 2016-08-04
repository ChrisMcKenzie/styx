variable "binary-template" {
  default = "${var.app_name}-${var.version}-${var.release}"
}

task "go-get" {
  script = <<EOF
  # do some bash stuf here
EOF
}


pipeline "build" {
  task "npm-install" {}
}

pipeline "test" {
  task "npm-test" {}
}

workflow "development" {
  task "cool-test" {
    script = "hello"
  }
  pipeline "build" {}
  pipeline "test" {}

  pipeline "build-failed" {
    requires {
      test = "fail"
    }

    task "local-notification" {
      message = "Local build failed on task: {{.Task}}"
    }
  }
}
