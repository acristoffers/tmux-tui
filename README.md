# tmux-tui

Manages tmux sessions, windows and panes in a Terminal User Interface inspired by lazygit.

![screenshot](.screenshot/screenshot.png)

## Implemented actions

|         | Create | Destroy | Rename| Go to | Swap |
| :---    | :---   | :---    | :---  | :---  | :--- |
| Session | ✓      | ✓       | ✓     | ✓     | ✗    |
| Window  | ✓      | ✓       | ✓     | ✓     | ✓    |
| Pane    | ✗      | ✓       | ✗     | ✓     | ✓    |

Desired, but not yet done: move

## Installation

Install with go:

```bash
go install github.com/acristoffers/tmux-tui@latest
```

or use nix:

```bash
nix profile install github:acristoffers/tmux-tui
```

And add the following line to your `tmux.conf` to open it in a popup window (scratch window):

```tmux
bind-key O display-popup -E -w '80%' -h '80%' tmux-tui
```

You can experiment with the width and height (`-w` and `-h`, respectively).
