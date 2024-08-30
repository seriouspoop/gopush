package utils

import (
	"github.com/manifoldco/promptui"
)

func Prompt(label string) (res string, err error) {
	template := &promptui.PromptTemplates{
		Valid:   "{{ . }}: ",
		Success: "{{ `\U00002714` | green }} {{ . | faint}}{{ `:` | faint}} ",
	}
	prompt := promptui.Prompt{
		Label:     label,
		Templates: template,
	}

	res, err = prompt.Run()
	return
}
