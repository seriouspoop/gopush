package svc

import "errors"

var (
	ErrBranchInvalid       = errors.New("invalid branch")
	ErrBranchAlreadyExist  = errors.New("branch already exist")
	ErrTestsFailed         = errors.New("tests failed")
	ErrRemoteNotLoaded     = errors.New("remote not loaded")
	ErrRemoteNotFound      = errors.New("no remotes found")
	ErrRemoteAlreadyExists = errors.New("remote already exists")
	ErrRepoAlreadyExists   = errors.New("repository already exists")
	ErrRepoNotFound        = errors.New("repository not found")
	ErrConfigNotLoaded     = errors.New("config not loaded")
	ErrFileNotFound        = errors.New("file not found")
)
