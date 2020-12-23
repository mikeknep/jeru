terraform {
  required_providers {
    local = {
      version = "2.0.0"
    }
  }

  required_version = "0.14.3"
}

// Change this resource name from "main" to "test" to simulate a refactor that
// requires `terraform state mv` to avoid destroying and recreating the file
resource "local_file" "main" {
  content  = "example"
  filename = "${path.module}/example"
}
