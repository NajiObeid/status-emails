package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
)

const (
	red    = "🔴"
	yellow = "🟡"
	green  = "🟢"
	done   = "✅"
	wip    = "🛠️"

	day = 24 * time.Hour
)

type pullRequest struct {
	title     string
	state     string
	url       string
	createdAt time.Time
	updatedAt time.Time
	closedAt  time.Time
}

func newPullRequest(issue *github.Issue) pullRequest {
	if issue == nil {
		return pullRequest{}
	}

	return pullRequest{
		title:     deref(issue.Title),
		state:     deref(issue.State),
		url:       deref(issue.HTMLURL),
		createdAt: deref(issue.CreatedAt),
		updatedAt: deref(issue.UpdatedAt),
		closedAt:  deref(issue.ClosedAt),
	}
}

func (pr *pullRequest) isOpen() bool {
	return pr.state == "open"
}

func (pr *pullRequest) closedThisWeek() bool {
	return pr.state == "closed" && pr.closedAt.After(lastMonday(time.Now()))
}

func (pr *pullRequest) String() string {
	format := `
%s %s
Title: %s
Age: %s
# Comment: --
`

	age := pr.closedAt.Sub(pr.createdAt)
	marker := done
	if pr.isOpen() {
		age = time.Now().Sub(pr.createdAt)
		if age >= 0*day && age < 7*day {
			marker = green
		} else if age >= 7*day && age < 21*day {
			marker = yellow
		} else {
			marker = red
		}
	}

	ageDays := int(age.Hours() / 24)
	return strings.TrimSpace(fmt.Sprintf(format, marker, pr.url, pr.title, pluralizeDays(ageDays)))
}

func deref[T any](v *T) T {
	if v != nil {
		return *v
	}

	return zeroValue[T]()
}

func zeroValue[T any]() T {
	var result T
	return result
}

func pluralizeDays(n int) string {
	if n == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", n)
}

func lastMonday(t time.Time) time.Time {
	y, m, d := t.Date()
	daysPastLastMonday := int((7 + t.Weekday() - time.Monday) % 7)
	return time.Date(y, m, d-daysPastLastMonday, 0, 0, 0, 0, time.UTC)
}
