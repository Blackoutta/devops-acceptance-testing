package main

import (
	"flag"
	"fmt"
	"github.com/Blackoutta/profari"
	"github.com/common-nighthawk/go-figure"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/scenarios/build"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/assertion"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/prep"
	"os"
)

var teardown = flag.Bool("skip-teardown", false, "prevent the test suite from tearing down")

func init() {
	flag.Parse()
}

func main() {
	var successCount int
	success := true
	testNum := 1
	exitChan := make(chan assertion.TestResult)

	go func() {
		build.RunBuildAndDeployTest(exitChan)
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

func newApi() {
	t1 := &build.BuildTest{
		Name:         "简单构建测试套件",
		ErrChan:      make(chan error),
		SkipTeardown: prep.CheckSkipTeardown(),
	}
	result, exit := profari.RunTests(t1)

	fmt.Println(result)
	os.Exit(exit)
}
