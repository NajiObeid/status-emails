package main

import "time"

type issue struct {
	title     string
	state     string
	url       string
	createdAt time.Time
	updatedAt time.Time
	closedAt  time.Time
}

func (i *issue) String() string {
	return i.title + " " + i.state
}
