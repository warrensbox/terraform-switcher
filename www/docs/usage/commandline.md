## Use dropdown menu to select version
![tfswitch](../static/tfswitch.gif "tfswitch")

1.  You can switch between different versions of terraform by typing the command `tfswitch` on your terminal.
2.  Select the version of terraform you require by using the up and down arrow.
3.  Hit **Enter** to select the desired version.

The most recently selected versions are presented at the top of the dropdown.

## Supply version on command line
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v4.gif" alt="drawing" style="width: 600px;"/>

1. You can also supply the desired version as an argument on the command line.
2. For example, `tfswitch 0.10.5` for version 0.10.5 of terraform.
3. Hit **Enter** to switch.

## See all versions including beta, alpha and release candidates(rc)
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v5.gif" alt="drawing" style="width: 600px;"/>

1. Display all versions including beta, alpha and release candidates(rc).
2. For example, `tfswitch -l` or `tfswitch --list-all` to see all versions.
3. Hit **Enter** to select the desired version.

## Use environment variables

You can also set environment variables for tfswitch to override some configurations:

### `TF_VERSION`

`TF_VERSION` environment variable can be set to your desired terraform version.

For example:

```bash
export TF_VERSION=0.14.4
tfswitch # Will automatically switch to terraform version 0.14.4
```

### `TF_DEFAULT_VERSION`

`TF_DEFAULT_VERSION` environment variable can be set to your desired terraform version that will be used as a fallback version, if not other sources are found.

For example:

```bash
export TF_DEFAULT_VERSION=0.14.4
tfswitch # Will automatically switch to terraform version 0.14.4
```

### `TF_PRODUCT`

`TF_PRODUCT` environment variable can be set to your desired product/tool.

This can either be set to:

 * `terraform`
 * `opentofu`

For example:

```bash
export TF_PRODUCT=opentofu
tfswitch # Will install opentofu instead of terraform
```

### `TF_ARCH`

`TF_ARCH` environment variable can be set to override default CPU architecture for downloaded Terraform binary.

- This can be set to any string, though incorrect values will result in download failure.
- Suggested values: `amd64`, `arm64`, `386`.
- Check available Arch types at:
  - [Terraform Downloads](https://releases.hashicorp.com/terraform/)
  - [OpenTofu Downloads](https://get.opentofu.org/tofu/)

For example:

```bash
export TF_ARCH=amd64
tfswitch # Will install Terraform binary for amd64 architecture
```

## Install latest version only

1. Install the latest stable version only.
2. Run `tfswitch -u` or `tfswitch --latest`.
3. Hit **Enter** to install.

## Install latest implicit version for stable releases

1. Install the latest implicit stable version.
2. Ex: `tfswitch -s 0.13` or `tfswitch --latest-stable 0.13` downloads 0.13.6 (latest) version.
3. Hit **Enter** to install.

## Install latest implicit version for beta, alpha and release candidates(rc)

1. Install the latest implicit pre-release version.
2. Ex: `tfswitch -p 0.13` or `tfswitch --latest-pre 0.13` downloads 0.13.0-rc1 (latest) version.
3. Hit **Enter** to install.

## Show latest version only

1. Just show what the latest version is.
2. Run `tfswitch -U` or `tfswitch --show-latest`
3. Hit **Enter** to show.

## Show latest implicit version for stable releases

1. Show the latest implicit stable version.
2. Ex: `tfswitch -S 0.13` or `tfswitch --show-latest-stable 0.13` shows 0.13.6 (latest) version.
3. Hit **Enter** to show.

## Show latest implicit version for beta, alpha and release candidates(rc)

1. Show the latest implicit pre-release version.
2. Ex: `tfswitch -P 0.13` or `tfswitch --show-latest-pre 0.13` shows 0.13.0-rc1 (latest) version.
3. Hit **Enter** to show.

## Use custom mirror

To install from a remote mirror other than the default(https://releases.hashicorp.com/terraform). Use the `-m` or `--mirror` parameter.

```bash
tfswitch --mirror https://example.jfrog.io/artifactory/hashicorp`
```

## Install to non-default location

By default `tfswitch` will download the Terraform binary to the user home directory under this path: `$HOME/.terraform.versions`

If you want to install the binaries outside of the home directory then you can provide the `-i` or `--install` to install Terraform binaries to a non-standard path. Useful if you want to install versions of Terraform that can be shared with multiple users.

The Terraform binaries will then be placed in the directory `.terraform.versions` under the custom install path e.g. `/opt/terraform/.terraform.versions`

```bash
tfswitch -i /opt/terraform
```

**NOTE**: The directory passed in `-i`/`--install` must be created before running `tfswitch`

## Install binary for non-default architecture

By default `tfswitch` will download the binary for the architecture of the host machine.

If you want to download the binary for non-default CPU architecture then you can provide the `-A` or `--arch` command line argument to download binaries for custom CPU architecture. Useful if you need to override binary architecture for whatever reason.

```bash
tfswitch --arch amd64
```

**NOTE**: If the target file already exists in the download directory (See [Install to non-default location](#install-to-non-default-location) section above), it will be not downloaded. Downloaded files are stored without the architecture in the filename. Format of the filenames in download directory: `<product>_<version>`. E.g. `terraform_1.10.4`.
