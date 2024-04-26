# Installation
`tfswitch` is available for Windows, macOS and Linux based operating systems.

## Windows
Download and extract the Windows version of `tfswitch` that is compatible with your system.  
We are building binaries for 386, amd64, arm6 and arm7 CPU structure.  
See the [release page](https://github.com/warrensbox/terraform-switcher/releases/latest) for your download.

## Homebrew
Installation for macOS is the easiest with Homebrew. <a href="https://brew.sh/" target="_blank">If you do not have homebrew installed, click here</a>.

```ruby
brew install warrensbox/tap/tfswitch
```

## Linux
Installation for Linux operating systems.

```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh | bash
```

By default installer script will try to download `tfswitch` binary into `/usr/local/bin`  
To install at custom path use below:
```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh | bash -s -- -b $HOME/.local/bin
```

By default installer script will try to download latest version of `tfswitch` binary  
To install custom (not latest) version use:
```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh | bash -s -- 1.1.1
```

Both options can be combined though:
```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh | bash -s -- -b $HOME/.local/bin 1.1.1
```

## Arch User Repository (AUR) packages for Arch Linux

```sh
# compiled from source
yay tfswitch

# precompiled
yay tfswitch-bin
```

## Install from source

Alternatively, you can install the binary from the source <a href="https://github.com/warrensbox/terraform-switcher/releases" target="_blank">here</a>.

[Having trouble installing](https://tfswitch.warrensbox.com/Troubleshoot/).


