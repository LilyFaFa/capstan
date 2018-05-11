/*
Copyright (c) 2018 The ZJU-SEL Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package push

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	propush "github.com/prometheus/client_golang/prometheus/push"
)

func Push(result string, endpoint string) error {
	entries := strings.Split(result, "\\n")
	for _, entry := range entries {
		item := strings.Split(entry, " ")
		if len(item) < 3 {
			return errors.New("The data format does not meet the requirements.")
		}
		err := PushEntry(item, endpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

func PushEntry(item []string, endpoint string) error {
	pushData := item[2]
	if len(item) > 3 {
		for i := 3; i < len(item); i++ {
			pushData += " "
			pushData += item[i]
		}
	}

	collections, err := ConvertToCollection(pushData)
	if err != nil {
		return err
	}
	fmt.Println(collections)
	// Get the job name
	jobName, ok := collections["job"]
	if !ok {
		jobName = "Unknown job"
	}

	helpMessage := fmt.Sprintf("The result of job %s ,named %s \n", jobName, item[0])
	tpmc := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: item[0],
		Help: helpMessage,
	})

	result, err := strconv.ParseFloat(item[1], 64)
	if err != nil {
		return err
	}
	tpmc.Set(result)

	// Push the data to pushgateway
	if err := propush.Collectors(
		jobName,
		collections,
		endpoint,
		tpmc,
	); err != nil {
		return err
	}
	return nil
}

func ConvertToCollection(pushData string) (map[string]string, error) {
	collections := make(map[string]string)
	data := strings.Split(pushData, ",")
	for _, value := range data {
		str := strings.Split(value, "=")
		if len(str) != 2 {
			return nil, errors.New("The data format does not meet the requirements.")
		}
		collections[str[0]] = str[1]
	}
	return collections, nil
}
