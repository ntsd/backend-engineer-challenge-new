package scanner

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/ntsd/backend-engineer-challenge-new/internal/model"
	"github.com/sirupsen/logrus"
)

// scan start scan repository
func (s *scanner) scan(logger *logrus.Entry, scan model.Scan) (model.Findings, error) {
	// TODO: make it faster by clone to temp folder
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:      scan.RepositoryURL,
		Progress: logger.Logger.Out,
	})
	if err != nil {
		return nil, fmt.Errorf("error cloning repository: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("error to get repository head: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("error to get commit object: %w", err)
	}

	files, err := commit.Files()
	if err != nil {
		return nil, fmt.Errorf("error to get commit files: %w", err)
	}

	finding := model.Findings{}

	scanRegex := regexp.MustCompile(`\s(?P<key>(private_key|public_key)(\w+|\s))`)
	if err := files.ForEach(func(f *object.File) error {
		// skip binary files
		isBin, err := f.IsBinary()
		if err != nil {
			return fmt.Errorf("error to check is binary: %w", err)
		}
		if isBin {
			return nil
		}

		lines, err := f.Lines()
		if err != nil {
			return fmt.Errorf("error to get lines: %w", err)
		}

		for lineNum, line := range lines {
			founds := scanLine(scanRegex, line)
			for _, found := range founds {
				finding = append(finding, model.Finding{
					Path:        f.Name,
					Line:        lineNum,
					Description: fmt.Sprintf("`%s` might be a secret", found),
				})
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error in for each files: %w", err)
	}

	return finding, nil
}

func scanLine(scanRegex *regexp.Regexp, line string) []string {
	matchs := scanRegex.FindAllStringSubmatch(line, -1)
	found := []string{}
	for _, match := range matchs {
		for i, group := range match {
			// only get the 2nd group
			if i < 1 {
				continue
			}
			found = append(found, strings.TrimSpace(group))
			break
		}
	}
	return found
}
