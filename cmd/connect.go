package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"github.com/tryoasnafi/gate/internal/session"
	"github.com/tryoasnafi/gate/internal/store.go"
)

func ConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect <label>",
		Short: "Connect to a saved SSH label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConnect(args[0])
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	return cmd
}

func runConnect(label string) error {
	pwd := session.Require()

	entry, err := store.GetEntry(pwd, label)
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            entry.User,
		Auth:            []ssh.AuthMethod{ssh.Password(entry.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", entry.Host, entry.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("SSH connection failed: %w", err)
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer sess.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to set terminal raw mode: %w", err)
	}
	defer term.Restore(fd, oldState)

	width, height, err := term.GetSize(fd)
	if err != nil {
		width, height = 80, 24
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := sess.RequestPty("xterm", height, width, modes); err != nil {
		return fmt.Errorf("RequestPty failed: %w", err)
	}

	sess.Stdin = os.Stdin
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	if err := sess.Shell(); err != nil {
		return fmt.Errorf("shell error: %w", err)
	}

	go handleResize(sess, fd)
	go handleSignals(sess, fd, oldState)

	if err := sess.Wait(); err != nil {
		var exitErr *ssh.ExitError
		if errors.As(err, &exitErr) {
			// Return non-zero exit code for parent process
			return fmt.Errorf("remote exited with status %d", exitErr.ExitStatus())
		}
		return fmt.Errorf("session wait error: %w", err)
	}

	return nil
}

func handleSignals(sess *ssh.Session, fd int, oldState *term.State) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for sig := range signals {
		switch sig {
		case syscall.SIGINT:
			// Forward Ctrl+C to remote session (remote receives ^C)
			sess.Signal(ssh.SIGINT)
		case syscall.SIGTERM, syscall.SIGQUIT:
			// Cleanup terminal, notify and exit
			term.Restore(fd, oldState)
			fmt.Fprintf(os.Stderr, "\nReceived %s. Exiting gracefully...\n", sig)
			os.Exit(0)
		}
	}
}

func handleResize(sess *ssh.Session, fd int) {
	sigWinch := make(chan os.Signal, 1)
	signal.Notify(sigWinch, syscall.SIGWINCH)
	for range sigWinch {
		width, height, err := term.GetSize(fd)
		if err == nil {
			sess.WindowChange(height, width)
		}
	}
}
