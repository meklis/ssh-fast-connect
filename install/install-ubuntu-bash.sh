#!/usr/bin/env bash

# Enable stop script if error
# http://web.archive.org/web/20110314180918/http://www.davidpashley.com/articles/writing-robust-shell-scripts.html
set -e

echo "Install sfc binary"
wget -O ~/.local/bin/sfc https://github.com/meklis/ssh-fast-connect/releases/download/0.1/sfc-linux
chmod +x ~/.local/bin/sfc

echo "Add word complete to bashrc"
cat <<EOF >>  ~/.bashrc
#complete function for sfc
_sfc_complete(){
    COMPREPLY=(`sfc -h`)
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ ${COMP_CWORD} == 1 ]] ; then
        COMPREPLY=( $(compgen -W "`sfc -h`" -- ${cur}) )
        return 0
    fi
}
complete -F _sfc_complete sfc

EOF

cat <<EOF > ~/.sfc.conf.yml
# Profile executed when 'sfc <server name>'
profiles:
  gnome: gnome-terminal  --title='%name%' --tab --active -e "ssh -i ~/.ssh/id_rsa  %username%@%host%"
  gnome-password: gnome-terminal  --title='%name%' --tab --active -e "sshpass -p %password% ssh -o StrictHostKeyChecking=no %username%@%host%"

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
    servers:
      - {name: office.pc, host: 10.0.10.10, username: username, password: password }

EOF

cat <<EOF
sfc v0.1 installed!
Configure you file before using!
Config file - ~/.sfc.conf.yml

Usage: sfc <server name 1> [<server name 2>...]
EOF
set +e