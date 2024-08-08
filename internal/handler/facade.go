package handler

type servicer interface {
	LoadProject() error
	InitializeRepo() error
	InitializeRemote() error
	SetUserPreference() error
	SetRemoteAuth() error
	LoadConfig() error
	// FetchAndMerge() error
	Pull(force bool) error
	StageChanges() error
}
