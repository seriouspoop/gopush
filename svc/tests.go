package svc

import (
	"github.com/seriouspoop/gopush/utils"
)

func (s *Svc) CheckTestsAndRun() (bool, error) {
	present, err := s.bash.TestsPresent()
	if err != nil {
		return false, err
	}
	if present {
		s.bash.GenerateMocks()
		output, err := s.bash.RunTests()
		if err != nil {
			utils.Logger(utils.LOG_INFO, utils.Faint(output))
			return false, ErrTestsFailed
		}
		return true, nil
	}
	return false, nil
}
