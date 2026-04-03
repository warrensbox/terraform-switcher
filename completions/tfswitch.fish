# Remove once https://github.com/warrensbox/terraform-switcher/issues/537 is implemented
set -l COMMAND tfswitch

complete -c $COMMAND -f

complete -c $COMMAND -s A -l arch                       -d "Override CPU architecture type for downloaded binary" -r -f -a "386 amd64 arm arm64"
complete -c $COMMAND -s b -l bin                        -d "Custom binary path" -r -a "(__fish_complete_directories)"
complete -c $COMMAND -s c -l chdir                      -d "Switch to a different working directory" -r -a "(__fish_complete_directories)"
complete -c $COMMAND -s d -l default                    -d "Default to this version if none detected"
complete -c $COMMAND -s g -l log-level                  -d "Log level" -r -f -a "DEBUG ERROR FATAL INFO NOTICE OFF PANIC TRACE WARN"
complete -c $COMMAND -s h -l help                       -d "Show help"
complete -c $COMMAND -s i -l install                    -d "Custom install path" -r -a "(__fish_complete_directories)"
complete -c $COMMAND -s K -l force-color                -d "Force color output if terminal supports it"
complete -c $COMMAND -s k -l no-color                   -d "Disable color output"
complete -c $COMMAND -s l -l list-all                   -d "List all versions of a product"
complete -c $COMMAND -s m -l mirror                     -d "Install from a remote API other than the default"
complete -c $COMMAND -s n -l match-version-requirement  -d "Check if the requested version matches the requirement mandated by the configuration"
complete -c $COMMAND -s p -l latest-pre                 -d "Latest pre-release implicit version"
complete -c $COMMAND -s P -l show-latest-pre            -d "Show latest pre-release implicit version"
complete -c $COMMAND -s r -l dry-run                    -d "Only show what tfswitch would do"
complete -c $COMMAND -s R -l show-required              -d "Show required (or explicitly requested) version"
complete -c $COMMAND -s s -l latest-stable              -d "Latest implicit version based on a constraint"
complete -c $COMMAND -s S -l show-latest-stable         -d "Show latest implicit version"
complete -c $COMMAND -s t -l product                    -d "Specify which product to use" -r -f -a "opentofu terraform"
complete -c $COMMAND -s u -l latest                     -d "Get latest stable version"
complete -c $COMMAND -s U -l show-latest                -d "Show latest stable version"
complete -c $COMMAND -s v -l version                    -d "Show version"
