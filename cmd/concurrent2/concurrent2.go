package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/common"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/order"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/probe"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/scenarios/test"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/prep"
)

var teardown = flag.Bool("skip-teardown", false, "prevent the test suite from tearing down")

func init() {

	flag.Parse()
}

func main() {
	t1 := &probe.ProbeTest{
		Name:         "探针功能测试套件",
		ErrChan:      make(chan error),
		SkipTeardown: prep.CheckSkipTeardown(),
	}

	t2 := &common.CommonCoreTest{
		Name:         "Core模块通用测试套件",
		ErrChan:      make(chan error),
		SkipTeardown: prep.CheckSkipTeardown(),
	}

	t3 := &test.TestFeatureTest{
		Name:         "测试功能测试套件",
		ErrChan:      make(chan error),
		SkipTeardown: prep.CheckSkipTeardown(),
	}

	t4 := &order.OrderTest{
		Name:         "工单测试套件",
		ErrChan:      make(chan error),
		SkipTeardown: prep.CheckSkipTeardown(),
	}

	result, exit := profari.RunTests(t1, t2, t3, t4)
	fmt.Println(result)
	os.Exit(exit)
}
