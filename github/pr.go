package github

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"text/template"

	"github.com/google/go-github/github"
	"github.com/igorbrites/neptune/terraform"
	"golang.org/x/oauth2"
)

type PullRequest struct {
	Owner string
	Repo string
	Number int
}

type TemplateInput struct {
	Folder string
	Output string
	Template string
	Workspace string
}

func (pr *PullRequest) Comment(plan terraform.Plan) {
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
		Folder: folder,
		Template: "no-changes.tmpl",
		Workspace: plan.Workspace,
	}

	if plan.Type == terraform.Error {
		tmpl.Template = "error.tmpl"
		tmpl.Output = plan.ProcessedError()
	} else if plan.Type == terraform.Changed {
		tmpl.Template = "plan.tmpl"
		tmpl.Output = plan.ProcessedOutput()
	}

	_, filename, _, _ := runtime.Caller(1)

	var b bytes.Buffer
	t, _ := template.New(tmpl.Template).ParseFiles(path.Join(path.Dir(filename), tmpl.Template))
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
