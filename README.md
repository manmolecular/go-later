# go-later
Simple terminal "To Do" application that I came up with to simplify my small temporary work tasks.

## Requirements
- Latest Go release (tested on 1.21, most probably won't work below 1.18)

## Usage
1. Clone the repo using `git clone`
2. Build the binary using the following command in the root directory of the cloned repository: 
```shell
make  # which is equal to: "go build -mod vendor -o later ./cmd/later/."
```
3. Make sure that the build process was successful and the binary file `later` exists in the root directory of the repository
4. Validate that the application works - run the `later` binary without any arguments to see available commands:
```shell
./later
```
5. Integrate `later` with your shell of choice; add to your shell configuration file (for `zsh` it is the `~/.zshrc` file) the following lines. **Note**: please, adjust `<YOUR_PATH>` in the snippet below to your own path including the `later` binary. Feel free to modify commands as you wish:
```shell
# later (to do task list) binary and its aliases
# ! NOTE (1): please, make sure that `later` binary is built first
# ! NOTE (2): please, make sure that you set the correct path to the `later` binary (<YOUR_PATH>)
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
6. Restart the terminal, verify that integration with `later` works (macOS iTerm2 example below):
```shell
Last login: Sat Sep 30 01:47:51 on ttys001
Tasks to do: 0 (use "tdl" to see)
➜  ~ td do this later
➜  ~ tdl
1. do this later (created at: 2023-09-30 01:49:12)
➜  ~ tdh
td (add), tdl (list), tdp (pop), tdd (delete), tdc (clean)
```
