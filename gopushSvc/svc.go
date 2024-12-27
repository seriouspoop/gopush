package gopushSvc

import (
	"github.com/seriouspoop/gopush/config"
	"github.com/seriouspoop/gopush/model"
)

type Svc struct {
	git        gitHelper
	bash       scriptHelper
	cfg        *config.Config
	passphrase model.Password
}

func New(git gitHelper, bash scriptHelper) *Svc {
	return &Svc{
		git:  git,
		bash: bash,
	}
}
