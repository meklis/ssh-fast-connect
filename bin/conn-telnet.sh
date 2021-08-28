#!/usr/bin/expect -f
set host [lindex $argv 0]
set user [lindex $argv 1]
set password [lindex $argv 2]
spawn telnet "$host"
expect ":"
exp_send "$user\r"
expect "ord:"
exp_send "$password\r"
#expect "(.*)"
interact
