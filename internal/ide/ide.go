package ide

import (
	"errors"
	"os"

	"github.com/eskelinenantti/tmuxide/internal/project"
)

type Tmux interface {
	HasSession(name string) bool
	New(session string, dir string) error
	NewWindow(session string, dir string) error
	Attach(session string) error
	Switch(session string) error
}

type ShellPath interface {
	Contains(path string) bool
}

var ErrTmuxNotInstalled = errors.New(
	"Did not find tmux, which is a required dependency for ide command.\n\n" +

		"You can install tmux e.g. via homebrew by running\n" +
		"brew install tmux\n",
)

func Start(project project.Project, tmux Tmux, path ShellPath) error {
	if !path.Contains("tmux") {
		return ErrTmuxNotInstalled
	}

	if !tmux.HasSession(project.Name) {
		if err := tmux.New(project.Name, project.WorkingDir); err != nil {
			return err
		}
	}

	if isAttached() {
		return tmux.Switch(project.Name)
	}

	return tmux.Attach(project.Name)
}

func isAttached() bool {
	_, isAttached := os.LookupEnv("TMUX")
	return isAttached
}
