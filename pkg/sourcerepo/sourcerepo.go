package sourcerepo

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type SourceRepo struct {
	repo *git.Repository
}

// New returns initialized and cloned *SourceRepo
func New(url string, directory string) (sr *SourceRepo, err error) {
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	sourceRepo := &SourceRepo{
		repo: r,
	}
	return sourceRepo, nil
}

// Checkout checks out the commitId of the current repository
func (sr *SourceRepo) Checkout(commitId string) error {
	wt, err := sr.repo.Worktree()
	if err != nil {
		return err
	}
	err = wt.Pull(&git.PullOptions{
		Force:    true,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	err = wt.Checkout(&git.CheckoutOptions{
		Force: true,
		Hash:  plumbing.NewHash(commitId),
	})
	if err != nil {
		return err
	}
	return nil
}
