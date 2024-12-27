package gopushSvc

import "errors"

var (
	ErrPullFailed           = errors.New("remote pull failed")
	ErrMergeFailed          = errors.New("merge failed")
	ErrBranchInvalid        = errors.New("invalid branch")
	ErrBranchAlreadyExist   = errors.New("branch already exist")
	ErrTestsFailed          = errors.New("tests failed")
	ErrRemoteNotLoaded      = errors.New("remote not loaded")
	ErrRemoteNotFound       = errors.New("no remotes found")
	ErrRemoteAlreadyExists  = errors.New("remote already exists")
	ErrRepoAlreadyExists    = errors.New("repository already exists")
	ErrRepoNotFound         = errors.New("repository not found")
	ErrConfigNotLoaded      = errors.New("config not loaded")
	ErrFileNotFound         = errors.New("file not found")
	ErrAuthNotFound         = errors.New("auth credentials not found")
	ErrInvalidAuthMethod    = errors.New("invalid auth method")
	ErrWaitExit             = errors.New("waiting for user")
	ErrInvalidPassphrase    = errors.New("invalid passphrase")
	ErrKeyNotSupported      = errors.New("invalid key on remote")
	ErrAuthLoadFailed       = errors.New("failed to load auth")
	ErrAlreadyUpToDate      = errors.New("already up to date")
	ErrRemoteBranchNotFound = errors.New("remote branch not found")
)
