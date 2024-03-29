package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/common"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/prep"
)

var teardown = flag.Bool("skip-teardown", false, "prevent the test suite from tearing down")

func init() {
	flag.Parse()
}

func main() {
	t1 := &common.CommonCoreTest{
		Name:         "Core模块通用测试套件",
		ErrChan:      make(chan error),
		SkipTeardown: prep.CheckSkipTeardown(),
	}

	result, exit := profari.RunTests(t1)
	fmt.Println(result)
	os.Exit(exit)
}
