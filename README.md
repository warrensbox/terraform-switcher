[![Build Status](https://travis-ci.org/warrensbox/terraform-switcher.svg?branch=master)](https://travis-ci.org/warrensbox/terraform-switcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/warrensbox/terraform-switcher)](https://goreportcard.com/report/github.com/warrensbox/terraform-switcher)
[![CircleCI](https://circleci.com/gh/warrensbox/terraform-switcher/tree/master.svg?style=shield&circle-token=55ddceec95ff67eb38269152282f8a7d761c79a5)](https://circleci.com/gh/warrensbox/terraform-switcher)

# Terraform Switcher 

<img style="text-allign:center" src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/smallerlogo.png" alt="drawing"/>

<!-- ![gopher](https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/logo.png =100x20) -->

The tfswitch command lets you switch between different terraform versions. 
If you do not have a particular version installed, tfswitch will download the version you desire.
The installation is minimal and easy. 
Simply select the version you require from the dropdown and start using terraform with ease. 

## Installation

At the moment, installation is available for most unix/linux based operating systems.

### Homebrew

Installation for MacOS is the easiest with Homebrew. [If you do not have homebrew installed, click here](https://brew.sh/). 


```ruby
brew install warrensbox/tap/tfswitch
```

To upgrade, simply run `brew upgrade warrensbox/tap/tfswitch`

### Linux

Installation for other linux operation systems.

```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh | bash
```

To upgrade, simply run 

### Install from source

Alternatively, you can install the binary from source [here](https://github.com/warrensbox/terraform-switcher/releases) 

## How to use:

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch.gif" alt="drawing" style="width: 180px;"/>

1.  You can start using by typing the command `tfswitch` on your terminal. 
2.  Select the version of terraform you require by using the up and down arrow.
3.  Hit **Enter** to select the desired version

## Additional Info

[Visit Site](https://warrensbox.github.io/terraform-switcher/)


## Issues

Please open  *issues* here: [New Issue](https://github.com/warrensbox/terraform-switcher/issues)







