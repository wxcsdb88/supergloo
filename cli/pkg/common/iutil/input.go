package iutil

import (
	"gopkg.in/AlecAivazis/survey.v1"
)

func GetStringInput(msg string, value *string) error {
	prompt := &survey.Input{Message: msg}
	if err := survey.AskOne(prompt, value, nil); err != nil {
		return err
	}
	return nil
}
