# Remove once https://github.com/warrensbox/terraform-switcher/issues/537 is implemented

complete -c tfswitch -f

# --- Flag: -A / --arch ---
complete -c tfswitch -s A -l arch -d "Select architecture" -r -f -a "386 amd64 arm arm64"

# --- Flag: -b / --bin ---
complete -c tfswitch -s b -l bin -d "Custom binary path" -r

# --- Flag: -c / --chdir AND -i / --install ---
complete -c tfswitch -s c -l chdir -d "Switch to a different directory" -r -a "(__fish_complete_directories)"
complete -c tfswitch -s i -l install -d "Install specific version" -r -a "(__fish_complete_directories)"

# --- Flag: -g / --log-level ---
complete -c tfswitch -s g -l log-level -d "Log level" -r -f -a "DEBUG ERROR FATAL INFO NOTICE OFF PANIC TRACE WARN"

# --- Flag: -t / --product ---
complete -c tfswitch -s t -l product -d "Product to switch" -r -f -a "opentofu terraform"

# --- General flags ---
complete -c tfswitch -s l -l list -d "List all installed versions"
complete -c tfswitch -s u -l list-all -d "List all available versions"
complete -c tfswitch -s h -l help -d "Show help"
complete -c tfswitch -s v -l version -d "Show version"
