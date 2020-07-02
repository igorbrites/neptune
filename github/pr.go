package github

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/google/go-github/github"
	"github.com/igorbrites/neptune/terraform"
	"golang.org/x/oauth2"
)

type PullRequest struct {
	Owner  string
	Repo   string
	Number int
}

type TemplateInput struct {
	Folder    string
	Output    string
	Template  string
	Workspace string
}

const (
	errorTemplate     = ":rotating_light: There was an error running plan on directory `{{.Folder}}`, using workspace `{{.Workspace}}`:\n```\n{{.Output}}\n```"
	noChangesTemplate = ":heavy_check_mark: No changes found on directory `{{.Folder}}`, using workspace `{{.Workspace}}`!"
	planTemplate      = ":heavy_check_mark: Here is the plan for directory `{{.Folder}}`, using workspace `{{.Workspace}}`:\n\n<details><summary>:warning: Click here to show the plan fully</summary>\n\n```diff\n{{.Output}}\n```\n</details>"
)

func (pr *PullRequest) Comment(plan terraform.Plan) {
	if pr.Number <= 0 {
		fmt.Println("Skipping GitHub comment...")
		return
	}

	comment := pr.GenerateCommentText(plan)

	_, _, err := pr.GetService().CreateComment(context.Background(), pr.Owner, pr.Repo, pr.Number, comment)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (pr *PullRequest) GenerateCommentText(plan terraform.Plan) *github.IssueComment {
	var comment string
	folder, _ := os.Getwd()
	tmpl := TemplateInput{
		Folder:    folder,
		Template:  noChangesTemplate,
		Workspace: plan.Workspace,
	}

	if plan.Type == terraform.Error {
		tmpl.Template = errorTemplate
		tmpl.Output = plan.ProcessedError()
	} else if plan.Type == terraform.Changed {
		tmpl.Template = planTemplate
		tmpl.Output = plan.ProcessedOutput()
	}

	var b bytes.Buffer
	t, _ := template.New("comment").Parse(tmpl.Template)
	err := t.Execute(&b, tmpl)
	if err != nil {
		panic(err)
	}
	comment = b.String()

	return &github.IssueComment{
		Body: &comment,
	}
}

func (pr *PullRequest) GetService() *github.IssuesService {
	token, present := os.LookupEnv("GITHUB_TOKEN")
	if !present || token == "" {
		panic("You must set yout GitHub Token on the GITHUB_TOKEN environment variable!")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return client.Issues
}
