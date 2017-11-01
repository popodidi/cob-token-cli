package private

import (
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"io/ioutil"
	"strings"
	"path"
	"encoding/csv"
	"bufio"
	"strconv"
	"fmt"
	"errors"
	"gopkg.in/cheggaaa/pb.v2"
	"github.com/popodidi/cob-token-cli/utils"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

type toSend struct {
	address string
	value   float64
}

func allocateCOBAction(c *cli.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var files []os.FileInfo
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var csvFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".csv") {
			csvFiles = append(csvFiles, f.Name())
		}
	}

	csvFileName := ""
	csvFilePrompt := &survey.Select{
		Message: "Choose a .csv file",
		Options: csvFiles,
	}
	survey.AskOne(csvFilePrompt, &csvFileName, nil)

	var csvFile *os.File
	csvFile, err = os.Open(path.Join(dir, csvFileName))

	r := csv.NewReader(bufio.NewReader(csvFile))

	var result [][]string
	result, err = r.ReadAll()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	toSends := make([]toSend, 0)
	for i, row := range result {
		if len(row) != 2 {
			return cli.NewExitError(err.Error(), 1)
		}
		if i == 0 {
			isValidTitle := row[0] == "address" && row[1] == "value"
			if isValidTitle {
				continue
			} else {
				return cli.NewExitError(err.Error(), 1)
			}
		} else {
			addr := row[0]
			value, _err := strconv.ParseFloat(row[1], 64)
			if _err != nil {
				return cli.NewExitError(_err.Error(), 1)
			}
			toSends = append(toSends, toSend{addr, value})
		}
	}

	var totalValue float64 = 0
	for _, s := range toSends {
		totalValue += s.value
	}
	confirmTotalValue := false
	confirmTotalValuePrompt := &survey.Confirm{
		Message: fmt.Sprintf("Total count: %d / Total value: %f COBs", len(toSends), totalValue),
	}
	survey.AskOne(confirmTotalValuePrompt, &confirmTotalValue, nil)

	if !confirmTotalValue {
		return cli.NewExitError(errors.New("user stopped"), 1)
	}

	var qs = []*survey.Question{
		{
			Name:     "from-private-key",
			Prompt:   &survey.Password{Message: "From private key"},
			Validate: survey.Required,
		},
		{
			Name:     "gas-price",
			Prompt:   &survey.Input{Message: "Gas Price (Gwei)"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		FromPrivKey string `survey:"from-private-key"`
		GasPrice    int64  `survey:"gas-price"`
	}{}

	err = survey.Ask(qs, &answers)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	shouldStart := false
	shouldStartPrompt := &survey.Confirm{
		Message: "Start sending COBs",
	}
	survey.AskOne(shouldStartPrompt, &shouldStart, nil)

	if !shouldStart {
		return cli.NewExitError(errors.New("user stopped"), 1)
	}

	var logs = [][]string{[]string{"address", "value", "tx"}}

	gasPrice := big.NewInt(1)
	gasPrice.Mul(big.NewInt(answers.GasPrice), big.NewInt(1000000000))
	count := len(toSends)
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		cobValue := decimal.NewFromFloat(toSends[i].value)
		cobValue = cobValue.Mul(decimal.New(1, 18))
		cobAmount := big.NewInt(cobValue.IntPart())

		_tx, _err := utils.SendCOB(answers.FromPrivKey, toSends[i].address, cobAmount, big.NewInt(500000), gasPrice)
		if _err != nil {
			logs = append(logs, []string{toSends[i].address, fmt.Sprintf("%f", toSends[i].value), "ERROR"})
		}
		logs = append(logs, []string{toSends[i].address, fmt.Sprintf("%f", toSends[i].value), _tx.Hash().Hex()})
		bar.Increment()
	}
	bar.Finish()

	var logFile *os.File
	logFileName := csvFileName + "." + fmt.Sprint(time.Now().Unix()) + ".log"
	logFile, err = os.Create(path.Join(dir, logFileName))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer logFile.Close()

	writer := csv.NewWriter(logFile)
	defer writer.Flush()

	for _, value := range logs {
		err = writer.Write(value)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	return nil
}
