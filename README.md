[![Build Status](https://github.com/warrensbox/terraform-switcher/actions/workflows/build.yml/badge.svg)](https://github.com/warrensbox/terraform-switcher/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/warrensbox/terraform-switcher)](https://goreportcard.com/report/github.com/warrensbox/terraform-switcher)
![GitHub Release](https://img.shields.io/github/v/release/warrensbox/terraform-switcher)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/warrensbox/terraform-switcher)


# Terraform Switcher

<img style="text-align:center" src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/smallerlogo.png" alt="drawing" width="120" height="130"/>

The `tfswitch` command line tool lets you switch between different versions of [terraform](https://www.terraform.io/).  
If you do not have a particular version of terraform installed, `tfswitch` will download and verify the version you desire.  
The installation is minimal and easy.  
Once installed, simply select the version you require from the dropdown and start using terraform.

## Documentation
Click [here](https://tfswitch.warrensbox.com) for our extended documentation.

## Quick Start
### Dropdown Menu
Just call `tfswitch` and select the desired terraform version via the dropdown menu.
### Version on command line
Call `tfswitch 1.7.0` to install terraform in the version 1.7.0. Replace version string with your desired version.

## NOTE
Going forward we will change the version identifier of `tfswitch` to align with the common go package versioning.  
Please be advised to change any automated implementation you might have that is relying on the `tfswitch` version string.  
**Old version string:** `0.1.2412`  
**New version string:** `v1.0.0` Note the `v` that is preceding all version numbers.


## How to contribute
An open source project becomes meaningful when people collaborate to improve the code.    
Feel free to look at the code, critique and make suggestions. Let's make `tfswitch` better!   

See step-by-step instructions on how to contribute here: [Contribute](https://tfswitch.warrensbox.com/How-to-Contribute/)      

## Additional Info
See how to *upgrade*, *uninstall*, *troubleshoot* here: [More info](https://tfswitch.warrensbox.com/Upgrade-or-Uninstall/)   

## Issues
Please open  *issues* here: [New Issue](https://github.com/warrensbox/terraform-switcher/issues)
