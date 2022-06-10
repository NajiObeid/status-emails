package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

func main() {
	gh := github.NewClient(
		oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken}),
		),
	)

	editor, err := newEditor()
	if err != nil {
		fmt.Printf("Could not init editor: %+v\n", err)
		os.Exit(1)
	}

	for _, repo := range split(repositories) {
		query := fmt.Sprintf("repo:%s author:%s is:pr", repo, githubUsername)
		searchResult, resp, err := gh.Search.Issues(context.Background(), query, nil)
		if err != nil {
			fmt.Printf("Could not search issues %+v\n", resp.Response)
			os.Exit(1)
		}

		for _, issue := range searchResult.Issues {
			editor.addPullRequest(newPullRequest(issue))
		}
	}

	rawDocument := editor.generateRawDocument()
	statusEmail, err := editor.editDocument(rawDocument)
	if err != nil {
		fmt.Printf("Error editing document: %+v\n", err)
		fmt.Printf("Printing raw status email to console\n\n")
		fmt.Printf("%+v\n", rawDocument)
		os.Exit(1)
	}

	if dryRun {
		fmt.Printf("Dry Run: printing status email to console\n\n")
		fmt.Printf("%+v\n", statusEmail)
		os.Exit(0)
	}

	email := newEmail(googleUsername, googleAppPassword)
	err = email.send(split(mailingList), statusEmail)
	if err != nil {
		fmt.Printf("Could not send email: %+v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully sent status email!")
}

func split(s string) []string {
	res := strings.Split(s, ",")
	for i := range res {
		res[i] = strings.TrimSpace(res[i])
	}
	return res
}
