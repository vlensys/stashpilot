# stashpilot

A TUI tool for managing git stashes, the part of git that `git stash list` makes miserable.


```
╭────────────────────────────╮╭─────────────────────────────────────────╮
│ stash@{0}       just now   ││ stash@{0}                           12%  │
│ main: WIP: add feature foo ││                                          │
│ stash@{1}       2h ago     ││  file.txt | 2 +-                         │
│ main: WIP: refactor bar    ││  1 file changed, 1 insertion(+), 1 del.. │
│ stash@{2}       3d ago     ││                                          │
│ main: WIP: new go file     ││ diff --git a/file.txt b/file.txt         │
│                            ││ index ...                                │
│                            ││ --- a/file.txt                           │
│                            ││ +++ b/file.txt                           │
│                            ││ @@ -1 +1,2 @@                            │
│                            ││  hello                                   │
│                            ││ +change2                                 │
╰────────────────────────────╯╰─────────────────────────────────────────╯
 stash@{1} applied successfully  ↑/↓ nav  a apply  p pop  d drop  n new  r refresh  q quit
```

## The problem

`git stash list` shows you a flat list with no content. To see what's in a stash you need `git stash show -p stash@{2}`, and to apply the right one without losing it you need to remember the index, type `git stash apply stash@{2}`, then check it worked. You've done this wrong at least once.

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
