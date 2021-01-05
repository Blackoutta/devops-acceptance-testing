package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Blackoutta/profari"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/scenarios/test"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/prep"
)

var teardown = flag.Bool("skip-teardown", false, "prevent the test suite from tearing down")

func init() {
	flag.Parse()
}

func main() {
	t1 := &test.TestFeatureTest{
		Name:         "测试功能测试套件",
		ErrChan:      make(chan error),
		SkipTeardown: prep.CheckSkipTeardown(),
	}

	result, exit := profari.RunTests(t1)
	fmt.Println(result)
	os.Exit(exit)
}
