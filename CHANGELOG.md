# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.13.1201] - 2021-11-28
### Bug fixes
- No matter what users pass to --bin or -b, the local binary is called terraform. User can resume old behavior where -b is custom

### Added
-  -c, --chdir=value : Switch to a different working directory before executing the given command. Ex: tfswitch --chdir terraform_project will run tfswitch in the terraform_project directory
- -P, --show-latest-pre=value : Show latest pre-release implicit version. Ex: tfswitch --show-latest-pre 0.13 prints 0.13.0-rc1 (latest)
- -S, --show-latest-stable=value : Show latest implicit version. Ex: tfswitch --show-latest-stable 0.13 prints 0.13.7 (latest)
- -U, --show-latest : Show latest stable version

### Removed
- snapcraft installation option