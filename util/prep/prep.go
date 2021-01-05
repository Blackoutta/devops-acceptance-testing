package prep

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/assertion"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/param"
)

func SetupTest(suiteName string) (*os.File, *assertion.Assertion, *param.SuiteParams, http.Client) {
	// 日志目录
	logpath := "./logs"
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.Mkdir(logpath, os.ModePerm|os.ModeDir)
	}

	// 设置日志
	f, err := os.Create(logpath + "/" + suiteName + ".log")
	if err != nil {
		log.Fatalf("error creating log file: %v\n", err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	logger := log.New(mw, suiteName+" ", 7)

	// 初始化断言
	ast := assertion.NewAssertion(suiteName, logger)

	// 初始化测试参数
	sp := param.NewParamSet()

	// 初始化http客户端
	c := http.Client{
		Timeout: 35 * time.Second,
	}

	return f, ast, sp, c
}

func CheckSkipTeardown() bool {
	var skip bool
	flag.Visit(func(f *flag.Flag) {
		fmt.Println(f.Name)
		if f.Name == "skip-teardown" {
			skip = true
		}
	})
	return skip
}
