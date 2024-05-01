## Order of Terraform version definition precedence

| Order | Method |
| --- | ----------- |
| 1 | `$HOME/.tfswitch.toml` (`version` parameter) |
| 2 | `.tfswitch.toml` (`version` parameter) |
| 3 | `.tfswitchrc` (version as a string) |
| 4 | `.terraform-version` (version as a string) |
| 5 | Terraform root module (`required_version` constraint) |
| 6 | `terragrunt.hcl` (`terraform_version_constraint` parameter) |
| 7 | Environment variable (`TF_VERSION`) |

With 1 being the highest precedence and 7 — the lowest  
*(If you disagree with this order of precedence, please open an issue)*
