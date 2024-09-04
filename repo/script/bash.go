package script

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/seriouspoop/gopush/model"
)

type Error struct {
	FileNotExists error
}

type Bash struct {
	err *Error
}

func New(svcError *Error) *Bash {
	return &Bash{
		err: svcError,
	}
}

func (b *Bash) GetCurrentBranch() (model.Branch, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return model.Branch(""), err
	}
	return model.Branch(string(output[:len(output)-1])), nil
}

func (b *Bash) PullBranch(remoteName string, branch model.Branch, force bool) (string, error) {
	var cmd *exec.Cmd
	if force {
		cmd = exec.Command("git", "pull", remoteName, branch.String(), "--allow-unrelated-histories")
	} else {
		cmd = exec.Command("git", "pull", remoteName, branch.String())
	}
	output, err := cmd.CombinedOutput()
	return string(output[:len(output)-1]), err
}

func (b *Bash) GenerateMocks() (string, error) {
	cmd := exec.Command("go", "generate", "./...")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (b *Bash) TestsPresent() (bool, error) {
	cmd := exec.Command("find", ".", "-name", "*.test.go", "-or", "-name", "*_test.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return len(output) > 0, err
}

func (b *Bash) RunTests() (string, error) {
	cmd := exec.Command("go", "test", "./...")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// TODO -> shift this to git repo
func (b *Bash) Push(branch model.Branch, withUpStream bool) (string, error) {
	var cmd *exec.Cmd
	if withUpStream {
		cmd = exec.Command("git", "push", "-u", "origin", branch.String())
	} else {
		cmd = exec.Command("git", "push", "origin", branch.String())
	}
	output, err := cmd.CombinedOutput()
	return string(output[:len(output)-1]), err
}

func (b *Bash) Exists(name, path string) bool {
	fpath := filepath.Join(path, name)
	_, err := os.Stat(fpath)
	return !errors.Is(err, os.ErrNotExist)
}

func (b *Bash) CreateFile(name, path string) (*os.File, error) {
	fpath := filepath.Join(path, name)
	return os.Create(fpath)
}

func (b *Bash) CreateDir(name, path string) error {
	dpath := filepath.Join(path, name)
	return os.Mkdir(dpath, os.ModePerm)
}

func (b *Bash) SetUpstream(remoteName string, branch model.Branch) error {
	remoteArg := fmt.Sprintf("%s/%s", remoteName, branch)
	cmd := exec.Command("git", "branch", "--set-upstream-to", remoteArg)
	_, err := cmd.CombinedOutput()
	return err
}

func (b *Bash) GenerateSSHKey(keyName, path, mail, passphrase string) error {
	filePath := filepath.Join(path, keyName)
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", mail, "-f", filePath, "-P", passphrase)
	_, err := cmd.CombinedOutput()
	return err
}

func (b *Bash) ShowFileContent(filename, path string) (string, error) {
	filePath := filepath.Join(path, filename)
	cmd := exec.Command("cat", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output[:len(output)-1]), nil
}
