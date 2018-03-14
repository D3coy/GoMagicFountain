package main

const (
	usage = `
Usage: !roll [options] [d#] [quantity]

Examples:
    - !roll d4
    - !roll d10 3
    - !roll d20 d10 3 d5
    - !roll -t 3 d4 5
    - !roll -s d5

Options:
	-t [#]      Gets the highest [#] rolls from the rolls made
	-s          Used to Sort the resulting dice rolls
`
)
