// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func main() { //nolint: funlen
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

	if err := os.RemoveAll(dirName); err != nil {
		fmt.Println("Couldn't remove test dir -", err)

		os.Exit(1)
	}

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

	count := 0

	for clientOpts, reqData := range requestMap {
		c, err := NewClient(clientOpts)

		if err != nil {
			fmt.Println("Couldn't create new client -", err)
			os.Exit(1)
		}

		for _, reqDatum := range reqData {
			testData, err := c.GetTestData(ctx, reqDatum)

			if err != nil {
				fmt.Printf("Couldn't get raw data for %s %s - %s\n", reqDatum.Module, reqDatum.Call, err)

				continue
			}

			if err := writeTestData(testData, dirName, metaFileName, storageFileName, count); err != nil {
				fmt.Println("Couldn't create test files -", err)
				os.Exit(1)
			}

			count++
		}
	}
}

func writeTestData(
	testData *TestData,
	dirName string,
	metaFileName string,
	storageFileName string,
	count int,
) error {
	if err := os.MkdirAll(fmt.Sprintf("%s/%d", dirName, count), os.ModePerm); err != nil {
		return err
	}

	metaFile, err := os.Create(fmt.Sprintf("%s/%d/%s", dirName, count, metaFileName))

	if err != nil {
		return err
	}

	defer metaFile.Close()

	if _, err := metaFile.Write(testData.Meta); err != nil {
		return err
	}

	storageFile, err := os.Create(fmt.Sprintf("%s/%d/%s", dirName, count, storageFileName))

	if err != nil {
		return err
	}

	defer storageFile.Close()

	if _, err := storageFile.Write(testData.StorageData); err != nil {
		return err
	}

	infoFile, err := os.Create(fmt.Sprintf("%s/%d/info.txt", dirName, count))

	if err != nil {
		return err
	}

	defer infoFile.Close()

	infoStr := fmt.Sprintf(
		"Blockchain - %s\nBlock number - %d\nApi URL - %s\nWs URL - %s\n",
		testData.Blockchain,
		testData.BlockNumber,
		testData.APIURL,
		testData.WsURL,
	)

	if _, err := infoFile.WriteString(infoStr); err != nil {
		return err
	}

	return nil
}
