<!-- markdownlint-disable MD041 -->

## Use dropdown menu to select version

![tfswitch](../static/tfswitch.gif "tfswitch")

1. You can switch between different versions of terraform by typing the command
   `tfswitch` on your terminal.
2. Select the version of terraform you require by using the up and down arrow.
3. Hit **Enter** to select the desired version.

The most recently selected versions are presented at the top of the dropdown.

## Supply version on command line

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v4.gif" alt="drawing" style="width: 600px;"/>

1. You can also supply the desired version as an argument on the command line.
2. For example, `tfswitch 0.10.5` for version 0.10.5.

## See all versions including beta, alpha and release candidates(rc)

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v5.gif" alt="drawing" style="width: 600px;"/>

1. Display all versions including beta, alpha and release candidates(rc).
2. For example, `tfswitch -l` or `tfswitch --list-all` to see all versions.
3. Hit **Enter** to select the desired version.

## Install latest version only

1. Install the latest stable version only.
2. Run `tfswitch -u` or `tfswitch --latest`

## Install latest implicit version for stable releases

1. Install the latest implicit stable version.
2. Ex: `tfswitch -s 0.13` or `tfswitch --latest-stable 0.13` downloads latest
   on `0.*` branch (`~> 0.13`), while `tfswitch -s 0.13.5` or `tfswitch
--latest-stable 0.13.5` downloads latest on `0.13.*` branch (`~> 0.13.5`) and
   `tfswitch -s 0` or `tfswitch --latest-stable 0` downloads latest on `0` branch
   (`~> 0`).

## Install latest implicit version for beta, alpha and release candidates(rc)

1. Install the latest implicit version, including prereleases versions.
2. Ex: `tfswitch -p 0.13` or `tfswitch --latest-pre 0.13` downloads 0.13.0-rc1
   (latest) version.
3. See examples for `--latest-stable` option above.

## Show latest version only

1. Just show what the latest version is.
2. Run `tfswitch -U` or `tfswitch --show-latest`

## Show latest implicit version for stable releases

1. Show the latest implicit stable version.
2. Ex: `tfswitch -S 0.13` or `tfswitch --show-latest-stable 0.13` shows latest
   on `0.*` branch (`~> 0.13`), while `tfswitch -S 0.13.5` or `tfswitch
--show-latest-stable 0.13.5` shows latest on `0.13.*` branch (`~> 0.13.5`).

## Show latest implicit version for beta, alpha and release candidates(rc)

1. Show the latest implicit version, including prereleases versions.
2. Ex: `tfswitch -P 0.13` or `tfswitch --show-latest-pre 0.13`

## Show required (or explicitly requested) version

1. Show the version required by version constraints.
2. Takes into account version from module version constraint, command line,
   configuration file(s), env var, etc. See [General](general.md) for options.
3. Defaults to latest version if no constraints found.
4. Ex: `tfswitch -R` or `tfswitch --show-required`
5. Can be combined with options like `--latest-stable` and `--latest-pre` to
   show the required version that would be installed by those options.

## Use custom mirror

To install from a remote mirror other than the default
(<https://releases.hashicorp.com/terraform>). Use the `-m` or `--mirror`
parameter.

```bash
tfswitch --mirror https://example.jfrog.io/artifactory/hashicorp`
```

## Install to non-default location

By default `tfswitch` will download the Terraform binary to the user home
directory under this path: `$HOME/.terraform.versions`

If you want to install the binaries outside of the home directory then you can
provide the `-i` or `--install` to install Terraform binaries to a non-standard
path. Useful if you want to install versions of Terraform that can be shared
with multiple users.

The Terraform binaries will then be placed in the directory
`.terraform.versions` under the custom install path e.g.
`/opt/terraform/.terraform.versions`

```bash
tfswitch -i /opt/terraform
```

**NOTE**: The directory passed in `-i`/`--install` must be created before
running `tfswitch`

## Install binary for CPU architecture that doesn't match the host

By default `tfswitch` will download the binary for the CPU architecture of the
host machine.

If you want to download the binary for CPU architecture that doesn't match the
host then you can provide the `-A` or `--arch` command line argument to
download binaries for custom CPU architecture. Useful if you need to override
binary architecture for whatever reason.

```bash
tfswitch --arch amd64
```

**NOTE**: If the target file already exists in the download directory (See
[Install to non-default location](#install-to-non-default-location) section
above), it will be not downloaded. Downloaded files are stored without the
architecture in the filename. Format of the filenames in download directory:
`<product>_<version>`. E.g. `terraform_1.10.4`.

## Disable color output / Force color output

`tfswitch` defaults to color output if the terminal supports it and if the TTY
is allocated (interactive session).

Disabling color output can be useful in non-interactive sessions, such as when
running scripts in CI/CD pipeline or when piping output to other commands.

- If you want to disable color output, you can use the `--no-color` (`-k`) flag.
- If you want to force color output even if the TTY is not allocated
  (non-interactive session), you can use the `--force-color` (`-K`) flag.

**NOTE**: `--no-color` and `--force-color` flags are mutually exclusive.
