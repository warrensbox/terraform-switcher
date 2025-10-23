<!-- markdownlint-disable MD041 -->

## Order of Terraform version definition precedence

| Order | Method                                                      |
| ----- | ----------------------------------------------------------- |
| 1     | `$HOME/.tfswitch.toml` (`version` parameter)                |
| 2     | `.tfswitchrc` (version as a string)                         |
| 3     | `.terraform-version` (version as a string)                  |
| 4     | Terraform root module (`required_version` constraint)       |
| 5     | `terragrunt.hcl` (`terraform_version_constraint` parameter) |
| 6     | Environment variable (`TF_VERSION`)                         |
| 7     | Version provided as command line argument                   |

With 1 being the **lowest** precedence and 7 — the **highest**  
_(If you disagree with this order of precedence, please open an issue)_
