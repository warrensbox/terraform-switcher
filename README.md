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

See installation guide here: [tfswitch installation](https://warrensbox.github.io/terraform-switcher/)

## Installation

`tfswitch` is available for MacOS and Linux based operating systems.

### Homebrew

Installation for MacOS is the easiest with Homebrew. [If you do not have homebrew installed, click here](https://brew.sh/).


```ruby
brew install warrensbox/tap/tfswitch
```

### Linux

Installation for other linux operation systems.

```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh | bash
```

### Install from source

Alternatively, you can install the binary from source [here](https://github.com/warrensbox/terraform-switcher/releases)

## How to use:
### Use dropdown menu to select version
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch.gi" alt="drawing" style="width: 370px;"/>

1.  You can switch between different versions of terraform by typing the command `tfswitch` on your terminal.
2.  Select the version of terraform you require by using the up and down arrow.
3.  Hit **Enter** to select the desired version.

The most recently selected versions are presented at the top of the dropdown.

### Supply version on command line
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v4.gi" alt="drawing" style="width: 370px;"/>

1. You can also supply the desired version as an argument on the command line.
2. For example, `tfswitch 0.10.5` for version 0.10.5 of terraform.
3. Hit **Enter** to switch.

### See all versions including beta, alpha and release candidates(rc)
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v5.gi" alt="drawing" style="width: 370px;"/>

1. Display all versions including beta, alpha and release candidates(rc). 
2. For example, `tfswitch -l` or `tfswitch --list-all` to see all versions.
3. Hit **Enter** to select the desired version.

### Use version.tf file  
If a .tf file with the terraform constrain is included in the current directory, it should automatically download or switch to that terraform version. For example, the following should automatically switch terraform to the lastest version:     
```ruby
terraform {
  required_version = ">= 0.12.9"

  required_providers {
    aws        = ">= 2.52.0"
    kubernetes = ">= 1.11.1"
  }
}
```
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/versiontf.gi" alt="drawing" style="width: 370px;"/>


### Use .tfswitch.toml file  (For non-admin - users with limited privilege on their computers)
This is similiar to using a .tfswitchrc file, but you can specify a custom binary path for your terraform installation

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v7.gi" alt="drawing" style="width: 370px;"/>     

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v8.gi" alt="drawing" style="width: 370px;"/>

1. Create a custom binary path. Ex: `mkdir /Users/warrenveerasingam/bin` (replace warrenveerasingam with your username)
2. Add the path to your PATH. Ex: `export PATH=$PATH:/Users/warrenveerasingam/bin` (add this to your bash profile or zsh profile)
3. Pass -b or --bin parameter with your custom path to install terraform. Ex: `tfswitch -b /Users/warrenveerasingam/bin/terraform 0.10.8 `
4. Optionally, you can create a `.tfswitch.toml` file in your terraform directory.
5. Your `.tfswitch.toml` file should look like this:
```
bin = "/Users/warrenveerasingam/bin/terraform"
version = "0.11.3"
```
4. Run `tfswitch` and it should automatically install the required terraform version in the specified binary path

### Use .tfswitchrc file
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v6.gi" alt="drawing" style="width: 370px;"/>

1. Create a `.tfswitchrc` file containing the desired version
2. For example, `echo "0.10.5" >> .tfswitchrc` for version 0.10.5 of terraform
3. Run the command `tfswitch` in the same directory as your `.tfswitchrc`

#### *Instead of a `.tfswitchrc` file, a `.terraform-version` file may be used for compatibility with [`tfenv`](https://github.com/tfutils/tfenv#terraform-version-file) and other tools which use it*

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

### Jenkins setup
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/jenkins_tfswitch.png" alt="drawing" style="width: 170px;"/>

```sh
#!/bin/bash 

echo "Installing tfswitch locally"
wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh 
chmod 755 install.sh
./install.sh -b bin-directory

./bin-directory/tfswitch
```

If you have limited permission, try:

```sh
#!/bin/bash 

echo "Installing tfswitch locally"
wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh 
chmod 755 install.sh
./install.sh -b bin-directory

CUSTOMBIN=`pwd`/bin             #set custom bin path
mkdir $CUSTOMBIN                #create custom bin path
export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

./bin-directory/tfswitch -b $CUSTOMBIN/terraform 0.11.7

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

            wget https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh 
            chmod 755 install.sh
            ./install.sh -b bin-directory

            CUSTOMBIN=`pwd`/bin             #set custom bin path
            mkdir $CUSTOMBIN                #create custom bin path
            export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

            ./bin-directory/tfswitch -b $CUSTOMBIN/terraform 0.11.7

            terraform -v                    #testing version
```

## How to contribute    
An open source project becomes meaningful when people collaborate to improve the code.    
Feel free to look at the code, critique and make suggestions. Lets make `tfswitch` better!   

See step-by-step instructions on how to contribute here: [Contribute](https://tfswitch.warrensbox.com/How-to-Contribute/)      

## Additional Info

See how to *upgrade*, *uninstall*, *troubleshoot* here: [More info](https://tfswitch.warrensbox.com/Upgrade-or-Uninstall/)   


## Issues

Please open  *issues* here: [New Issue](https://github.com/warrensbox/terraform-switcher/issues)
