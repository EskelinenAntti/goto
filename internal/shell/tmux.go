package shell

import (
	"os"
	"os/exec"
)

type Tmux struct{}

func (Tmux) HasSession(session string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", session)
	return cmd.Run() == nil
}

func (Tmux) New(session string, dir string) error {
	return exec.Command("tmux", "new-session", "-ds", session, "-c", dir).Run()
}

func (Tmux) NewWindow(session string, dir string) error {
	return exec.Command("tmux", "new-window", "-d", "-t", session, "-c", dir).Run()
}

func (Tmux) Attach(session string) error {
	tmuxCmd := exec.Command("tmux", "attach", "-t", session)
	return attachAndRun(tmuxCmd)
}

func (Tmux) Switch(session string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", session)
	return attachAndRun(cmd)
}

func attachAndRun(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
