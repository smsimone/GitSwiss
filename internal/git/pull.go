package git

import (
	"context"
)

func Pull(ctx context.Context, path string) error {
	err := CheckGitRepo(path)
	if err == nil {
		return Execute(GitCommand{
			Path: path, Options: []string{"pull"},
		})
	}
	return err
}

func MergePull(ctx context.Context, path string, remote string, source string) error {
	err := CheckGitRepo(path)
	if err == nil {
		return Execute(GitCommand{
			Path: path, Options: []string{"pull", remote, source},
		})
	}
	return err
}
