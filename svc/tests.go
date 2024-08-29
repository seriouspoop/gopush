package svc

import "fmt"

func (s *Svc) CheckTestsAndRun() (bool, error) {
	present, err := s.bash.TestsPresent()
	if err != nil {
		return false, err
	}
	if present {
		output, _ := s.bash.GenerateMocks()
		fmt.Print(output)
		output, err := s.bash.RunTests()
		fmt.Print(output)
		if err != nil {
			return true, ErrTestsFailed
		}
	}
	return false, nil
}
