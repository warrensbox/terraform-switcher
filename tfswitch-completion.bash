# Remove once https://github.com/warrensbox/terraform-switcher/issues/537 is implemented
# shellcheck disable=SC2207 # W: Prefer mapfile or read -a to split command output (or quote to avoid splitting).
_tfswitch() {
	local cur prev
	cur=${COMP_WORDS[COMP_CWORD]}
	prev=${COMP_WORDS[COMP_CWORD - 1]}

	if [[ ${cur} == -* ]]; then
		COMPREPLY=($(compgen -W "$(tfswitch --help 2>&1 | grep -Eo '[[:space:]]+(-{1,2}[a-zA-Z0-9-]+)')" -- "$cur"))
		return 0
	fi

	case "${prev}" in
	-A | --arch)
		COMPREPLY=($(compgen -W "386 amd64 arm arm64" -- "$cur"))
		return 0
		;;
	-b | --bin)
		_filedir
		return 0
		;;
	-c | --chdir | -i | --install)
		_filedir -d
		return 0
		;;
	-g | --log-level)
		COMPREPLY=($(compgen -W "DEBUG ERROR INFO NOTICE TRACE" -- "$cur"))
		return 0
		;;
	-t | --product)
		COMPREPLY=($(compgen -W "opentofu terraform" -- "$cur"))
		return 0
		;;
	esac

} && complete -F _tfswitch tfswitch

# vim: filetype=bash shiftwidth=4 tabstop=4 expandtab autoindent:
