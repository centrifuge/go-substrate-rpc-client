package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func main() {
	ctx := context.Background()
	args := os.Args

	if len(args) != 4 {
		fmt.Println("Expected 4 args")
		os.Exit(1)
	}

	dirName := args[1]
	metaFileName := args[2]
	storageFileName := args[3]

	fmt.Printf("Test data directory name - %s\n", dirName)
	fmt.Printf("Test data meta file name - %s\n", metaFileName)
	fmt.Printf("Test data storage file name - %s\n", storageFileName)

	parser := NewParser()

	t := reflect.TypeOf(types.EventRecords{})

	requestMap := make(map[ClientOpts][]*ReqData)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if parser.CanSkip(field) {
			fmt.Println("Skipping field -", field.Name)
			continue
		}

		fieldInfo, err := parser.ParseField(field)

		if err != nil {
			fmt.Println("Couldn't parse field -", err)
			os.Exit(1)
		}

		if reqData, ok := requestMap[*fieldInfo.ClientOpts]; ok {
			reqData = append(reqData, fieldInfo.ReqData)
			requestMap[*fieldInfo.ClientOpts] = reqData
		} else {
			requestMap[*fieldInfo.ClientOpts] = []*ReqData{fieldInfo.ReqData}
		}
	}

	var rawData []*TestData

	for clientOpts, reqData := range requestMap {
		c, err := NewClient(clientOpts)

		if err != nil {
			fmt.Println("Couldn't create new client -", err)
			os.Exit(1)
		}

		for _, reqDatum := range reqData {
			rawDatum, err := c.GetTestData(ctx, reqDatum)

			if err != nil {
				fmt.Printf("Couldn't get raw data for %s %s - %s\n", reqDatum.Module, reqDatum.Call, err)

				continue
			}

			rawData = append(rawData, rawDatum)
		}
	}

	if err := writeTestFiles(rawData, dirName, metaFileName, storageFileName); err != nil {
		fmt.Println("Couldn't create test files -", err)
		os.Exit(1)
	}
}

func writeTestFiles(rawData []*TestData, dirName string, metaFileName string, storageFileName string) error {
	if err := os.RemoveAll(dirName); err != nil {
		return err
	}

	for i, rawDatum := range rawData {
		if err := os.MkdirAll(fmt.Sprintf("%s/%d", dirName, i), os.ModePerm); err != nil {
			return err
		}

		metaFile, err := os.Create(fmt.Sprintf("%s/%d/%s", dirName, i, metaFileName))

		if err != nil {
			return err
		}

		defer metaFile.Close()

		if _, err := metaFile.Write(rawDatum.Meta); err != nil {
			return err
		}

		storageFile, err := os.Create(fmt.Sprintf("%s/%d/%s", dirName, i, storageFileName))

		if err != nil {
			return err
		}

		defer storageFile.Close()

		if _, err := storageFile.Write(rawDatum.StorageData); err != nil {
			return err
		}

		infoFile, err := os.Create(fmt.Sprintf("%s/%d/info.txt", dirName, i))

		if err != nil {
			return err
		}

		defer infoFile.Close()

		infoStr := fmt.Sprintf(
			"Blockchain - %s\nBlock number - %d\nApi URL - %s\nWs URL - %s\n",
			rawDatum.Blockchain,
			rawDatum.BlockNumber,
			rawDatum.ApiURL,
			rawDatum.WsURL,
		)

		if _, err := infoFile.WriteString(infoStr); err != nil {
			return err
		}
	}

	return nil
}
