package probe

import (
	"fmt"
	"os"
	"time"

	"github.com/Blackoutta/profari"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/API"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/newapi"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/random"
)

type ProbeTest struct {
	Name         string
	ErrChan      chan error
	SkipTeardown bool
	*profari.Client
	logFile *os.File
	*newapi.SetupParams
}

func (t *ProbeTest) GetName() string {
	return t.Name
}

func (t *ProbeTest) GetErrChan() chan error {
	return t.ErrChan
}

func (t *ProbeTest) Run() {
	var err error
	// You must initialize the profari client before test starts
	t.Client, t.logFile, err = profari.NewClient(t.Name, t.ErrChan)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 开始

	// 准备工作
	t.SetupParams = newapi.BasicSetup(t.Client)
	fmt.Printf("%+v\n", t.SetupParams)

	// 编辑部署，加入readiness + liveness HTTP
	ingressHost := fmt.Sprintf("devops.%s.com", random.ShortGUID())
	t.Send(newapi.NewEditK8sDeployConfig(t.DeployID, ingressHost)).AssertContainString("编辑部署，加入http get类型的readiness和liveness probe", t.Resp, "success")

	// 执行部署
	var gen API.GeneralResp
	t.Send(
		newapi.ExecuteDeploy{DeployId: t.DeployID},
	).DecodeJSON(&gen).AssertContainString("执行部署应成功", t.Resp, "success")
	deployHistoryID := gen.Data

	statusChan := make(chan string)

	// 检查部署状态应成功
	go func() {
		var history API.DeployHistory

		t.Send(newapi.GetDeployHistory{
			DeployHistoryID: deployHistoryID,
		}).DecodeJSON(&history).AssertContainString("查询部署状态接口应成功", t.Resp, "success")

		for {
			t.Send(newapi.GetDeployHistory{
				DeployHistoryID: deployHistoryID,
			}).DecodeJSON(&history)

			status := history.Data.Status
			statusChan <- status

			time.Sleep(2 * time.Second)
		}
	}()

	var finish bool

	for !finish {
		select {
		case s := <-statusChan:
			switch s {
			case "SUCCESS":
				t.Printf("部署成功！部署状态为: %s", s)
				finish = true
			case "FAILED":
				t.Printf("部署失败！部署状态为: %s", s)
				t.ErrChan <- fmt.Errorf("部署失败！部署状态为: %s", s)
				finish = true
			case "WAITING", "RUNNING":
				t.Printf("部署进行中，部署状态为: %s", s)
			default:
				t.Printf("遇到未知的部署状态: %s, 部署异常，请检查！", s)
				finish = true
			}
		case <-time.After(10 * time.Second):
			t.Println("长时间未获取到部署状态信息，测试失败，请检查接口。")
			t.ErrChan <- fmt.Errorf("长时间未获取到部署状态信息，测试失败，请检查接口。")
			finish = true
		}
	}

	// 检查K8S描述文件
	var sd API.StringData
	t.Send(
		newapi.GetK8SDeploymentFile{
			DeployID: t.DeployID,
			AppID:    t.AppID,
			EnvID:    t.EnvID,
		},
	).DecodeJSON(&sd).AssertEqualInt("查看Pod的k8s deployment接口调用应成功", 0, gen.ErrorCode)

	t.AssertContainString("k8s deployment文件中含有readiness的信息", sd.Data, "readiness")
	t.AssertContainString("k8s deployment文件中含有liveness的信息", sd.Data, "liveness")

	// 获取pod名
	var pods API.PodList
	t.Send(
		newapi.GetPodList{
			DeployID: t.DeployID,
			AppID:    t.AppID,
			EnvID:    t.EnvID,
		},
	).DecodeJSON(&pods).AssertEqualInt("获取Pod名", 1, len(pods.KubernetesDeploy))
	podName := pods.KubernetesDeploy[0].PodInstances[0].PodName

	// 检查describe
	t.Send(
		newapi.GetPodDescribeInfo{
			DeployID: t.DeployID,
			AppID:    t.AppID,
			EnvID:    t.EnvID,
			PodName:  podName,
		},
	).DecodeJSON(&sd).AssertEqualInt("获取pod的describe接口应调用成功", 0, sd.ErrorCode)

	t.AssertContainString("describe信息中显示pod启动状态正常", sd.Data, "Started container app")

}

func (t *ProbeTest) Teardown() {
	if t.SkipTeardown == true {
		t.Println("Skipping Teardown...")
		t.logFile.Close()
		t.EndTest()
		return
	}

	defer t.logFile.Close()
	defer t.EndTest()

	t.Send(
		newapi.DeleteDeploy{
			DeployID: t.DeployID,
		}).AssertContainString("删除部署应成功", t.Resp, "success")

	t.Send(newapi.DeleteK8sEnv{
		EnvID: t.EnvID,
	}).AssertContainString("删除k8s环境应成功", t.Resp, "success")

	t.Send(
		newapi.DeleteApp{
			Id:   t.AppID,
			Name: t.AppName,
		}).AssertContainString("删除应用应成功", t.Resp, "success")

	t.Send(newapi.DeleteProject{
		ProjectID: t.ProjectID,
	}).AssertContainString("删除项目应成功", t.Resp, "success")

}
