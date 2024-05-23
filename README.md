# bspwm-scratchpad

A scratchpad implementation for bspwm

## Usage
This implementation comes in 2 parts. The daemon, and the interface. In order for your windows to be managed, the daemon must be alive.

### Daemon
The daemon runs as a long running process, boot it up with your preferred init system at launch.

### CLI
The CLI interacts with the daemon over a unix socket (default: `/tmp/scratch.sock`). It is recommended to add calls to the 
CLI in your sxhkd config. There is no 'required calls' to the CLI other than to `remove` **before** a window is killed.
This gives the daemon the chance to adjust its internal state accordingly. 

Such a config may look like this:
```
mod1 + shift + q
  scratch-cli remove && bspc node -c
```
The CLI will fail gracefully (exit 0) if the daemon cannot be reached.