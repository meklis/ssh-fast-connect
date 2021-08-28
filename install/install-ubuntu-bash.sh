#!/usr/bin/env bash

# Enable stop script if error
# http://web.archive.org/web/20110314180918/http://www.davidpashley.com/articles/writing-robust-shell-scripts.html
set -e

echo "Install sfc binary"
wget -O ~/.local/bin/fc https://github.com/meklis/ssh-fast-connect/releases/download/0.3/fc-linux
chmod +x ~/.local/bin/fc

echo "Add word complete to bashrc"
cat <<EOF >>  ~/.bashrc
#complete function for fx
_fc_complete(){
    COMPREPLY=(`fc -h`)
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ ${COMP_CWORD} == 1 ]] ; then
        COMPREPLY=( $(compgen -W "`fc -h`" -- ${cur}) )
        return 0
    fi
}
complete -F _fc_complete fc

EOF

cat <<EOF > ~/.fc.conf.yml
# Profile executed when 'fc <server name>'
profiles:
  gnome: gnome-terminal  --title='%name%' --tab --active -e "ssh -i ~/.ssh/id_rsa  %username%@%address%"
  gnome-password: gnome-terminal  --title='%name%' --tab --active -e "sshpass -p %password% ssh -o StrictHostKeyChecking=no %username%@%address%"


# 'Defaults' set in commands by default
# If default parameters not setted in server - they will be set by from defaults
# Parameter 'command' is required for work
# Another parameters is variable and depend on from command
groups:
  - name: General
    defaults:
      profile: gnome
      username: user
      password: password
      ssh_key: ~/.ssh/id_rsa
# Sources must return json or yaml content
# Source can be executable script and return content of yaml or json
#    servers_source: ./servers.yml
    static_servers:
      - {name: office.pc, address: 10.0.10.10, username: username, password: password }
EOF

cat <<EOF
sfc v0.3 installed!
Configure you file before using!
Config file - ~/.fc.conf.yml

Usage: fc <server name 1> [<server name 2>...]
EOF
set +e