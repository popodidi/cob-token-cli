package utils

import (
	"os"
	"io/ioutil"
	"strings"
	"gopkg.in/AlecAivazis/survey.v1"
	"path"
	"encoding/csv"
	"bufio"
	"errors"
)

func SelectCSV(dir string, message string) ([][]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var csvFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".csv") {
			csvFiles = append(csvFiles, f.Name())
		}
	}

	if len(csvFiles) == 0 {
		return nil, errors.New("no .csv file found")
	}

	csvFileName := ""
	csvFilePrompt := &survey.Select{
		Message: message,
		Options: csvFiles,
	}
	survey.AskOne(csvFilePrompt, &csvFileName, nil)

	var csvFile *os.File
	csvFile, err = os.Open(path.Join(dir, csvFileName))
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(bufio.NewReader(csvFile))

	var result [][]string
	result, err = r.ReadAll()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func WriteDataToCsv(csvData [][]string, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	for _, value := range csvData {
		err = writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}
