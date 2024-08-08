package svc

import (
	"bufio"
	"os"

	"github.com/seriouspoop/gopush/config"
)

type Svc struct {
	git  gitHelper
	bash scriptHelper
	r    *bufio.Reader
	cfg  *config.Config
}

func New(git gitHelper, bash scriptHelper) *Svc {
	r := bufio.NewReader(os.Stdin)
	return &Svc{
		git:  git,
		bash: bash,
		r:    r,
	}
}
