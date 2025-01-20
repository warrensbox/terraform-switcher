## Set a default TF version for CI/CD pipeline

1. When using a CI/CD pipeline, you may want a default or fallback version to avoid the pipeline from hanging.
2. Ex: `tfswitch -d 1.2.3` or `tfswitch --default 1.2.3` installs version `1.2.3` when no other versions could be detected.
[Also, see CICD example](../Continuous-Integration.md)

## Automatically switch with bash

Add the following to the end of your `~/.bashrc` file:
(Use either `.tfswitchrc` or `.terraform-version`)

```sh
cdtfswitch(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -e "$cdir/.tfswitchrc" ]; then
    tfswitch
  fi
}
alias cd='cdtfswitch'
```

## Automatically switch with zsh

Add the following to the end of your `~/.zshrc` file:

```sh
load-tfswitch() {
  local tfswitchrc_path=".tfswitchrc"

  if [ -f "$tfswitchrc_path" ]; then
    tfswitch
  fi
}
add-zsh-hook chpwd load-tfswitch
load-tfswitch
```

> NOTE: if you see an error like this: `command not found: add-zsh-hook`, then you might be on an older version of zsh (see below), or you simply need to load `add-zsh-hook` by adding this to your `.zshrc`:
>    ```
>    autoload -U add-zsh-hook
>    ```

### Older version of zsh

```sh
cd(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -e "$cdir/.tfswitchrc" ]; then
    tfswitch
  fi
}
```

## Automatically switch with fish shell

Add the following to the end of your `~/.config/fish/config.fish` file:

```sh
function switch_terraform --on-event fish_postexec
    string match --regex '^cd\s' "$argv" > /dev/null
    set --local is_command_cd $status

    if test $is_command_cd -eq 0
      if count *.tf > /dev/null

        grep -c "required_version" *.tf > /dev/null
        set --local tf_contains_version $status
        
        if test $tf_contains_version -eq 0
            command tfswitch
        end
      end
    end
end
```
