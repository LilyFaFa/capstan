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

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ZJU-SEL/capstan/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/pflag"
)

var (
	PGWEndpoint = pflag.String("PGWEndpoint", "", "Path to the pushGateWay, You can also provide this variable as an environment variable.")
)

func main() {
	util.InitFlags()

	if len(pflag.Args()) != 4 {
		fmt.Fprintf(os.Stderr, "%s\n", "The length of the data does not meet the requirements")
		os.Exit(0)
	}

	if *PGWEndpoint == "" && os.Getenv("PGWEndpoint") == "" {
		fmt.Fprintf(os.Stderr, "%s\n", "No env or flag named PGWEndpoint found")
		os.Exit(0)
	} else if *PGWEndpoint == "" {
		*PGWEndpoint = os.Getenv("PGWEndpoint")
	}

	pushData := pflag.Args()[0:2]
	pushData = append(pushData, pflag.Args()[2]+" "+pflag.Args()[3])
	data := strings.Split(pushData[2], ",")
	collections := make(map[string]string)
	for _, value := range data {
		str := strings.Split(value, "=")
		if len(str) != 2 {
			fmt.Fprintf(os.Stderr, "%s\n", "The data format does not meet the requirements.")
			os.Exit(0)
		}
		collections[str[0]] = str[1]
	}

	// Get the job name
	jobName, ok := collections["job"]
	if !ok {
		jobName = "Unknown job"
	}
	helpMessage := fmt.Sprintf("The result of job %s ,named %s \n", jobName, pushData[0])

	// Parse result
	result, err := strconv.ParseFloat(pushData[1], 64)
	tpmc := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: pushData[0],
		Help: helpMessage,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(0)
	}
	tpmc.Set(result)

	// Push the data to pushgateway
	if err := push.Collectors(
		jobName,
		collections,
		*PGWEndpoint,
		tpmc,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(0)
	}

}
