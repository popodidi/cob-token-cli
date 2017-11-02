package private

import (
	"github.com/urfave/cli"
	"os"
	"path"
	"fmt"
	"errors"
	"gopkg.in/cheggaaa/pb.v2"
	"github.com/popodidi/cob-token-cli/utils"
	"math/big"
	"time"
	"github.com/shopspring/decimal"
)

type toSend struct {
	address string
	value   string
}

func allocateCOBAction(c *cli.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var csvData [][]string
	csvData, err = utils.SelectCSV(dir, "Choose a .csv file")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var toSends []toSend
	toSends, err = readFromCsvData(csvData)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	totalValue := decimal.NewFromFloat(0)
	errCount := 0
	for _, s := range toSends {
		_amount, _err := utils.StringToWei(s.value)
		if _err != nil {
			errCount += 1
		} else {
			totalValue = totalValue.Add(decimal.NewFromBigInt(_amount, 0).Div(decimal.New(1, 18)))
		}
	}

	if !utils.AskForConfirm(
		fmt.Sprintf("Total count: %d / Total value: %s COBs / Err coune: %d",
			len(toSends)-errCount, totalValue.String(), errCount)) {
		return cli.NewExitError(errors.New("user stopped"), 1)
	}

	var privateKey string
	var gasPrice *big.Int
	privateKey, err = utils.AskForPrivateKey()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	gasPrice, err = utils.AskForGasPriceGwei()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if !utils.AskForConfirm("START") {
		return cli.NewExitError(errors.New("user stopped"), 1)
	}

	var logs = [][]string{[]string{"address", "value", "tx"}}
	defer writeLogsToFile(logs, dir)

	count := len(toSends)
	bar := pb.StartNew(count)
	updateLogsAndBar := func(_l []string) {
		logs = append(logs, _l)
		bar.Increment()
	}

	for i := 0; i < count; i++ {
		var log []string
		var cobAmount *big.Int
		cobAmount, err = utils.StringToWei(toSends[i].value)

		if err != nil {
			log = []string{toSends[i].address, fmt.Sprintf("%f", toSends[i].value), "ERROR"}
			updateLogsAndBar(log)
			continue
		}
		_tx, _err := utils.SendCOB(privateKey, toSends[i].address, cobAmount, big.NewInt(500000), gasPrice)

		if _err != nil {
			log = []string{toSends[i].address, fmt.Sprintf("%f", toSends[i].value), "ERROR"}
			updateLogsAndBar(log)
			continue
		}
		log = []string{toSends[i].address, fmt.Sprintf("%f", toSends[i].value), _tx.Hash().Hex()}
		updateLogsAndBar(log)
	}
	bar.Finish()

	return nil
}

func readFromCsvData(csvData [][]string) ([]toSend, error) {
	toSends := make([]toSend, 0)
	for i, row := range csvData {
		if len(row) != 2 {
			return nil, errors.New("invalid column number (must be 2)")
		}
		if i == 0 {
			isValidTitle := row[0] == "address" && row[1] == "value"
			if isValidTitle {
				continue
			} else {
				return nil, errors.New("invalid column title (must be \"address\",\"value\")")
			}
		} else {
			addr := row[0]
			value := row[1]
			toSends = append(toSends, toSend{addr, value})
		}
	}
	return toSends, nil
}

func writeLogsToFile(logs [][]string, dir string) error {
	logFilePath := "log." + fmt.Sprint(time.Now().Unix()) + ".csv"
	err := utils.WriteDataToCsv(logs, path.Join(dir, logFilePath))
	if err != nil {
		fmt.Printf("log file written to %s", logFilePath)
	}
	return err
}
