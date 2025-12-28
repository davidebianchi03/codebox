package commands

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"syscall"

	"golang.org/x/term"
)

func PrepareTerminal() *term.Terminal {
	if !term.IsTerminal(syscall.Stdin) {
		slog.Warn("std::cin is not a terminal")
	}
	if !term.IsTerminal(syscall.Stdout) {
		slog.Warn("std::cout is not a terminal")
	}

	oldState, err := term.MakeRaw(syscall.Stdin)
	defer func() {
		err = errors.Join(err, term.Restore(syscall.Stdin, oldState))
	}()

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	terminal := term.NewTerminal(screen, "")
	return terminal
}
