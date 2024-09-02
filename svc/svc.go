package svc

import (
	"github.com/seriouspoop/gopush/config"
)

type Svc struct {
	git  gitHelper
	bash scriptHelper
	cfg  *config.Config
}

func New(git gitHelper, bash scriptHelper) *Svc {
	return &Svc{
		git:  git,
		bash: bash,
	}
}
