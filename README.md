# go-later
The simplest possible "Do it later" that I can come up with to use every day in a terminal for my small tasks.

## Requirements
- Go (tested on 1.21)

## Usage
1. Build the binary first: 
```shell
make
```
2. Add to your shell configuration (for me it is `.zshrc`). I use a set of additional aliases for simplification:
```shell
# later (to do task list) binary and its aliases
export PATH="$HOME/<YOUR_PATH>/go-later:$PATH"
# 'tdh' is a short reminder of available aliases
alias tdh="echo 'td (add), tdl (list), tdp (pop), tdd (delete), tdc (clean)'"
# 'td' is a default alias to add tasks to the list
alias td="later push"
# 'tdl' lists all the saved for later tasks
alias tdl="later list"
# 'tdp' removes the last task from from the list
alias tdp="later pop"
# 'tdd' removes the exact task (by ID) from the list
alias tdd="later delete"
# 'tdc' cleans up the tasks storage (homedir/.later)
alias tdc="later clean"
echo "Tasks to do: $(later count) (use \"tdl\" to see)"
```
3. Restart the shell/terminal, verify that it works:
```shell
Last login: Sat Sep 30 01:47:51 on ttys001
Tasks to do: 0 (use "tdl" to see)
➜  ~ td do this later
➜  ~ tdl
1. do this later (created at: 2023-09-30 01:49:12)
➜  ~ tdh
td (add), tdl (list), tdp (pop), tdd (delete), tdc (clean)
```