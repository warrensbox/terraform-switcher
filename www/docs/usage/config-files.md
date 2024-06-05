## Get the version from a subdirectory

While using the file configuration it might be necessary to change the working directory. You can do that with the `--chdir` or `-c` parameter.

```bash
tfswitch --chdir terraform_dir
tfswitch -c terraform_dir
```

## Use version.tf file

If a .tf file with the terraform constraints is included in the current directory, it should automatically download or switch to that terraform version.  
For example, the following should automatically switch terraform to the lastest version:  

```
terraform {
  required_version = ">= 0.12.9"

  required_providers {
    aws        = ">= 2.52.0"
    kubernetes = ">= 1.11.1"
  }
}
```

![versiontf](../static/versiontf.gif "Use version.tf")

## Use .tfswitchrc file

![tfswitchrc](../static/tfswitch-v6.gif)

1. Create a `.tfswitchrc` file containing the desired version
2. For example, `echo "0.10.5" >> .tfswitchrc` for version 0.10.5 of terraform
3. Run the command `tfswitch` in the same directory as your `.tfswitchrc`

*Instead of a `.tfswitchrc` file, a `.terraform-version` file may be used for compatibility with [`tfenv`](https://github.com/tfutils/tfenv#terraform-version-file) and other tools which use it*

## Use .tfswitch.toml file  (For non-admin - users with limited privilege on their computers)

This is similiar to using a .tfswitchrc file, but you can specify a custom binary path for your terraform installation

![toml1](../static/tfswitch-v7.gif)
![toml2](../static/tfswitch-v8.gif)

1. Create a custom binary path. Ex: `mkdir $HOME/bin`
2. Add the path to your PATH. Ex: `export PATH=$PATH:$HOME/bin` (add this to your bash profile or zsh profile)
3. Pass -b or --bin parameter with your custom path to install terraform. Ex: `tfswitch -b $HOME/bin/terraform 0.10.8 `
4. Optionally, you can create a `.tfswitch.toml` file in your terraform directory(current directory) OR in your home directory(~/.tfswitch.toml). The toml file in the current directory has a higher precedence than toml file in the home directory
5. Your `.tfswitch.toml` file should look like this:

```toml
bin = "$HOME/bin/terraform"
version = "0.11.3"
```

6. Run `tfswitch` and it should automatically install the required terraform version in the specified binary path

**NOTE**

1. For linux users that do not have write permission to `/usr/local/bin/`, `tfswitch` will attempt to install terraform at `$HOME/bin`. Run `export PATH=$PATH:$HOME/bin` to append bin to PATH  
2. For windows host, `tfswitch` need to be run under `Administrator` mode, and `$HOME/.tfswitch.toml` with `bin` must be defined (with a valid path) as minimum, below is an example for `$HOME/.tfswitch.toml` on windows

```toml
bin = "C:\\Users\\<%USRNAME%>\\bin\\terraform.exe"
```

## Setting product using .tfswitch.toml file

The .tfswitch.toml file can be configured with a `product` attribute to configure tfswitch to use Terraform or OpenTofu, by default:

```toml
product = "opentofu"
```

or

```toml
product = "terraform"
```

## Use terragrunt.hcl file

If a terragrunt.hcl file with the terraform constraint is included in the current directory, it should automatically download or switch to that terraform version.  
For example, the following should automatically switch terraform to the lastest version 0.13:

```hcl
terragrunt_version_constraint = ">= 0.26, < 0.27"
terraform_version_constraint  = ">= 0.13, < 0.14"
...
```
