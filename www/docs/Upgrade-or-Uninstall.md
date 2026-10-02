<!-- markdownlint-disable MD041 MD024 -->

## Upgrade

### Homebrew

```shell
brew upgrade tfswitch
```

If installed from the tap, run:
```shell
brew upgrade warrensbox/tap/tfswitch
```

### Linux

Rerun:

```sh
curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/master/install.sh | bash
```

## Uninstall

### Homebrew

```shell
brew uninstall tfswitch
```

If installed from the tap, run:
```shell
brew uninstall warrensbox/tap/tfswitch
```

### Linux

Run (replace `/usr/local/bin` if you installed `tfswitch` to a custom location):

```sh
rm /usr/local/bin/tfswitch
```
