package utils

import (
	"gopkg.in/AlecAivazis/survey.v1"
	"math/big"
)

func AskForConfirm(message string) bool {
	result := false
	prompt := &survey.Confirm{
		Message: message,
	}
	survey.AskOne(prompt, &result, nil)
	return result
}

func AskForETHAddress() (string, error) {
	result := ""
	prompt := &survey.Input{
		Message: "ETH Address",
	}
	err := survey.AskOne(prompt, &result, survey.Required)
	if err != nil {
		return "", err
	}
	return result, nil
}

func AskForString(message string) (string, error) {
	result := ""
	prompt := &survey.Input{
		Message: message,
	}
	err := survey.AskOne(prompt, &result, survey.Required)
	if err != nil {
		return "", err
	}
	return result, nil
}

func AskForFloat(message string) (float64, error) {
	var result float64
	prompt := &survey.Input{
		Message: message,
	}
	err := survey.AskOne(prompt, &result, survey.Required)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func AskForPrivateKey() (string, error) {
	result := ""
	prompt := &survey.Password{
		Message: "From private key",
	}

	err := survey.AskOne(prompt, &result, survey.Required)
	if err != nil {
		return "", err
	}
	return result, nil
}

func AskForGasPriceGwei() (*big.Int, error) {
	var result int64 = 0
	prompt := &survey.Input{
		Message: "Gas Price (Gwei)",
	}
	err := survey.AskOne(prompt, &result, survey.Required)
	if err != nil {
		return nil, err
	}
	gasPrice := big.NewInt(1)
	gasPrice.Mul(big.NewInt(result), big.NewInt(1000000000))
	return gasPrice, nil

}
