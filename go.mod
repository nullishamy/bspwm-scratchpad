module github.com/nullishamy/bspwm-scratchpad/v2

go 1.20

require (
	github.com/alecthomas/kong v0.9.0
	github.com/nullishamy/bspwm-scratchpad/transport v0.0.0
)

replace github.com/nullishamy/bspwm-scratchpad/transport => ./transport
