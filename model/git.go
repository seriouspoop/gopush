package model

import (
	"strings"
)

type Branch string

func (b Branch) String() string {
	return string(b)
}

func (b Branch) Valid() bool {
	return len(b.String()) > 0
}

type Provider int

const (
	UNKOWN Provider = iota
	GITHUB
	BITBUCKET
	GITLAB
)

func (p Provider) String() string {
	if p == GITHUB {
		return "GitHub"
	} else if p == BITBUCKET {
		return "BitBucket"
	} else if p == GITLAB {
		return "GitLab"
	}
	return ""
}

type Remote struct {
	Name string
	Url  string
}

func (r *Remote) Provider() Provider {
	if strings.Contains(r.Url, "github") {
		return GITHUB
	} else if strings.Contains(r.Url, "bitbucket") {
		return BITBUCKET
	} else if strings.Contains(r.Url, "gitlab") {
		return GITLAB
	} else {
		return UNKOWN
	}
}
