package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Stash struct {
	Index   int
	Ref     string
	Branch  string
	Message string
	Date    time.Time
}

func gitCmd(args ...string) *exec.Cmd {
	cwd, _ := os.Getwd()
	cmd := exec.Command("git", args...)
	cmd.Dir = cwd
	return cmd
}

func IsRepo() bool {
	return gitCmd("rev-parse", "--git-dir").Run() == nil
}

func List() ([]Stash, error) {
	out, err := gitCmd("stash", "list", "--format=%gd\x00%gs\x00%ci").Output()
	if err != nil {
		if len(out) == 0 {
			return nil, fmt.Errorf("git error")
		}
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return nil, nil
	}

	lines := strings.Split(raw, "\n")
	stashes := make([]Stash, 0, len(lines))

	for i, line := range lines {
		parts := strings.SplitN(line, "\x00", 3)
		if len(parts) < 3 {
			continue
		}

		ref := parts[0]
		msg := parts[1]
		dateStr := strings.TrimSpace(parts[2])

		date, _ := time.Parse("2006-01-02 15:04:05 -0700", dateStr)

		stashes = append(stashes, Stash{
			Index:   i,
			Ref:     ref,
			Branch:  extractBranch(msg),
			Message: msg,
			Date:    date,
		})
	}

	return stashes, nil
}

func extractBranch(msg string) string {
	msg = strings.TrimPrefix(msg, "WIP on ")
	msg = strings.TrimPrefix(msg, "On ")
	if idx := strings.Index(msg, ":"); idx != -1 {
		return msg[:idx]
	}
	return ""
}

func Diff(ref string) (string, error) {
	stat, err := gitCmd("stash", "show", "--stat", ref).Output()
	if err != nil {
		return "", err
	}

	patch, err := gitCmd("stash", "show", "-p", ref).Output()
	if err != nil {
		return string(stat), nil
	}

	return string(bytes.TrimRight(stat, "\n")) + "\n\n" + string(patch), nil
}

func Apply(ref string) error {
	out, err := gitCmd("stash", "apply", ref).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func Pop(index int) error {
	out, err := gitCmd("stash", "pop", "stash@{"+strconv.Itoa(index)+"}").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func Drop(index int) error {
	out, err := gitCmd("stash", "drop", "stash@{"+strconv.Itoa(index)+"}").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func Push(message string) error {
	args := []string{"stash", "push"}
	if message != "" {
		args = append(args, "-m", message)
	}
	out, err := gitCmd(args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}
