package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/eskelinenantti/tmuxide/internal/ide"
	"github.com/eskelinenantti/tmuxide/internal/project"
	"github.com/eskelinenantti/tmuxide/internal/test/mock"
	"github.com/eskelinenantti/tmuxide/internal/test/spy"
)

const command string = "cs"

func TestRunWithFile(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	file := dir + "/file.txt"
	os.WriteFile(file, []byte{}, 0644)

	tmux := &spy.Tmux{}

	shell := shellEnv{
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, file}, shell)

	if got, want := err, project.ErrNotADirectory; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	var expectedCalls [][]string
	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithDirectory(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	tmux := &spy.Tmux{}

	shell := shellEnv{
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, dir},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithoutArguments(t *testing.T) {
	os.Unsetenv("TMUX")

	dir := t.TempDir()
	t.Chdir(dir)

	tmux := &spy.Tmux{}

	shell := shellEnv{
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, dir},
		{"Attach", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunHelp(t *testing.T) {
	os.Unsetenv("EDITOR")
	os.Unsetenv("TMUX")

	tmux := &spy.Tmux{}
	dir := t.TempDir()

	shell := shellEnv{
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir, "-h"}, shell)

	if got, want := err.Error(), fmt.Sprintf(helpMsgTemplate, command); got != want {
		t.Fatalf("got=%v, want=%v", got, want)
	}

	var expectedCalls [][]string
	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithTmuxSessionExists(t *testing.T) {
	t.Setenv("TMUX", "test")

	dir := t.TempDir()
	session := project.Name(dir)

	tmux := &spy.Tmux{
		Sessions: session,
	}

	shell := shellEnv{
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	expectedCalls := [][]string{
		{"HasSession", session},
		{"Switch", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunInsideTmux(t *testing.T) {
	t.Setenv("TMUX", "test")

	dir := t.TempDir()
	tmux := &spy.Tmux{}

	shell := shellEnv{
		Tmux: tmux,
		Path: mock.Path{},
	}

	err := run([]string{command, dir}, shell)

	if err != nil {
		t.Fatalf("err=%v", err)
	}

	session := project.Name(dir)

	expectedCalls := [][]string{
		{"HasSession", session},
		{"New", session, dir},
		{"Switch", session},
	}

	if got, want := tmux.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}

func TestRunWithoutTmux(t *testing.T) {
	os.Unsetenv("TMUX")

	tmuxSpy := &spy.Tmux{}
	dir := t.TempDir()

	shell := shellEnv{
		Tmux: tmuxSpy,
		Path: mock.Path{Missing: []string{"tmux"}},
	}

	err := run([]string{command, dir}, shell)

	if got, want := err, ide.ErrTmuxNotInstalled; !errors.Is(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
	var expectedCalls [][]string
	if got, want := tmuxSpy.Calls, expectedCalls; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
