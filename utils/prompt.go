package utils

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func Prompt(label string, opts ...interface{}) (res string, err error) {
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

	res, err = prompt.Run()
	return
}
