package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type editor struct {
	exe  string
	file string

	openPRs   []pullRequest
	closedPRs []pullRequest
}

func newEditor() (*editor, error) {
	e := &editor{}

	for _, name := range split(editors) {
		exe, err := exec.LookPath(name)
		if err == nil {
			e.exe = exe
			break
		}
	}

	if e.exe == "" {
		return e, nil
	}

	file, err := os.CreateTemp("", "status-email")
	if err != nil {
		return nil, err
	}
	file.Close()
	e.file = file.Name()

	return e, nil
}

func (e *editor) cleanup() {
	os.Remove(e.file)
}

func (e *editor) addPullRequest(pr pullRequest) {
	if pr.isOpen() {
		e.openPRs = append(e.openPRs, pr)
	} else if pr.closedThisWeek() {
		e.closedPRs = append(e.closedPRs, pr)
	}
}

func (e *editor) generateRawDocument() string {
	docFormat := `
# Status Email Editor
# Please review the autogenerated pull request data and edit accordingly.
# Lines starting with '#' will be ignored, and exiting without saving
# aborts the program.

Active pull requests 🚀:
%s


Closed pull request (this week) 🔒:
%s


# Blockers 🚫:
# No blockers.
#
#
# Plans 🗓:
# The same thing we do every night, Pinky.
#
`

	activePullRequests := ""
	for _, issue := range e.openPRs {
		activePullRequests += issue.String() + "\n\n"
	}
	if activePullRequests == "" {
		activePullRequests = wip
	}

	closedPullRequests := ""
	for _, issue := range e.closedPRs {
		closedPullRequests += issue.String() + "\n\n"
	}
	if closedPullRequests == "" {
		closedPullRequests = wip
	}

	return strings.TrimSpace(fmt.Sprintf(
		docFormat,
		strings.TrimSpace(activePullRequests),
		strings.TrimSpace(closedPullRequests),
	))
}

func (e *editor) editDocument(raw string) (string, error) {
	if e.exe == "" {
		return "", errors.New("editor not found")
	}

	if err := os.WriteFile(e.file, []byte(raw), 0o644); err != nil {
		return "", err
	}

	initialFileStat, err := os.Stat(e.file)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(e.exe, e.file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	fileStat, err := os.Stat(e.file)
	if err != nil {
		return "", err
	}

	if fileStat.Size() == initialFileStat.Size() && fileStat.ModTime() == initialFileStat.ModTime() {
		return "", errors.New("status email not saved")
	}

	readFile, err := os.Open(e.file)
	if err != nil {
		return "", err
	}
	defer readFile.Close()

	var result string
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if !strings.HasPrefix(line, "#") {
			result += line + "\n"
		}
	}

	return strings.TrimSpace(result), nil
}
