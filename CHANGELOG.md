<!-- markdownlint-disable MD013 MD024 MD043 -->
<!-- textlint-disable -->

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://semver.org).

## [v1.11.0](https://github.com/warrensbox/terraform-switcher/tree/v1.11.0) - 2025-12-13

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.10.0...v1.11.0)

### Fixed

- fix: Fail fast if `chdir` directory is not readable [#671](https://github.com/warrensbox/terraform-switcher/pull/671) ([yermulnik](https://github.com/yermulnik))
- fix: Don't PromptUI in non-interactive terminal [#669](https://github.com/warrensbox/terraform-switcher/pull/669) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.11.0` [#672](https://github.com/warrensbox/terraform-switcher/pull/672) ([yermulnik](https://github.com/yermulnik))

## [v1.10.0](https://github.com/warrensbox/terraform-switcher/tree/v1.10.0) - 2025-11-26

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.9.0...v1.10.0)

### Fixed

- fix(symlink): Improve symlinking [#648](https://github.com/warrensbox/terraform-switcher/pull/648) ([yermulnik](https://github.com/yermulnik))

### Other

- go: Bump Go version to 1.25 :warning: [#653](https://github.com/warrensbox/terraform-switcher/pull/653) ([yermulnik](https://github.com/yermulnik))
- docs: Update CHANGELOG with `v1.10.0` [#661](https://github.com/warrensbox/terraform-switcher/pull/661) ([yermulnik](https://github.com/yermulnik))

## [v1.9.0](https://github.com/warrensbox/terraform-switcher/tree/v1.9.0) - 2025-10-31

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.8.0...v1.9.0)

### Added

- feat(terragrunt): Support `root.hcl` and custom file name [#644](https://github.com/warrensbox/terraform-switcher/pull/644) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.9.0` [#647](https://github.com/warrensbox/terraform-switcher/pull/647) ([yermulnik](https://github.com/yermulnik))

## [v1.8.0](https://github.com/warrensbox/terraform-switcher/tree/v1.8.0) - 2025-10-23

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.7.0...v1.8.0)

### Added

- feat: Check if version matches version requirement [#640](https://github.com/warrensbox/terraform-switcher/pull/640) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.8.0` [#643](https://github.com/warrensbox/terraform-switcher/pull/643) ([yermulnik](https://github.com/yermulnik))

## [v1.7.0](https://github.com/warrensbox/terraform-switcher/tree/v1.7.0) - 2025-09-30

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.6.0...v1.7.0)

### Added

- feat: Add «Show Required» flag [#631](https://github.com/warrensbox/terraform-switcher/pull/631) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.7.0` [#634](https://github.com/warrensbox/terraform-switcher/pull/634) ([yermulnik](https://github.com/yermulnik))

## [v1.6.0](https://github.com/warrensbox/terraform-switcher/tree/v1.6.0) - 2025-09-07

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.5.1...v1.6.0)

### Added

- feat: Add Keybase as fallback source of Hashicorp's public PGP-key [#624](https://github.com/warrensbox/terraform-switcher/pull/624) ([yermulnik](https://github.com/yermulnik))

### Fixed

- fix: [release workflow] `mkdocs gh-deploy` requires creds to persist [#623](https://github.com/warrensbox/terraform-switcher/pull/623) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.6.0` [#625](https://github.com/warrensbox/terraform-switcher/pull/625) ([yermulnik](https://github.com/yermulnik))

## [v1.5.1](https://github.com/warrensbox/terraform-switcher/tree/v1.5.1) - 2025-09-06

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.5.0...v1.5.1)

### Other

- chore: DependaBot updates
- docs: Update CHANGELOG with `v1.5.1` [#622](https://github.com/warrensbox/terraform-switcher/pull/622) ([yermulnik](https://github.com/yermulnik))

## [v1.5.0](https://github.com/warrensbox/terraform-switcher/tree/v1.5.0) - 2025-07-30

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.4.7...v1.5.0)

### Fixed

- fix: Rectify incomplete/incorrect logging library logic [#602](https://github.com/warrensbox/terraform-switcher/pull/602) ([yermulnik](https://github.com/yermulnik))
- fix(Goreleaser): Do not exclude `"^.*?test(ing)?"` from log (overlaps with word `latest`) [#605](https://github.com/warrensbox/terraform-switcher/pull/605) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.5.0` [#608](https://github.com/warrensbox/terraform-switcher/pull/608) ([yermulnik](https://github.com/yermulnik))

## [v1.4.7](https://github.com/warrensbox/terraform-switcher/tree/v1.4.7) - 2025-07-24

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.4.6...v1.4.7)

### Fixed

- fix: Command-line flags to show/install latest stable versions are broken [#601](https://github.com/warrensbox/terraform-switcher/pull/601) ([yermulnik](https://github.com/yermulnik))

### Other

- chore(`Makefile`): another small set of improvements (no functional changes) [#600](https://github.com/warrensbox/terraform-switcher/pull/600) ([yermulnik](https://github.com/yermulnik))
- chore(`Makefile`): Improve `super-linter` target for a plain repo and a worktree [#597](https://github.com/warrensbox/terraform-switcher/pull/597) ([yermulnik](https://github.com/yermulnik))
- docs: Update CHANGELOG with `v1.4.7` [#604](https://github.com/warrensbox/terraform-switcher/pull/604) ([yermulnik](https://github.com/yermulnik))

## [v1.4.6](https://github.com/warrensbox/terraform-switcher/tree/v1.4.6) - 2025-07-03

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.4.5...v1.4.6)

### Added

- feat: Togglable color logging [#594](https://github.com/warrensbox/terraform-switcher/pull/594) ([yermulnik](https://github.com/yermulnik))
- feat(shell-completion): Add Bash completion script [#586](https://github.com/warrensbox/terraform-switcher/pull/586) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.4.6` [#596](https://github.com/warrensbox/terraform-switcher/pull/596) ([yermulnik](https://github.com/yermulnik))
- chore(`Makefile`): Add `govulncheck` target [#593](https://github.com/warrensbox/terraform-switcher/pull/593) ([yermulnik](https://github.com/yermulnik))
- chore(`.github/linters/.golangci.yml`): Sync config with super-linter [#592](https://github.com/warrensbox/terraform-switcher/pull/592) ([yermulnik](https://github.com/yermulnik))
- chore(goreleaser): Brews: install Bash completion script [#590](https://github.com/warrensbox/terraform-switcher/pull/590) ([yermulnik](https://github.com/yermulnik))

## [v1.4.5](https://github.com/warrensbox/terraform-switcher/tree/v1.4.5) - 2025-05-05

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.4.4...v1.4.5)

### Changed

- feat(gpg): Switch from deprecated Go `crypto` to `ProtonMail/gopenpgp@v3` [#579](https://github.com/warrensbox/terraform-switcher/pull/579) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.4.5` [#580](https://github.com/warrensbox/terraform-switcher/pull/580) ([yermulnik](https://github.com/yermulnik))

## [v1.4.4](https://github.com/warrensbox/terraform-switcher/tree/v1.4.4) - 2025-04-07

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.4.3...v1.4.4)

### Fixed

- chore: Add `super-linter` [#573](https://github.com/warrensbox/terraform-switcher/pull/573) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.4.4` [#575](https://github.com/warrensbox/terraform-switcher/pull/575) ([yermulnik](https://github.com/yermulnik))

## [v1.4.3](https://github.com/warrensbox/terraform-switcher/tree/v1.4.3) - 2025-03-19

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.4.2...v1.4.3)

### Added

- feat: Add env vars for `install` and `bin` args [#566](https://github.com/warrensbox/terraform-switcher/pull/566) ([yermulnik](https://github.com/yermulnik))
- fix: reinstate env var expansion in `bin` and add the same for `install` TOML params [#570](https://github.com/warrensbox/terraform-switcher/pull/570) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.4.3` [#569](https://github.com/warrensbox/terraform-switcher/pull/569) ([yermulnik](https://github.com/yermulnik))

## [v1.4.2](https://github.com/warrensbox/terraform-switcher/tree/v1.4.2) - 2025-03-14

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.4.1...v1.4.2)

### Added

- fix: Improve logging and fix symlink installation issues [#559](https://github.com/warrensbox/terraform-switcher/pull/559) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.4.2` [#565](https://github.com/warrensbox/terraform-switcher/pull/565) ([yermulnik](https://github.com/yermulnik))

## [v1.4.1](https://github.com/warrensbox/terraform-switcher/tree/v1.4.1) - 2025-03-05

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.4.0...v1.4.1)

### Added

- feat: Implement exclusive locking during download and install [#551](https://github.com/warrensbox/terraform-switcher/pull/551) ([yermulnik](https://github.com/yermulnik))
- feat: recommend updating `$PATH` only when necessary [#549](https://github.com/warrensbox/terraform-switcher/pull/549) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.4.1` [#554](https://github.com/warrensbox/terraform-switcher/pull/554) ([yermulnik](https://github.com/yermulnik))

## [v1.4.0](https://github.com/warrensbox/terraform-switcher/tree/v1.4.0) - 2025-02-26

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.3.0...v1.4.0)

### Added

- feat: Improve messaging [#546](https://github.com/warrensbox/terraform-switcher/pull/546) ([yermulnik](https://github.com/yermulnik))
- fix: Create home bin directory if it does not already exist (reinstate historical behavior) [#544](https://github.com/warrensbox/terraform-switcher/pull/544) ([MatthewJohn](https://github.com/MatthewJohn))
- feat: Override default build format to Zip when building for Windows [#527](https://github.com/warrensbox/terraform-switcher/pull/527) ([felipebraga](https://github.com/felipebraga))

### Other

- :warning: go: bump golang.org/x/crypto from 0.33.0 to 0.34.0 (`go.mod`: upgrade go directive to at least 1.23.0) [#545](https://github.com/warrensbox/terraform-switcher/pull/545) ([dependabot[bot]](https://github.com/apps/dependabot))
- docs: Update CHANGELOG with `v1.4.0` [#548](https://github.com/warrensbox/terraform-switcher/pull/548) ([yermulnik](https://github.com/yermulnik))

## [v1.3.0](https://github.com/warrensbox/terraform-switcher/tree/v1.3.0) - 2025-01-26

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.2.4...v1.3.0)

### Added

- feat: Allow to download binary with custom CPU arch [#532](https://github.com/warrensbox/terraform-switcher/pull/532) ([yermulnik](https://github.com/yermulnik))

### Fixed

- fix(Makefile): bring up to date [#535](https://github.com/warrensbox/terraform-switcher/pull/535) ([yermulnik](https://github.com/yermulnik))
- fix: Exclude CHANGELOG howto from distribution archives [#528](https://github.com/warrensbox/terraform-switcher/pull/528) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Update CHANGELOG with `v1.3.0` [#539](https://github.com/warrensbox/terraform-switcher/pull/539) ([yermulnik](https://github.com/yermulnik))
- docs: TOML file belongs in Home dir only [#534](https://github.com/warrensbox/terraform-switcher/pull/534) ([yermulnik](https://github.com/yermulnik))
- docs: Update CHANGELOG with `v1.2.4` [#520](https://github.com/warrensbox/terraform-switcher/pull/520) ([yermulnik](https://github.com/yermulnik))

## [v1.2.4](https://github.com/warrensbox/terraform-switcher/tree/v1.2.4) - 2024-11-23

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.2.3...v1.2.4)

### Other

- :dependabot: Update Go dependencies ([dependabot[bot]](https://github.com/apps/dependabot))
- docs: Update CHANGELOG with `v1.2.4` [#520](https://github.com/warrensbox/terraform-switcher/pull/520) ([yermulnik](https://github.com/yermulnik))
- docs: Update CHANGELOG with `v1.2.3` [#509](https://github.com/warrensbox/terraform-switcher/pull/509) ([yermulnik](https://github.com/yermulnik))

## [v1.2.3](https://github.com/warrensbox/terraform-switcher/tree/v1.2.3) - 2024-09-30

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.2.2...v1.2.3)

### Added

- Split output and logs [#506](https://github.com/warrensbox/terraform-switcher/pull/506) ([MatthewJohn](https://github.com/MatthewJohn))

### Fixed

- fix: Improve wording when installing binary [#494](https://github.com/warrensbox/terraform-switcher/pull/494) ([yermulnik](https://github.com/yermulnik))

### Other

- docs: Adjust version definition order of precedence [#492](https://github.com/warrensbox/terraform-switcher/pull/492) ([yermulnik](https://github.com/yermulnik))
- docs: Update CHANGELOG with v1.2.1 and v1.2.2 [#487](https://github.com/warrensbox/terraform-switcher/pull/487) ([warrensbox](https://github.com/warrensbox))

## [v1.2.2](https://github.com/warrensbox/terraform-switcher/tree/v1.2.2) - 2024-07-07

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.2.1...v1.2.2)

### Added

- default: Fix default binary symlink name for Opentofu to be tofu.exe, rather than terraform.exe [#483](https://github.com/warrensbox/terraform-switcher/pull/483) ([MatthewJohn](https://github.com/MatthewJohn))

### Other

- chore: Add Multi Labeler workflow [#486](https://github.com/warrensbox/terraform-switcher/pull/486) ([yermulnik](https://github.com/yermulnik))
- docs: Update CHANGELOG with `v1.2.0` and `v1.2.1` [#484](https://github.com/warrensbox/terraform-switcher/pull/484) ([yermulnik](https://github.com/yermulnik))
- fix: Remove duplicate .exe extension added to paths for windows inside symlink, as this is already handled by ConvertExecutableExt [#481](https://github.com/warrensbox/terraform-switcher/pull/481) ([MatthewJohn](https://github.com/MatthewJohn))

## [v1.2.1](https://github.com/warrensbox/terraform-switcher/tree/v1.2.1) - 2024-07-06

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.2.0...v1.2.1)

### Fixed

- fix: Remove duplicate .exe extension added to paths for windows inside symlink, as this is already handled by ConvertExecutableExt [#481](https://github.com/warrensbox/terraform-switcher/pull/481) ([MatthewJohn](https://github.com/MatthewJohn))

## [v1.2.0](https://github.com/warrensbox/terraform-switcher/tree/v1.2.0) - 2024-07-05

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.2.0-alpha2...v1.2.0)

### Fixed

- fix: Rectify typo from #468 (+ go fmt) [#478](https://github.com/warrensbox/terraform-switcher/pull/478) ([yermulnik](https://github.com/yermulnik))
- Fix user bin path on windows [#468](https://github.com/warrensbox/terraform-switcher/pull/468) ([PierreTechoueyres](https://github.com/PierreTechoueyres))

### Other

- docs: Update CHANGELOG with `v1.2.0-alpha2` [#479](https://github.com/warrensbox/terraform-switcher/pull/479) ([yermulnik](https://github.com/yermulnik))

### This release (v1.2.0) containes changes from v1.2.0-alpha2, v1.2.0-alpha1 and v1.2.0-alpha

#### v1.2.0-alpha2

### Added

- Added ERRORlog-level [#447](https://github.com/warrensbox/terraform-switcher/pull/447) ([MatrixCrawler](https://github.com/MatrixCrawler))
- skip any parsing of config files if the -v flag is provided [#446](https://github.com/warrensbox/terraform-switcher/pull/446) ([MatrixCrawler](https://github.com/MatrixCrawler))
- RECENT to JSON Proposal [#437](https://github.com/warrensbox/terraform-switcher/pull/437) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Convert recent file to json [#436](https://github.com/warrensbox/terraform-switcher/pull/436) ([MatthewJohn](https://github.com/MatthewJohn))
- OpenTofu support and JSON recentFile support [#435](https://github.com/warrensbox/terraform-switcher/pull/435) ([MatthewJohn](https://github.com/MatthewJohn))
- Check version sources based on whether verison has been set [#429](https://github.com/warrensbox/terraform-switcher/pull/429) ([MatthewJohn](https://github.com/MatthewJohn))
- 154 feature request option to dry run [#416](https://github.com/warrensbox/terraform-switcher/pull/416) ([MatrixCrawler](https://github.com/MatrixCrawler))

### Fixed

- Fix for #443 - Error unzipping on windows systems [#450](https://github.com/warrensbox/terraform-switcher/pull/450) ([MatrixCrawler](https://github.com/MatrixCrawler))
- fix: Only extract terraform file from Terraform release archive [#433](https://github.com/warrensbox/terraform-switcher/pull/433) ([MatthewJohn](https://github.com/MatthewJohn))
- docs: Fix links to `install.sh` [#432](https://github.com/warrensbox/terraform-switcher/pull/432) ([yermulnik](https://github.com/yermulnik))
- docs: Streamline/fix docu site files and README [#419](https://github.com/warrensbox/terraform-switcher/pull/419) ([yermulnik](https://github.com/yermulnik))

### Other

- Update README.md [#448](https://github.com/warrensbox/terraform-switcher/pull/448) ([MatrixCrawler](https://github.com/MatrixCrawler))
- feat(goreleaser): Enable Release Changelog [#418](https://github.com/warrensbox/terraform-switcher/pull/418) ([yermulnik](https://github.com/yermulnik))
- added changelog [#415](https://github.com/warrensbox/terraform-switcher/pull/415) ([warrensbox](https://github.com/warrensbox))

#### v1.2.0-alpha1

### Added

- Add debug logging when successfully obtaining parameter from location… [#464](https://github.com/warrensbox/terraform-switcher/pull/464) ([MatthewJohn](https://github.com/MatthewJohn))
- Add TOML configuration and environment variable for "default version" [#463](https://github.com/warrensbox/terraform-switcher/pull/463) ([MatthewJohn](https://github.com/MatthewJohn))
- feat: Allow for case-insensitve matching of products [#458](https://github.com/warrensbox/terraform-switcher/pull/458) ([MatthewJohn](https://github.com/MatthewJohn))

### Fixed

- fix: Allow Env vars in TOML `bin` value (re-introduce feature) [#467](https://github.com/warrensbox/terraform-switcher/pull/467) ([yermulnik](https://github.com/yermulnik))
- Fix usage for product as passing opentofu to --product does not install Terraform [#461](https://github.com/warrensbox/terraform-switcher/pull/461) ([MatthewJohn](https://github.com/MatthewJohn))

### Other

- docs: Update CHANGELOG with `v1.2.0-alpha` [#460](https://github.com/warrensbox/terraform-switcher/pull/460) ([yermulnik](https://github.com/yermulnik))
- Add error return values to signatures of public methods to allow migrating to returning errors rather than Fatals in future [#457](https://github.com/warrensbox/terraform-switcher/pull/457) ([MatthewJohn](https://github.com/MatthewJohn))

#### v1.2.0-alpha

### Added

- Added ERRORlog-level [#447](https://github.com/warrensbox/terraform-switcher/pull/447) ([MatrixCrawler](https://github.com/MatrixCrawler))
- skip any parsing of config files if the -v flag is provided [#446](https://github.com/warrensbox/terraform-switcher/pull/446) ([MatrixCrawler](https://github.com/MatrixCrawler))
- RECENT to JSON Proposal [#437](https://github.com/warrensbox/terraform-switcher/pull/437) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Convert recent file to json [#436](https://github.com/warrensbox/terraform-switcher/pull/436) ([MatthewJohn](https://github.com/MatthewJohn))
- OpenTofu support and JSON recentFile support [#435](https://github.com/warrensbox/terraform-switcher/pull/435) ([MatthewJohn](https://github.com/MatthewJohn))
- Check version sources based on whether verison has been set [#429](https://github.com/warrensbox/terraform-switcher/pull/429) ([MatthewJohn](https://github.com/MatthewJohn))
- 154 feature request option to dry run [#416](https://github.com/warrensbox/terraform-switcher/pull/416) ([MatrixCrawler](https://github.com/MatrixCrawler))

### Fixed

- Fix for #443 - Error unzipping on windows systems [#450](https://github.com/warrensbox/terraform-switcher/pull/450) ([MatrixCrawler](https://github.com/MatrixCrawler))
- fix: Only extract terraform file from Terraform release archive [#433](https://github.com/warrensbox/terraform-switcher/pull/433) ([MatthewJohn](https://github.com/MatthewJohn))
- docs: Fix links to `install.sh` [#432](https://github.com/warrensbox/terraform-switcher/pull/432) ([yermulnik](https://github.com/yermulnik))
- docs: Streamline/fix docu site files and README [#419](https://github.com/warrensbox/terraform-switcher/pull/419) ([yermulnik](https://github.com/yermulnik))

### Other

- Update README.md [#448](https://github.com/warrensbox/terraform-switcher/pull/448) ([MatrixCrawler](https://github.com/MatrixCrawler))
- feat(goreleaser): Enable Release Changelog [#418](https://github.com/warrensbox/terraform-switcher/pull/418) ([yermulnik](https://github.com/yermulnik))
- added changelog [#415](https://github.com/warrensbox/terraform-switcher/pull/415) ([warrensbox](https://github.com/warrensbox))

## [v1.2.0-alpha2](https://github.com/warrensbox/terraform-switcher/tree/v1.2.0-alpha2) - 2024-07-05

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.2.0-alpha1...v1.2.0-alpha2)

### Fixed

- Fix user bin path on windows [#468](https://github.com/warrensbox/terraform-switcher/pull/468) ([PierreTechoueyres](https://github.com/PierreTechoueyres))

### Other

- docs: Update CHANGELOG with `v1.2.0-alpha1` [#471](https://github.com/warrensbox/terraform-switcher/pull/471) ([yermulnik](https://github.com/yermulnik))

## [v1.2.0-alpha1](https://github.com/warrensbox/terraform-switcher/tree/v1.2.0-alpha1) - 2024-06-22

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.2.0-alpha...v1.2.0-alpha1)

### Added

- Add debug logging when successfully obtaining parameter from location… [#464](https://github.com/warrensbox/terraform-switcher/pull/464) ([MatthewJohn](https://github.com/MatthewJohn))
- Add TOML configuration and environment variable for "default version" [#463](https://github.com/warrensbox/terraform-switcher/pull/463) ([MatthewJohn](https://github.com/MatthewJohn))
- feat: Allow for case-insensitve matching of products [#458](https://github.com/warrensbox/terraform-switcher/pull/458) ([MatthewJohn](https://github.com/MatthewJohn))

### Fixed

- fix: Allow Env vars in TOML `bin` value (re-introduce feature) [#467](https://github.com/warrensbox/terraform-switcher/pull/467) ([yermulnik](https://github.com/yermulnik))
- Fix usage for product as passing opentofu to --product does not install Terraform [#461](https://github.com/warrensbox/terraform-switcher/pull/461) ([MatthewJohn](https://github.com/MatthewJohn))

### Other

- docs: Update CHANGELOG with `v1.2.0-alpha` [#460](https://github.com/warrensbox/terraform-switcher/pull/460) ([yermulnik](https://github.com/yermulnik))
- Add error return values to signatures of public methods to allow migrating to returning errors rather than Fatals in future [#457](https://github.com/warrensbox/terraform-switcher/pull/457) ([MatthewJohn](https://github.com/MatthewJohn))

## [v1.2.0-alpha](https://github.com/warrensbox/terraform-switcher/tree/v1.2.0-alpha) - 2024-06-08

[Full Changelog](https://github.com/warrensbox/terraform-switcher/compare/v1.1.1...v1.2.0-alpha)

### Added

- Added ERRORlog-level [#447](https://github.com/warrensbox/terraform-switcher/pull/447) ([MatrixCrawler](https://github.com/MatrixCrawler))
- skip any parsing of config files if the -v flag is provided [#446](https://github.com/warrensbox/terraform-switcher/pull/446) ([MatrixCrawler](https://github.com/MatrixCrawler))
- RECENT to JSON Proposal [#437](https://github.com/warrensbox/terraform-switcher/pull/437) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Convert recent file to json [#436](https://github.com/warrensbox/terraform-switcher/pull/436) ([MatthewJohn](https://github.com/MatthewJohn))
- OpenTofu support and JSON recentFile support [#435](https://github.com/warrensbox/terraform-switcher/pull/435) ([MatthewJohn](https://github.com/MatthewJohn))
- Check version sources based on whether verison has been set [#429](https://github.com/warrensbox/terraform-switcher/pull/429) ([MatthewJohn](https://github.com/MatthewJohn))
- 154 feature request option to dry run [#416](https://github.com/warrensbox/terraform-switcher/pull/416) ([MatrixCrawler](https://github.com/MatrixCrawler))

### Fixed

- Fix for #443 - Error unzipping on windows systems [#450](https://github.com/warrensbox/terraform-switcher/pull/450) ([MatrixCrawler](https://github.com/MatrixCrawler))
- fix: Only extract terraform file from Terraform release archive [#433](https://github.com/warrensbox/terraform-switcher/pull/433) ([MatthewJohn](https://github.com/MatthewJohn))
- docs: Fix links to `install.sh` [#432](https://github.com/warrensbox/terraform-switcher/pull/432) ([yermulnik](https://github.com/yermulnik))
- docs: Streamline/fix docu site files and README [#419](https://github.com/warrensbox/terraform-switcher/pull/419) ([yermulnik](https://github.com/yermulnik))

### Other

- Update README.md [#448](https://github.com/warrensbox/terraform-switcher/pull/448) ([MatrixCrawler](https://github.com/MatrixCrawler))
- feat(goreleaser): Enable Release Changelog [#418](https://github.com/warrensbox/terraform-switcher/pull/418) ([yermulnik](https://github.com/yermulnik))
- added changelog [#415](https://github.com/warrensbox/terraform-switcher/pull/415) ([warrensbox](https://github.com/warrensbox))

## [v1.1.1](https://github.com/warrensbox/terraform-switcher/tree/v1.1.1) - 2024-04-27

### Fixed

- Fix issue related to additional configuration apart from terraform_version_constraint failing #401 [#409](https://github.com/warrensbox/terraform-switcher/pull/409) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Fix issue related to terraform version constraint in version.tf not being parsed correctly #410 #402 [#403](https://github.com/warrensbox/terraform-switcher/pull/403) ([MatrixCrawler](https://github.com/MatrixCrawler))
- Fix issue where install.sh is unable to download tfswitch version(s) with 'v' appended to the version number #394 #413 #413 [#403](https://github.com/warrensbox/terraform-switcher/pull/405) ([yermulnik](https://github.com/yermulnik)) and ([d33psky](https://github.com/d33psky))

## [v1.1.0](https://github.com/warrensbox/terraform-switcher/tree/v1.1.0) - 2024-04-25

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
- Fixes issues related to install.sh - #339 [#341](https://github.com/warrensbox/terraform-switcher/pull/341) ([warrensbox](https://github.com/warrensbox))
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

- Upgrade libraries: cve's [#258](https://github.com/warrensbox/terraform-switcher/pull/258) ([warrensbox](https://github.com/warrensbox))
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

- Release 0.13 - Minor release [#245](https://github.com/warrensbox/terraform-switcher/pull/245) ([warrensbox](https://github.com/warrensbox))
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
- Fixes SemVer issue [#209](https://github.com/warrensbox/terraform-switcher/pull/209) ([warrensbox](https://github.com/warrensbox))

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
