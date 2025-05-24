# Remove once https://github.com/warrensbox/terraform-switcher/issues/537 is implemented
# - SC2015 (info): Note that A && B || C is not if-then-else. C may run when A is true.
# - SC2207 W: Prefer mapfile or read -a to split command output (or quote to avoid splitting).
# shellcheck disable=SC2015,SC2207
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
		[[ $(type -t _comp_compgen) == "function" ]] && _comp_compgen -a filedir || _filedir
		return 0
		;;
	-c | --chdir | -i | --install)
		[[ $(type -t _comp_compgen) == "function" ]] && _comp_compgen -a filedir -d || _filedir -d
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

}
complete -F _tfswitch tfswitch

# vim: set filetype=bash shiftwidth=4 tabstop=4 noexpandtab autoindent:
