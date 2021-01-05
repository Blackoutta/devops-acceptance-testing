package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Blackoutta/profari"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/scenarios/purge"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/prep"
)

var teardown bool

func init() {
	flag.BoolVar(&teardown, "skip-teardown", false, "prevent the test suite from tearing down")
	flag.Parse()
}

func main() {
	t1 := &purge.PurgeTest{
		Name:         "资源清理套件",
		ErrChan:      make(chan error),
		SkipTeardown: prep.CheckSkipTeardown(),
	}

	result, exit := profari.RunTests(t1)
	fmt.Println(result)
	os.Exit(exit)
}
