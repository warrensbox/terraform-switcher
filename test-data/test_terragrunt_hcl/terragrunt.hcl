include {
  path = "${find_in_parent_folders()}"
}

terraform {
  source = "..."

  extra_arguments "variables" {
    commands = get_terraform_commands_that_need_vars()
  }
}
inputs = merge(
  jsondecode(file("${find_in_parent_folders("general.tfvars")}"))
)

terraform_version_constraint = ">= 0.13, < 0.14"
