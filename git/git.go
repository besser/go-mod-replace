package git

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"

	git "github.com/go-git/go-git/v5"
)

func GetRepo(name string) (*git.Repository, error) {
	repository, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: name,
	})
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func GetCurrentBranchFromRepository(r *git.Repository) (string, error) {
	branchRefs, err := r.Branches()
	if err != nil {
		return "", err
	}

	headRef, err := r.Head()
	if err != nil {
		return "", err
	}

	var currentBranchName string
	err = branchRefs.ForEach(func(branchRef *plumbing.Reference) error {
		if branchRef.Hash() == headRef.Hash() {
			currentBranchName = branchRef.Name().String()

			return nil
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return currentBranchName, nil
}

func GetCurrentCommitFromRepository(repository *git.Repository) (string, error) {
	headRef, err := repository.Head()
	if err != nil {
		return "", err
	}
	headSha := headRef.Hash().String()

	return headSha, nil
}

func GetLatestTagFromRepository(r *git.Repository) (string, error) {
	tagRefs, err := r.Tags()
	if err != nil {
		return "", err
	}

	var (
		latestTagCommit *object.Commit
		latestTagName   string
	)

	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := r.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := r.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().Short()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().Short()
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return latestTagName, nil
}
