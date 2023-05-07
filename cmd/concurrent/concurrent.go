package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/grpc"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/pods"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/vm"

	"github.com/common-nighthawk/go-figure"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/artifact"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/build"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/pipeline"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/assertion"
)

func main() {
	success := true
	exitChan := make(chan assertion.TestResult)
	var successCount int

	// Just add test suites here
	var suiteList []func(chan assertion.TestResult)

	env := os.Getenv("CONFIGFILE")

	switch env {
	case "config-test.json", "/cf/config-test.json":
		suiteList = []func(chan assertion.TestResult){
			grpc.RunGrpcTest,
			pods.RunPodsTest,
			build.RunBuildAndDeployTest,
			pipeline.RunPipelineTest,
			vm.RunVmTest,
			artifact.RunArtifactTest,
		}
	case "config-prod.json", "/cf/config-prod.json":
		suiteList = []func(chan assertion.TestResult){
			pods.RunPodsTest,
			build.RunBuildAndDeployTest,
			pipeline.RunPipelineTest,
			artifact.RunArtifactTest,
		}
	}

	testNum := len(suiteList)

	for _, f := range suiteList {
		go f(exitChan)
		time.Sleep(15 * time.Second)
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
