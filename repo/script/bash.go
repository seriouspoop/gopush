package script

import (
	"errors"
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
	return string(output[:len(output)-1]), err
}

func (b *Bash) Exists(path, name string) bool {
	fpath := filepath.Join(path, name)
	_, err := os.Stat(fpath)
	return !errors.Is(err, os.ErrNotExist)
}

func (b *Bash) CreateFile(path, name string) (*os.File, error) {
	fpath := filepath.Join(path, name)
	return os.Create(fpath)
}

func (b *Bash) CreateDir(path, name string) error {
	dpath := filepath.Join(path, name)
	return os.Mkdir(dpath, os.ModePerm)
}

func (b *Bash) GenerateSSHKey(path, keyName, mail, passphrase string) error {
	filePath := filepath.Join(path, keyName)
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", mail, "-f", filePath, "-P", passphrase)
	_, err := cmd.CombinedOutput()
	return err
}

func (b *Bash) PullMerge() (string, error) {
	cmd := exec.Command("git", "merge")
	output, err := cmd.CombinedOutput()
	return string(output[:len(output)-1]), err
}
