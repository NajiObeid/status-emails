package main

const (

	// name is the name that will be used in email subjects.
	name = "Buzz"

	// googleUsername is your email address.
	googleUsername = "blightyear@toy.story"

	// googleAppPassword is used for authentication. Generate a new appPassword
	// from your google account settings.
	googleAppPassword = ""

	// mailingList is a comma separated list of recipients.
	mailingList = "woody@toy.story"

	// githubUsername is the github username used to fetch pull request data.
	githubUsername = "buzzlightyear"

	// githubToken is a personal access token used to access the github api.
	githubToken = ""

	// repositories is a comma separated list of repositories used to aggregate
	// pull request data.
	repositories = "golang/go"

	// editors is a comma separated list of editors used to review and edit emails
	// before sending.
	editors = "vi, nano"

	// dryRun will print the email body to console.
	dryRun = false
)
