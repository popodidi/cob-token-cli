package utils

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

func SelectCSV(dir string, message string) (fileName string, data [][]string, err error) {

	var files []os.FileInfo
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		fileName = ""
		data = nil
		return
	}

	var csvFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".csv") {
			csvFiles = append(csvFiles, f.Name())
		}
	}

	if len(csvFiles) == 0 {
		fileName = ""
		data = nil
		err = errors.New("no .csv file found")
		return
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
		fileName = csvFileName
		data = nil
		return
	}

	r := csv.NewReader(bufio.NewReader(csvFile))

	var result [][]string
	result, err = r.ReadAll()
	if err != nil {
		fileName = csvFileName
		data = nil
		return
	}

	fileName = csvFileName
	data = result
	return
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
