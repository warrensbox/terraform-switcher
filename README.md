[![Build Status](https://travis-ci.org/warrensbox/terraform-switcher.svg?branch=master)](https://travis-ci.org/warrensbox/terraform-switcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/warrensbox/terraform-switcher)](https://goreportcard.com/report/github.com/warrensbox/terraform-switcher)
[![CircleCI](https://circleci.com/gh/warrensbox/terraform-switcher/tree/master.svg?style=shield&circle-token=55ddceec95ff67eb38269152282f8a7d761c79a5)](https://circleci.com/gh/warrensbox/terraform-switcher)

# Terraform Switcher

<img style="text-allign:center" src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/smallerlogo.png" alt="drawing" width="120" height="130"/>

<!-- ![gopher](https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/logo.png =100x20) -->

The `tfswitch` command line tool lets you switch between different versions of [terraform](https://www.terraform.io/).
If you do not have a particular version of terraform installed, `tfswitch` will download the version you desire.
The installation is minimal and easy.
Once installed, simply select the version you require from the dropdown and start using terraform.

## Installation

`tfswitch` is available for MacOS and Linux based operating systems.

### Homebrew

Installation for MacOS is the easiest with Homebrew. [If you do not have homebrew installed, click here](https://brew.sh/).


```ruby
brew install warrensbox/tap/tfswitch
```

### General Linux

Installation for other linux operation systems.

```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh | bash
```

### Arch User Repository (AUR) packages for Arch Linux

```sh
# compiled from source
yay tfswitch
# precompiled
yay tfswitch-bin
```
### Install from source

Alternatively, you can install the binary from source [here](https://github.com/warrensbox/terraform-switcher/releases)

[Having trouble installing](https://tfswitch.warrensbox.com/Troubleshoot/).

## How to use:
### Use dropdown menu to select version
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch.gif#1" alt="drawing" style="width: 370px;"/>

1.  You can switch between different versions of terraform by typing the command `tfswitch` on your terminal.
2.  Select the version of terraform you require by using the up and down arrow.
3.  Hit **Enter** to select the desired version.

The most recently selected versions are presented at the top of the dropdown.

### Supply version on command line
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v4.gif#1" alt="drawing" style="width: 370px;"/>

1. You can also supply the desired version as an argument on the command line.
2. For example, `tfswitch 0.10.5` for version 0.10.5 of terraform.
3. Hit **Enter** to switch.

### See all versions including beta, alpha and release candidates(rc)
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v5.gif#1" alt="drawing" style="width: 370px;"/>

1. Display all versions including beta, alpha and release candidates(rc).
2. For example, `tfswitch -l` or `tfswitch --list-all` to see all versions.
3. Hit **Enter** to select the desired version.

### Use environment variable
You can also set the `TF_VERSION` environment variable to your desired terraform version.
For example:   
```bash
export TF_VERSION=0.14.4
tfswitch #will automatically switch to terraform version 0.14.4
```
### Install latest version only
1. Install the latest stable version only.
2. Run `tfswitch -u` or `tfswitch --latest`.
3. Hit **Enter** to install.
### Install latest implicit version for stable releases
1. Install the latest implicit stable version.
2. Ex: `tfswitch -s 0.13` or `tfswitch --latest-stable 0.13` downloads 0.13.6 (latest) version.
3. Hit **Enter** to install.
### Install latest implicit version for beta, alpha and release candidates(rc)
1. Install the latest implicit pre-release version.
2. Ex: `tfswitch -p 0.13` or `tfswitch --latest-pre 0.13` downloads 0.13.0-rc1 (latest) version.
3. Hit **Enter** to install.
### Show latest version only
1. Just show what the latest version is.
2. Run `tfswitch -U` or `tfswitch --show-latest`
3. Hit **Enter** to show.
### Show latest implicit version for stable releases
1. Show the latest implicit stable version.
2. Ex: `tfswitch -S 0.13` or `tfswitch --show-latest-stable 0.13` shows 0.13.6 (latest) version.
3. Hit **Enter** to show.
### Show latest implicit version for beta, alpha and release candidates(rc)
1. Show the latest implicit pre-release version.
2. Ex: `tfswitch -P 0.13` or `tfswitch --show-latest-pre 0.13` shows 0.13.0-rc1 (latest) version.
3. Hit **Enter** to show.
### Use version.tf file
If a .tf file with the terraform constrain is included in the current directory, it should automatically download or switch to that terraform version. For example, the following should automatically switch terraform to the latest version:
```ruby
terraform {
  required_version = ">= 0.12.9"

  required_providers {
    aws        = ">= 2.52.0"
    kubernetes = ">= 1.11.1"
  }
}
```
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/versiontf.gif#1" alt="drawing" style="width: 370px;"/>

### Use .tfswitch.toml file  (For non-admin - users with limited privilege on their computers)
This is similiar to using a .tfswitchrc file, but you can specify a custom binary path for your terraform installation

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v7.gif#1" alt="drawing" style="width: 370px;"/>

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v8.gif#1" alt="drawing" style="width: 370px;"/>

1. Create a custom binary path. Ex: `mkdir /Users/warrenveerasingam/bin` (replace warrenveerasingam with your username)
2. Add the path to your PATH. Ex: `export PATH=$PATH:/Users/warrenveerasingam/bin` (add this to your bash profile or zsh profile)
3. Pass -b or --bin parameter with your custom path to install terraform. Ex: `tfswitch -b /Users/warrenveerasingam/bin/terraform 0.10.8 `
4. Optionally, you can create a `.tfswitch.toml` file in your terraform directory(current directory) OR in your home directory(~/.tfswitch.toml). The toml file in the current directory has a higher precedence than toml file in the home directory
5. Your `.tfswitch.toml` file should look like this:
```ruby
bin = "/Users/warrenveerasingam/bin/terraform"
version = "0.11.3"
```
4. Run `tfswitch` and it should automatically install the required terraform version in the specified binary path

**NOTE** 
1. For linux users that do not have write permission to `/usr/local/bin/`, `tfswitch` will attempt to install terraform at `$HOME/bin`. Run `export PATH=$PATH:$HOME/bin` to append bin to PATH  
2. For windows host, `tfswitch` need to be run under `Administrator` mode, and `$HOME/.tfswitch.toml` with `bin` must be defined (with a valid path) as minimum, below is an example for `$HOME/.tfswitch.toml` on windows

```toml
bin = "C:\\Users\\<%USRNAME%>\\bin\\terraform.exe"
```
### Use .tfswitchrc file
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v6.gif#1" alt="drawing" style="width: 370px;"/>

1. Create a `.tfswitchrc` file containing the desired version
2. For example, `echo "0.10.5" >> .tfswitchrc` for version 0.10.5 of terraform
3. Run the command `tfswitch` in the same directory as your `.tfswitchrc`

#### *Instead of a `.tfswitchrc` file, a `.terraform-version` file may be used for compatibility with [`tfenv`](https://github.com/tfutils/tfenv#terraform-version-file) and other tools which use it*

### Use terragrunt.hcl file
If a terragrunt.hcl file with the terraform constrain is included in the current directory, it should automatically download or switch to that terraform version. For example, the following should automatically switch terraform to the latest version 0.13:
```ruby
terragrunt_version_constraint = ">= 0.26, < 0.27"
terraform_version_constraint  = ">= 0.13, < 0.14"
...
```

### Get the version from a subdirectory
```bash
tfswitch --chdir terraform_dir
tfswitch -c terraform_dir
```

### Use custom mirror 
To install from a remote mirror other than the default(https://releases.hashicorp.com/terraform). Use the `-m` or `--mirror` parameter.    
Ex: `tfswitch --mirror https://example.jfrog.io/artifactory/hashicorp`

## Automation
**Automatically switch with bash**

Add the following to the end of your `~/.bashrc` file:
(Use either `.tfswitchrc` or `.tfswitch.toml` or `.terraform-version`)

```sh
cdtfswitch(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -e "$cdir/.tfswitchrc" ]; then
    tfswitch
  fi
}
alias cd='cdtfswitch'
```

**Automatically switch with zsh**

Add the following to the end of your `~/.zshrc` file:

```sh
load-tfswitch() {
  local tfswitchrc_path=".tfswitchrc"

  if [ -f "$tfswitchrc_path" ]; then
    tfswitch
  fi
}
add-zsh-hook chpwd load-tfswitch
load-tfswitch
```
> NOTE: if you see an error like this: `command not found: add-zsh-hook`, then you might be on an older version of zsh (see below), or you simply need to load `add-zsh-hook` by adding this to your `.zshrc`:
>    ```
>    autoload -U add-zsh-hook
>    ```

*older version of zsh*
```sh
cd(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -e "$cdir/.tfswitchrc" ]; then
    tfswitch
  fi
}
```

**Automatically switch with fish shell**

Add the following to the end of your `~/.config/fish/config.fish` file:

```sh
function switch_terraform --on-event fish_postexec
    string match --regex '^cd\s' "$argv" > /dev/null
    set --local is_command_cd $status

    if test $is_command_cd -eq 0 
      if count *.tf > /dev/null

        grep -c "required_version" *.tf > /dev/null
        set --local tf_contains_version $status
        
        if test $tf_contains_version -eq 0      
            command tfswitch
        end
      end
    end
end
```

### Jenkins setup
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/jenkins_tfswitch.png" alt="drawing" style="width: 170px;"/>

```sh
#!/bin/bash 

echo "Installing tfswitch locally"
wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh #Get the installer on to your machine

chmod 755 install.sh            #Make installer executable

./install.sh -b `pwd`/.bin      #Install tfswitch in a location you have permission

CUSTOMBIN=`pwd`/.bin            #set custom bin path

export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

$CUSTOMBIN/tfswitch -b $CUSTOMBIN/terraform 0.11.7 #or simply tfswitch -b $CUSTOMBIN/terraform 0.11.7

terraform -v                    #testing version
```

### Circle CI setup

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/circleci_tfswitch.png" alt="drawing" style="width: 280px;"/>


Example config yaml
```yaml
version: 2
jobs:
  build:
    docker:
      - image: ubuntu

    working_directory: /go/src/github.com/warrensbox/terraform-switcher

    steps:
      - checkout
      - run: 
          command: |    
            set +e   
            apt-get update 
            apt-get install -y wget 
            rm -rf /var/lib/apt/lists/*

            echo "Installing tfswitch locally"

            wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh #Get the installer on to your machine

            chmod 755 install.sh            #Make installer executable

            ./install.sh -b `pwd`/.bin      #Install tfswitch in a location you have permission

            CUSTOMBIN=`pwd`/.bin            #set custom bin path

            export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

            $CUSTOMBIN/tfswitch -b $CUSTOMBIN/terraform 0.11.7 #or simply tfswitch -b $CUSTOMBIN/terraform 0.11.7

            terraform -v                    #testing version
```
## Order of precedence

| Order | Method |
| --- | ----------- |
| 1 | .tfswitch.toml |
| 2 | .tfswitchrc |
| 3 | .terraform-version |
| 4 | Environment variable |

With 1 being the highest precedence and 4 the lowest   
*(If you disagree with this order of precedence, please open an issue)*
## How to contribute    
An open source project becomes meaningful when people collaborate to improve the code.    
Feel free to look at the code, critique and make suggestions. Lets make `tfswitch` better!   

See step-by-step instructions on how to contribute here: [Contribute](https://tfswitch.warrensbox.com/How-to-Contribute/)      

## Additional Info

See how to *upgrade*, *uninstall*, *troubleshoot* here: [More info](https://tfswitch.warrensbox.com/Upgrade-or-Uninstall/)   


## Issues

Please open  *issues* here: [New Issue](https://github.com/warrensbox/terraform-switcher/issues)
