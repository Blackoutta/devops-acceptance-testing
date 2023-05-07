package main

import (
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/pipeline"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/assertion"
)

func main() {
	success := true
	exitChan := make(chan assertion.TestResult)
	var successCount int

	// Just add test suites here
	suiteList := []func(chan assertion.TestResult){
		pipeline.RunPipelineTest,
	}

	testNum := len(suiteList)

	for _, f := range suiteList {
		go f(exitChan)
	}

	var tr []assertion.TestResult

	for i := 0; i < testNum; i++ {
		var r assertion.TestResult
		if r = <-exitChan; r.Result != true {
			success = false
		} else {
			successCount++
		}
		tr = append(tr, r)
	}

	figure.NewFigure(fmt.Sprintf("SUCCEEDED TESTS: %v/%v", successCount, testNum), "", false).Print()

	assertion.PrintResult(tr)

	if success == false {
		os.Exit(1)
	}
}
