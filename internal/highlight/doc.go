// Package highlight provides ANSI terminal color support for logslice output.
//
// It wraps log lines or individual fields with escape sequences when output
// is directed to a TTY. Highlighting can be forced on or off regardless of
// terminal detection using ModeOn and ModeOff respectively.
//
// Usage:
//
//	h := highlight.New(highlight.ModeAuto, isatty.IsTerminal(os.Stdout.Fd()))
//	fmt.Println(h.Line(line, highlight.Green))
//	fmt.Println(h.Field(line, "level"))
package highlight
