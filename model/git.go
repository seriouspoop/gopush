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

type Password string

func (p Password) Valid() bool {
	return len(string(p)) > 0
}

func (p Password) String() string {
	return string(p)
}

type Provider int

const (
	ProviderUNKOWN Provider = iota
	ProviderGITHUB
	ProviderBITBUCKET
	ProviderGITLAB
)

func (p Provider) String() string {
	if p == ProviderGITHUB {
		return "GitHub"
	} else if p == ProviderBITBUCKET {
		return "BitBucket"
	} else if p == ProviderGITLAB {
		return "GitLab"
	}
	return ""
}

func (p Provider) HostURL() string {
	providerToHostMap := map[Provider]string{
		ProviderGITHUB:    "github.com",
		ProviderBITBUCKET: "bitbucket.org",
		ProviderGITLAB:    "gitlab.com",
	}

	return providerToHostMap[p]
}

type AuthMode int

const (
	AuthUNKNOWN AuthMode = iota
	AuthHTTP
	AuthSSH
)

type Remote struct {
	Name string
	Url  string
}

func (r *Remote) Provider() Provider {
	if strings.Contains(r.Url, "github") {
		return ProviderGITHUB
	} else if strings.Contains(r.Url, "bitbucket") {
		return ProviderBITBUCKET
	} else if strings.Contains(r.Url, "gitlab") {
		return ProviderGITLAB
	} else {
		return ProviderUNKOWN
	}
}

func (r *Remote) AuthMode() AuthMode {
	if strings.Contains(r.Url, "http") || strings.Contains(r.Url, "https") {
		return AuthHTTP
	} else if strings.Contains(r.Url, "git@") {
		return AuthSSH
	} else {
		return AuthUNKNOWN
	}
}
