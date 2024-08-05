/*
 * Removes files like .DS_Store and Thumb.db from the disk, as well as other "useless" files:
 *
 * 	- .DS_Store
 * 	- .cache
 * 	- .gradle
 * 	- .mypy_cache
 * 	- .sass-cache
 * 	- .textpadtmp
 * 	- Thumbs.db
 * 	- __pycache__
 * 	- _build
 * 	- build
 * 	- slprj
 * 	- zig-cache
 * 	- zig-out
 * 	- *.slxc
 * 	- *.bak
 * 	- ~*
 */
package main

import "github.com/acristoffers/tmux-tui/cmd"

func main() {
	cmd.Execute()
}
