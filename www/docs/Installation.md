# Installation

`tfswitch` is available for Windows, macOS and Linux based operating systems.

## Windows

Download and extract the Windows version of `tfswitch` that is compatible with your system.  
We are building binaries for 386, amd64, arm6 and arm7 CPU structure.  
See the [release page](https://github.com/warrensbox/terraform-switcher/releases/latest) for your download.

## Homebrew

Installation for macOS or Linux (where supported) is the easiest with Homebrew. <a href="https://brew.sh/" target="_blank">If you do not have Homebrew installed, click here</a>.

> **⚠** _Homebrew on Linux (formerly referred to as Linuxbrew) requires at least Homebrew [v5.0.6](https://github.com/Homebrew/brew/releases/tag/5.0.6) (released Dec 16, 2025)_

```shell
brew install tfswitch
```

> For users who want the latest release as soon as it is published, you can install `tfswitch` from our Homebrew tap:
>
> `brew install warrensbox/tap/tfswitch`
>
> The official Homebrew formula, `brew install tfswitch`, remains the recommended installation method for most users. The tap is mainly for users who want to get new releases before they are picked up by the main Homebrew repository.

## Linux

Installation for Linux operating systems.

```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/master/install.sh | bash
```

By default installer script will try to download `tfswitch` binary into `/usr/local/bin`  
To install at custom path use below:

```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/master/install.sh | bash -s -- -b $HOME/.local/bin
```

By default installer script will try to download latest version of `tfswitch` binary  
To install custom (not latest) version use:

```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/master/install.sh | bash -s -- 1.1.1
```

Both options can be combined though:

```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/master/install.sh | bash -s -- -b $HOME/.local/bin 1.1.1
```

## Arch User Repository (AUR) packages for Arch Linux

```sh
# compiled from source
yay tfswitch

# precompiled
yay tfswitch-bin
```

## Install from source

Alternatively, you can build from source by cloning the repository (prefer cloning over SSH to have `tfswitch` version defined correctly!) and running `make tfswitch` and have the built binaries put in the `build` directory (run `make install` to build and store built binary in `~/bin/` directory on supported platforms) or otherwise download and extract the pre-compiled binary from <a href="https://github.com/warrensbox/terraform-switcher/releases" target="_blank">releases</a>.

[Having trouble installing](https://tfswitch.warrensbox.com/Troubleshoot/).
