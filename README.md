# stashpilot

A TUI tool for managing git stashes, the part of git that `git stash list` makes miserable.

<img width="600" alt="image" src="https://github.com/user-attachments/assets/7c48449f-4fb0-401e-a3a9-7ecc9626a50d" />


## Install

```sh
go install github.com/vlensys/stashpilot@latest
```

If `stashpilot` isn't found after installing, add Go's bin directory to your PATH:

```sh
# fish
fish_add_path $HOME/go/bin

# bash
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc && source ~/.bashrc

# zsh
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.zshrc && source ~/.zshrc
```

Or clone and build:

```sh
git clone https://github.com/vlensys/stashpilot
cd stashpilot
go build -o stashpilot .
```

## Usage

Run it from inside any git repo:

```sh
stashpilot
```

| Key | Action |
|-----|--------|
| `↑` / `k` | Navigate up |
| `↓` / `j` | Navigate down |
| `a` | Apply stash (keeps it in the list) |
| `p` | Pop stash (apply + remove) |
| `d` | Drop stash (delete without applying) |
| `n` | Create a new stash from working changes |
| `ctrl+u` | Scroll diff up |
| `ctrl+d` | Scroll diff down |
| `r` | Refresh stash list |
| `q` | Quit |

Destructive actions (`pop`, `drop`) require confirmation with `y` or `Enter`.
