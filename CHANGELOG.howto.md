# How to update CHANGELOG with info on latest release

1. [Install GH CLI](https://github.com/cli/cli?tab=readme-ov-file#installation).
   - [Configure it](https://cli.github.com/manual/#configuration)
1. [Install `gh-changelog`](https://github.com/chelnak/gh-changelog?tab=readme-ov-file#installation-and-usage)
   - Ensure the `.changelog.yml` file is in the root of the repository:
     ```yaml
     ---
     file_name: CHANGELOG.md.new
     excluded_labels:
       - maintenance
       - dependencies
     sections:
       added:
         - enhancement
         - feature
         - new feature
       changed:
         - backwards-incompatible
         - depricated
       fixed:
         - bug
         - bugfix
         - fix
         - fixed
     skip_entries_without_label: false
     show_unreleased: true
     check_for_updates: true
     logger: console
     ```
1. Pull latest data from the `origin`
1. Create new branch and name it accordingly (e.g. `docs/Update_CHANGELOG_with_<latest_version_tag>`).
1. Run `gh changelog new --from-version <previous_version_tag> --next-version <latest_version_tag>` to generate CHANGELOG since the `<previous_version_tag>` to `<latest_version_tag>`.
1. Open `CHANGELOG.md.new`, re-arrange log entries to improve readability if applicable and copy everything under the `The format is based on […]` line (release version(s) with description of changes).
1. Open [`CHANGELOG.md`](CHANGELOG.md) and paste copied data right under `The format is based on […]` line (keep empty line between this line and pasted data).
1. Push your changes using conventional commit messages like ``docs: Update CHANGELOG with `<latest_version_tag>` ``, create PR and have someone from [CODEOWNERS](.github/CODEOWNERS) review and approve it.
