package sourcerepo

import (
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
)

// SourceRepo is the struct for the Source Repository
type SourceRepo struct {
	repo      *git.Repository
	URL       string
	Directory string
}

// New returns initialized and cloned *SourceRepo
func New(url string, directory string) (sr *SourceRepo, err error) {
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		r, err = git.PlainOpen(directory)
		if err != nil {
			log.Printf("Clone failed: %+v\n", err)
			return nil, err
		}
	}

	sourceRepo := &SourceRepo{
		repo:      r,
		URL:       url,
		Directory: directory,
	}
	return sourceRepo, nil
}

// Checkout checks out the commitID of the current repository
func (sr *SourceRepo) Checkout(commitID string) error {
	log.Printf("Checking out commit: %s\n", commitID)
	wt, err := sr.repo.Worktree()
	if err != nil {
		log.Printf("Getting Worktree failed: %+v\n", err)
		return err
	}
	err = wt.Pull(&git.PullOptions{
		Force:    true,
		Progress: os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Printf("Pull failed: %+v\n", err)
		return err
	}
	err = wt.Checkout(&git.CheckoutOptions{
		Force: true,
		Hash:  plumbing.NewHash(commitID),
	})
	if err != nil {
		log.Printf("Checkout failed: %+v\n", err)
		return err
	}
	return nil
}
