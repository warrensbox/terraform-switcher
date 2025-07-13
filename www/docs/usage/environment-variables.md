<!-- markdownlint-disable MD041 -->

## Environment variables

You can set environment variables for `tfswitch` to override some configurations:

### `FORCE_COLOR`

`tfswitch` defaults to color output if the terminal supports it and if the TTY
is allocated (interactive session).  
`FORCE_COLOR` environment variable can be set to force color output even if the
TTY is **not** allocated (non-interactive session).

- Any non-empty value enables color output.
- Is mutually exclusive with `NO_COLOR` environment variable (see
  [`NO_COLOR`](#no_color)).

### `NO_COLOR`

`tfswitch` defaults to color output if the terminal supports it and if the TTY
is allocated (interactive session).  
`NO_COLOR` environment variable can be set to disable color output forcefully.

- Can be useful in CI/CD pipelines or other non-interactive sessions where ANSI
  color (escape) codes are not desired or are not supported.
- Any non-empty value disables color output.
- Is mutually exclusive with `FORCE_COLOR` environment variable (see
  [`FORCE_COLOR`](#force_color)).

### `TF_ARCH`

`TF_ARCH` environment variable can be set to override default CPU architecture
of downloaded binaries.

- This can be set to any string, though incorrect values will result in
  download failure.
- Suggested values: `amd64`, `arm64`, `386`.
- Check available Arch types at:
  - [Terraform Downloads](https://releases.hashicorp.com/terraform/)
  - [OpenTofu Downloads](https://get.opentofu.org/tofu/)

For example:

```bash
export TF_ARCH="amd64"
tfswitch # Will install binary for amd64 architecture
```

### `TF_BINARY_PATH`

`tfswitch` defaults to install to the `/usr/local/bin/` directory (and falls
back to `$HOME/bin/` otherwise). The target filename is resolved automatically
based on the `product` parameter.  
`TF_BINARY_PATH` environment variable can be set to specify a **full
installation path** (directory + filename). If target directory does not exist,
`tfswitch` falls back to `$HOME/bin/` directory.

For example:

```bash
export TF_BINARY_PATH="$HOME/bin/terraform" # Path to the file
tfswitch # Will install binary as $HOME/bin/terraform
```

### `TF_DEFAULT_VERSION`

`TF_DEFAULT_VERSION` environment variable can be set to the desired product/tool
version that will be used as a fallback version, if not other sources are
found.

For example:

```bash
export TF_DEFAULT_VERSION="0.14.4"
tfswitch # Will automatically switch to terraform version 0.14.4
```

### `TF_INSTALL_PATH`

`tfswitch` defaults to download binaries to the `$HOME/.terraform.versions/`
directory.  
`TF_INSTALL_PATH` environment variable can be set to specify the parent
directory for `.terraform.versions` directory. Current user must have write
permissions to the target directory. If the target directory does not exist,
`tfswitch` will create it.

For example:

```bash
export TF_INSTALL_PATH="/var/cache" # Path to the directory where `.terraform.versions` directory resides
tfswitch # Will download actual binary to /var/cache/.terraform.versions/
```

### `TF_LOG_LEVEL`

`TF_LOG_LEVEL` environment variable can be set to override default log level.

- Supported log levels:
  - `OFF`: disables (suppresses) logging
  - `PANIC`: High severity, unrecoverable errors
  - `FATAL`: Fatal, unrecoverable errors + previous log level
  - `ERROR`: Runtime errors that should definitely be noted + previous log levels
  - `WARN`: Non-critical entries that deserve eyes + previous log levels
  - `INFO`: Default logging level, messages that highlight the progress + previous log levels
  - `NOTICE`: Normal operational entries, but not necessarily noteworthy + previous log levels
  - `DEBUG`: Verbose logging, useful for development + previous log levels
  - `TRACE`: Finer-grained informational events than DEBUG + previous log levels
  - Any other log level value results in error and `tfswitch` will exit with
    non-zero exit code.

For example:

```bash
export TF_LOG_LEVEL="DEBUG"
tfswitch # Will output debug logs
```

### `TF_PRODUCT`

`TF_PRODUCT` environment variable can be set to the desired product/tool.

This can either be set to:

- `terraform`
- `opentofu`

For example:

```bash
export TF_PRODUCT="opentofu"
tfswitch # Will install opentofu instead of terraform
```

### `TF_VERSION`

`TF_VERSION` environment variable can be set to the desired product/tool version.

For example:

```bash
export TF_VERSION="0.14.4"
tfswitch # Will automatically switch to terraform version 0.14.4
```
