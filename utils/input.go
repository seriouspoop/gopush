package utils

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
)

func valid(s string) error {
	if len(s) <= 0 {
		return errors.New("prompt cannot be empty")
	}
	return nil
}

func Prompt(password bool, allowEmpty bool, label string, opts ...interface{}) (res string, err error) {
	template := &promptui.PromptTemplates{
		Valid:   "{{ . }}: ",
		Success: "{{ `\U00002714` | green }} {{ . | faint}}{{ `:` | faint}} ",
	}

	if len(opts) > 0 {
		label = fmt.Sprintf(label, opts...)
	}

	prompt := promptui.Prompt{
		Label:     label,
		Templates: template,
	}

	if !allowEmpty {
		prompt.Validate = valid
	}

	if password {
		prompt.Mask = '*'
	}

	res, err = prompt.Run()
	return
}

func Select(items []string) (res string, err error) {
	selectTemplate := &promptui.SelectTemplates{
		Active: "\U0001F892 {{ . | green }}",
	}
	prompt := promptui.Select{
		Label:     "Select commit type",
		Items:     items,
		Templates: selectTemplate,
	}
	_, res, err = prompt.Run()
	return
}
