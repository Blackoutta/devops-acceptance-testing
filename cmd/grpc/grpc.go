package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/scenarios/grpc"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/assertion"
	"os"
)

func main() {
	var successCount int
	success := true
	testNum := 1
	exitChan := make(chan assertion.TestResult)

	go func() {
		grpc.RunGrpcTest(exitChan)
	}()

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
