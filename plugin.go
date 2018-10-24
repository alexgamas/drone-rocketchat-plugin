package main

import (
	"./rocketchat"
	"fmt"
	"github.com/drone/drone-template-lib/template"
)

type (
	Repo struct {
		Owner   string
		Name    string
		Link    string
		Avatar  string
		Branch  string
		Private bool
		Trusted bool
	}

	Build struct {
		Number   int
		Event    string
		Status   string
		Deploy   string
		Created  int64
		Started  int64
		Finished int64
		Link     string
	}

	Commit struct {
		Remote  string
		Sha     string
		Ref     string
		Link    string
		Pull    string
		Branch  string
		Message string
		Author  Author
	}

	Author struct {
		Name   string
		Email  string
		Avatar string
	}

	Config struct {
		// plugin-specific parameters and secrets
		Channel   string
		Text      string
		Username  string
		Password  string
		Url       string
		Template  string
		UserId    string
		AuthToken string
		IconURL   string
		IconEmoji string
		ImageURL  string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Commit Commit
		Config Config
	}
)

func (p Plugin) Exec() error {

	client := rocketchat.New(p.Config.Url, p.Config.UserId, p.Config.AuthToken)

	attachment := rocketchat.Attachment{
		Text:     message(p.Repo, p.Build, p.Commit),
		Color:    color(p.Build),
		ImageURL: p.Config.ImageURL,
	}

	payload := rocketchat.WebHookPostPayload{}
	payload.Username = p.Config.Username
	payload.Attachments = []*rocketchat.Attachment{&attachment}
	payload.IconUrl = p.Config.IconURL
	payload.IconEmoji = p.Config.IconEmoji

	if p.Config.Template != "" {

		txt, err := template.RenderTrim(p.Config.Template, p)

		if err != nil {
			return err
		}

		attachment.Text = txt
	}

	if p.Config.Username != "" {
		req := &rocketchat.LoginRequest{p.Config.Username, p.Config.Password}
		err := client.Login(req)
		if err != nil {
			return err
		}
	}
	return client.PostMessage(&payload)
}

func message(repo Repo, build Build, commit Commit) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		commit.Sha[:8],
		commit.Branch,
		commit.Author,
	)
}

func fallback(repo Repo, build Build, commit Commit) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		commit.Sha[:8],
		commit.Branch,
		commit.Author,
	)
}

func color(build Build) string {
	switch build.Status {
	case "success":
		return "good"
	case "failure", "error", "killed":
		return "danger"
	default:
		return "warning"
	}
}