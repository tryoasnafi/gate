package cmd

import (
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
	return &cobra.Command{
		Use:   "connect",
		Short: "Connect to a saved SSH label",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pwd := session.Require()
			entry, err := store.GetEntry(pwd, args[0])
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			config := &ssh.ClientConfig{
				User:            entry.User,
				Auth:            []ssh.AuthMethod{ssh.Password(entry.Password)},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}

			addr := fmt.Sprintf("%s:%d", entry.Host, entry.Port)
			client, err := ssh.Dial("tcp", addr, config)
			if err != nil {
				fmt.Println("SSH connection failed:", err)
				os.Exit(1)
			}
			defer client.Close()

			sess, err := client.NewSession()
			if err != nil {
				fmt.Println("Failed to create session:", err)
				os.Exit(1)
			}
			defer sess.Close()

			fd := int(os.Stdin.Fd())
			oldState, err := term.GetState(fd)
			if err != nil {
				fmt.Println("Failed to get terminal state:", err)
				os.Exit(1)
			}

			// Handle terminal resize (SIGWINCH)
			go func() {
				sigWinch := make(chan os.Signal, 1)
				signal.Notify(sigWinch, syscall.SIGWINCH)
				for range sigWinch {
					width, height, err := term.GetSize(fd)
					if err == nil {
						_ = sess.WindowChange(height, width)
					}
				}
			}()

			if _, err := term.MakeRaw(fd); err != nil {
				fmt.Println("Failed to set terminal raw mode:", err)
				os.Exit(1)
			}

			// Ensure terminal state restored no matter what
			defer func() {
				_ = term.Restore(fd, oldState)
			}()

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
				fmt.Println("RequestPty failed:", err)
				os.Exit(1)
			}

			sess.Stdin = os.Stdin
			sess.Stdout = os.Stdout
			sess.Stderr = os.Stderr

			// Forward Ctrl+C
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			go func() {
				for range ch {
					sess.Signal(ssh.SIGINT)
				}
			}()

			if err := sess.Shell(); err != nil {
				fmt.Println("Shell error:", err)
				os.Exit(1)
			}

			if err := sess.Wait(); err != nil {
				if exitErr, ok := err.(*ssh.ExitError); ok {
					os.Exit(exitErr.ExitStatus())
				} else {
					fmt.Println("Session wait error:", err)
					os.Exit(1)
				}
			}
		},
	}
}
