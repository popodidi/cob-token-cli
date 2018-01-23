package private

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"
	"time"

	"github.com/popodidi/cob-token-cli/utils"

	"github.com/shopspring/decimal"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v2"
)

type toSend struct {
	address string
	value   string
}

type sendLog struct {
	toSend    toSend
	txHash    string
	timestamp int64
	error     error
}

func (l *sendLog) Keys() []string {
	return []string{"address", "value", "tx", "timestamp", "error"}
}

func (l *sendLog) Value(key string) string {
	switch key {
	case "address":
		return l.toSend.address
	case "value":
		return l.toSend.value
	case "tx":
		return l.txHash
	case "timestamp":
		return fmt.Sprintf("%d", l.timestamp)
	case "error":
		if l.error != nil {
			return l.error.Error()
		} else {
			return ""
		}
	default:
		return ""
	}
}

func allocateCOBAction(c *cli.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var csvFileName string
	var csvData [][]string
	csvFileName, csvData, err = utils.SelectCSV(dir, "Choose a .csv file")
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
		fmt.Sprintf("Total count: %d / Total value: %s COBs / Err count: %d",
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

	var logs = []sendLog{}

	count := len(toSends)
	bar := pb.StartNew(count)
	updateLogsAndBar := func(_l sendLog) {
		logs = append(logs, _l)
		bar.Increment()
	}

	for i := 0; i < count; i++ {
		var log sendLog
		var cobAmount *big.Int
		var _err error
		cobAmount, _err = utils.StringToWei(toSends[i].value)

		if _err != nil {
			log = sendLog{
				toSend:    toSends[i],
				txHash:    "",
				timestamp: time.Now().Unix(),
				error:     _err,
			}
			updateLogsAndBar(log)
			continue
		}
		_tx, _err := utils.SendCOB(privateKey, toSends[i].address, cobAmount, big.NewInt(500000), gasPrice)

		if _err != nil {
			log = sendLog{
				toSend:    toSends[i],
				txHash:    "",
				timestamp: time.Now().Unix(),
				error:     _err,
			}
			updateLogsAndBar(log)
			continue
		}
		log = sendLog{
			toSend:    toSends[i],
			txHash:    _tx.Hash().Hex(),
			timestamp: time.Now().Unix(),
			error:     _err,
		}
		updateLogsAndBar(log)
	}
	bar.Finish()

	writeLogsToFile(logs, dir, csvFileName)
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

func writeLogsToFile(sendLogs []sendLog, dir string, csvName string) error {
	if len(sendLogs) == 0 {
		return errors.New("no logs to write")
	}

	timeStr := fmt.Sprintf("%d%d%d%d%d%d",
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	logFilePath := timeStr + ".log." + csvName
	err := writeLogToCsv(sendLogs, path.Join(dir, logFilePath))
	if err != nil {
		fmt.Printf("failed to write log file\n%s", err.Error())
	}
	return err
}

func writeLogToCsv(logs []sendLog, filePath string) error {
	if len(logs) == 0 {
		return errors.New("no logs to write")
	}

	csvData := make([][]string, 0)
	keys := logs[0].Keys()
	header := make([]string, 0)
	for _, k := range keys {
		header = append(header, k)
	}
	csvData = append(csvData, header)

	for _, l := range logs {
		row := make([]string, 0)
		for _, k := range keys {
			row = append(row, l.Value(k))
		}
		csvData = append(csvData, row)
	}

	return utils.WriteDataToCsv(csvData, filePath)
}
