<!-- markdownlint-disable MD013 MD024 MD043 -->
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://semver.org).

## [v1.1.0](https://github.com/warrensbox/terraform-switcher/tree/v1.1.0) - 2024-04-24

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/1.0.2...v1.1.0)

### Added

- feat: Support multiple TF version constraints from module and improve logging [#362](https://github.com/warrensbox/terraform-switcher/pull/362) ([yermulnik](https://github.com/yermulnik))
- Refactor parameter parsing [#356](https://github.com/warrensbox/terraform-switcher/pull/356) ([MatrixCrawler](https://github.com/MatrixCrawler))
- feat: Build statically linked binaries [#353](https://github.com/warrensbox/terraform-switcher/pull/353) ([yermulnik](https://github.com/yermulnik))
- Logging refactoring [#350](https://github.com/warrensbox/terraform-switcher/pull/350) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Checksum check for TF Binaries [#334](https://github.com/warrensbox/terraform-switcher/pull/334) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Feature: Add flag for install location (optional) [#309](https://github.com/warrensbox/terraform-switcher/pull/309) ([ArronaxKP](https://github.com/ArronaxKP))

### Fixed

- Fix for Version during build process [#374](https://github.com/warrensbox/terraform-switcher/pull/374) ([MatrixCrawler](https://github.com/MatrixCrawler))
- fix for #369 [#370](https://github.com/warrensbox/terraform-switcher/pull/370) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Release Workflow fix [#360](https://github.com/warrensbox/terraform-switcher/pull/360) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Fix/refactor deployment [#352](https://github.com/warrensbox/terraform-switcher/pull/352) ([warrensbox](https://github.com/warrensbox))
- Fix/refactor deployment [#351](https://github.com/warrensbox/terraform-switcher/pull/351) ([warrensbox](https://github.com/warrensbox))

### Documentation

- Added how to install on README [#378](https://github.com/warrensbox/terraform-switcher/pull/378) ([warrensbox](https://github.com/warrensbox))
- Readme and documentation update [#376](https://github.com/warrensbox/terraform-switcher/pull/376) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Update README.md [#364](https://github.com/warrensbox/terraform-switcher/pull/364) ([MatrixCrawler](https://github.com/MatrixCrawler))
- docs: Actualize CHANGELOG [#359](https://github.com/warrensbox/terraform-switcher/pull/359) ([yermulnik](https://github.com/yermulnik))

### Other

- feat(goreleaser): Update `changelog` section of `.goreleaser.yml` [#381](https://github.com/warrensbox/terraform-switcher/pull/381) ([yermulnik](https://github.com/yermulnik))
- Update dependabot.yml [#375](https://github.com/warrensbox/terraform-switcher/pull/375) ([MatrixCrawler](https://github.com/MatrixCrawler))
- optimization suggestion for #372 [#373](https://github.com/warrensbox/terraform-switcher/pull/373) ([MatrixCrawler](https://github.com/MatrixCrawler))
- feat: Add `CODEOWNERS` file [#368](https://github.com/warrensbox/terraform-switcher/pull/368) ([yermulnik](https://github.com/yermulnik))
- go: bump golang.org/x/crypto from 0.17.0 to 0.22.0 [#367](https://github.com/warrensbox/terraform-switcher/pull/367) ([dependabot](https://github.com/dependabot))
- go: bump golang.org/x/crypto from 0.16.0 to 0.17.0 [#366](https://github.com/warrensbox/terraform-switcher/pull/366) ([dependabot](https://github.com/dependabot))
- Create goreport.yml [#365](https://github.com/warrensbox/terraform-switcher/pull/365) ([MatrixCrawler](https://github.com/MatrixCrawler))
- go: bump golang.org/x/sys from 0.18.0 to 0.19.0 [#358](https://github.com/warrensbox/terraform-switcher/pull/358) ([dependabot](https://github.com/dependabot))
- Update Go Package Index [#354](https://github.com/warrensbox/terraform-switcher/pull/354) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Create codeql-analysis.yml [#256](https://github.com/warrensbox/terraform-switcher/pull/256) ([jukie](https://github.com/jukie))

## [1.0.2](https://github.com/warrensbox/terraform-switcher/releases/tag/1.0.2) - 2024-04-01

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/1.0.1...1.0.2)

### Other

- Downgrade ubuntu version [#347](https://github.com/warrensbox/terraform-switcher/pull/347) ([warrensbox](https://github.com/warrensbox))
- fix: Downgrade GH Ubuntu runner to `20.04` [#346](https://github.com/warrensbox/terraform-switcher/pull/346) ([yermulnik](https://github.com/yermulnik))
- Test pipeline - automating release [#343](https://github.com/warrensbox/terraform-switcher/pull/343) ([warrensbox](https://github.com/warrensbox))

## [1.0.1](https://github.com/warrensbox/terraform-switcher/releases/tag/1.0.1) - 2024-04-01

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/1.0.0...1.0.1)

### Other

- fix release - add TAG_CONTEXT [#342](https://github.com/warrensbox/terraform-switcher/pull/342) ([warrensbox](https://github.com/warrensbox))
- Fixes issues related to  install.sh - #339 [#341](https://github.com/warrensbox/terraform-switcher/pull/341) ([warrensbox](https://github.com/warrensbox))
- fix: Attempt to fix PR335 [#340](https://github.com/warrensbox/terraform-switcher/pull/340) ([yermulnik](https://github.com/yermulnik))
- #major release -update [#338](https://github.com/warrensbox/terraform-switcher/pull/338) ([warrensbox](https://github.com/warrensbox))

## [1.0.0](https://github.com/warrensbox/terraform-switcher/releases/tag/1.0.0) - 2024-04-01

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/0.13.1316...1.0.0)

### Fixed

- Update README.md [#305](https://github.com/warrensbox/terraform-switcher/pull/305) ([JayDoubleu](https://github.com/JayDoubleu))

### Other

- #major release -update [#338](https://github.com/warrensbox/terraform-switcher/pull/338) ([warrensbox](https://github.com/warrensbox))
- #major - create major release [#337](https://github.com/warrensbox/terraform-switcher/pull/337) ([warrensbox](https://github.com/warrensbox))
- Testing new pipeline [#336](https://github.com/warrensbox/terraform-switcher/pull/336) ([warrensbox](https://github.com/warrensbox))
- Fix/move circle GitHub ci [#335](https://github.com/warrensbox/terraform-switcher/pull/335) ([warrensbox](https://github.com/warrensbox))
- Use Windows User Directory for bin path [#327](https://github.com/warrensbox/terraform-switcher/pull/327) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Refactor re-use of PLATFORM [#321](https://github.com/warrensbox/terraform-switcher/pull/321) ([eternityduck](https://github.com/eternityduck))
- Use go-homedir.Dir() everywhere [#314](https://github.com/warrensbox/terraform-switcher/pull/314) ([kim0](https://github.com/kim0))
- Make leading slash optional in regex when looking for versions [#313](https://github.com/warrensbox/terraform-switcher/pull/313) ([tusv](https://github.com/tusv))
- Fix init function [#297](https://github.com/warrensbox/terraform-switcher/pull/297) ([jukie](https://github.com/jukie))
- Fix binPath when getting latest [#295](https://github.com/warrensbox/terraform-switcher/pull/295) ([jukie](https://github.com/jukie))
- Updating "x/text" to prevent CVE-2022-32149 [#288](https://github.com/warrensbox/terraform-switcher/pull/288) ([hknerts](https://github.com/hknerts))

## [0.13.1316](https://github.com/warrensbox/terraform-switcher/releases/tag/0.13.1316) - 2024-03-22

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/0.13.1308...0.13.1316)

### Other

- Upgrade GO version - 1.22 [#329](https://github.com/warrensbox/terraform-switcher/pull/329) ([warrensbox](https://github.com/warrensbox))
- Update the go version to 1.22 [#326](https://github.com/warrensbox/terraform-switcher/pull/326) ([surola](https://github.com/surola))

## [0.13.1308](https://github.com/warrensbox/terraform-switcher/releases/tag/0.13.1308) - 2023-02-06

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/0.13.1300...0.13.1308)

### Other

- Feature: Fallback or default version [#287](https://github.com/warrensbox/terraform-switcher/pull/287) ([warrensbox](https://github.com/warrensbox))
- Feature/add fallback option [#286](https://github.com/warrensbox/terraform-switcher/pull/286) ([warrensbox](https://github.com/warrensbox))
- Default version flag added [#275](https://github.com/warrensbox/terraform-switcher/pull/275) ([sivaramsajeev](https://github.com/sivaramsajeev))

## [0.13.1300](https://github.com/warrensbox/terraform-switcher/releases/tag/0.13.1300) - 2022-10-27

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/0.13.1288...0.13.1300)

### Other

- Release [#276](https://github.com/warrensbox/terraform-switcher/pull/276) ([warrensbox](https://github.com/warrensbox))
- Add Linux ARMv{6,7} arch to `install.sh` [#261](https://github.com/warrensbox/terraform-switcher/pull/261) ([yermulnik](https://github.com/yermulnik))
- Setup dependabot [#257](https://github.com/warrensbox/terraform-switcher/pull/257) ([jukie](https://github.com/jukie))

## [0.13.1288](https://github.com/warrensbox/terraform-switcher/releases/tag/0.13.1288) - 2022-07-04

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/0.13.1275...0.13.1288)

### Other

- Upgrade libraries: cve's   [#258](https://github.com/warrensbox/terraform-switcher/pull/258) ([warrensbox](https://github.com/warrensbox))
- Upgrade go libraries to resolve CVE's [#255](https://github.com/warrensbox/terraform-switcher/pull/255) ([jukie](https://github.com/jukie))

## [0.13.1275](https://github.com/warrensbox/terraform-switcher/releases/tag/0.13.1275) - 2022-06-20

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/0.13.1250...0.13.1275)

### Other

- Fix repo - update documentation [#253](https://github.com/warrensbox/terraform-switcher/pull/253) ([warrensbox](https://github.com/warrensbox))
- adjust regex to get all terraform versinos from mirror [#251](https://github.com/warrensbox/terraform-switcher/pull/251) ([chrispruitt](https://github.com/chrispruitt))
- :bug: Corrected bad installations if user doesn't input 2 dots in ver… [#249](https://github.com/warrensbox/terraform-switcher/pull/249) ([afreyermuth98](https://github.com/afreyermuth98))
- [Fix] Fail CircleCI pipeline on errors [#246](https://github.com/warrensbox/terraform-switcher/pull/246) ([yermulnik](https://github.com/yermulnik))

## [0.13.1250](https://github.com/warrensbox/terraform-switcher/releases/tag/0.13.1250) - 2022-05-27

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/0.14.0-alpha-1...0.13.1250)

### Other

- Release 0.13 - Minor release  [#245](https://github.com/warrensbox/terraform-switcher/pull/245) ([warrensbox](https://github.com/warrensbox))
- Fix chdirpath for absolute paths [#244](https://github.com/warrensbox/terraform-switcher/pull/244) ([jukie](https://github.com/jukie))
- Rebase 0.14 master round 3 [#242](https://github.com/warrensbox/terraform-switcher/pull/242) ([jukie](https://github.com/jukie))
- Use '.' vs full git ref to allow forks or other users to get their own tags [#241](https://github.com/warrensbox/terraform-switcher/pull/241) ([jukie](https://github.com/jukie))
- Fix chDirPath option for absolute paths [#240](https://github.com/warrensbox/terraform-switcher/pull/240) ([jukie](https://github.com/jukie))
- Rebase over master [#237](https://github.com/warrensbox/terraform-switcher/pull/237) ([jukie](https://github.com/jukie))
- Upgrade to Golang 1.18 [#235](https://github.com/warrensbox/terraform-switcher/pull/235) ([jukie](https://github.com/jukie))
- Use new api endpoint [#227](https://github.com/warrensbox/terraform-switcher/pull/227) ([jukie](https://github.com/jukie))

## [0.14.0-alpha-1](https://github.com/warrensbox/terraform-switcher/releases/tag/0.14.0-alpha-1) - 2022-05-23

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/0.13.1221...0.14.0-alpha-1)

### Other

- Upgrade to Golang 1.18 [#235](https://github.com/warrensbox/terraform-switcher/pull/235) ([jukie](https://github.com/jukie))
- Automate binary tests from test-data directory [#234](https://github.com/warrensbox/terraform-switcher/pull/234) ([jukie](https://github.com/jukie))
- Upgrade to Go 1.18 [#232](https://github.com/warrensbox/terraform-switcher/pull/232) ([jukie](https://github.com/jukie))
- MacOs -> macOS to be consistent with the standard casing [#230](https://github.com/warrensbox/terraform-switcher/pull/230) ([erikdw](https://github.com/erikdw))

## [0.13.1221](https://github.com/warrensbox/terraform-switcher/releases/tag/0.13.1221) - 2022-05-11

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/0.13.1218...0.13.1221)

### Other

- Fixes Hashicorp link issue [#229](https://github.com/warrensbox/terraform-switcher/pull/229) ([warrensbox](https://github.com/warrensbox))
- Tweaking URL regex based on comment from aderuelle [#225](https://github.com/warrensbox/terraform-switcher/pull/225) ([micahflattbuilt](https://github.com/micahflattbuilt))

## [0.13.1218](https://github.com/warrensbox/terraform-switcher/releases/tag/0.13.1218) - 2022-03-08

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/ee6ab3adac79d07213e5c8e35022c4dd04d68a1a...0.13.1218)

### Fixed

- Fix/semver constraints [#208](https://github.com/warrensbox/terraform-switcher/pull/208) ([warrensbox](https://github.com/warrensbox))

### Other

- Fixes Semver issue and M1 installation issue with homebrew [#216](https://github.com/warrensbox/terraform-switcher/pull/216) ([warrensbox](https://github.com/warrensbox))
- Fixes SemVer issue  [#209](https://github.com/warrensbox/terraform-switcher/pull/209) ([warrensbox](https://github.com/warrensbox))

## [0.13.1201](https://github.com/warrensbox/terraform-switcher/releases/tag/0.13.1201) - 2021-11-28

### Bug fixes

- No matter what users pass to --bin or -b, the local binary is called terraform. User can resume old behavior where -b is custom

### Added

- -c, --chdir=value : Switch to a different working directory before executing the given command. Ex: tfswitch --chdir terraform_project will run tfswitch in the terraform_project directory
- -P, --show-latest-pre=value : Show latest pre-release implicit version. Ex: tfswitch --show-latest-pre 0.13 prints 0.13.0-rc1 (latest)
- -S, --show-latest-stable=value : Show latest implicit version. Ex: tfswitch --show-latest-stable 0.13 prints 0.13.7 (latest)
- -U, --show-latest : Show latest stable version

### Removed

- snapcraft installation option
