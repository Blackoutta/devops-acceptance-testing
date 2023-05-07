package test

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/API"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/newapi"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/conf"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/random"
)

type TestFeatureTest struct {
	Name         string
	ErrChan      chan error
	SkipTeardown bool
	*profari.Client
	logFile *os.File
	params
}

type params struct {
	projectID int
	envID     int
	testID    int
}

func (t *TestFeatureTest) GetName() string {
	return t.Name
}

func (t *TestFeatureTest) GetErrChan() chan error {
	return t.ErrChan
}

func (t *TestFeatureTest) Run() {
	var err error
	// You must initialize the profari client before test starts
	t.Client, t.logFile, err = profari.NewClient(t.Name, t.ErrChan)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 开始
	// 创建项目
	var gen API.GeneralResp
	projectName := "proj-" + random.ShortGUID()
	t.Send(newapi.CreateProject{
		DepartmentId: 5,
		Description:  "这是一个项目",
		Forbid:       false,
		Identify:     projectName,
		Name:         projectName,
		UseDefault:   false,
	}).DecodeJSON(&gen).AssertContainString("创建项目应成功", t.Resp, "success")

	// 获取projectID
	var projects API.Projects
	t.Send(newapi.GetProject{ProjectName: projectName}).DecodeJSON(&projects).AssertEqualInt("获取项目ID应成功", len(projects.Data.Data), 1)
	t.projectID = projects.Data.Data[0].ID
	t.Printf("Project ID is : %v\n", t.projectID)

	// 创建Docker凭证
	var dockerCredentials API.DockerCredentials
	t.Send(newapi.CreateDockerCredential{
		Description: "this is a docker credential",
		Name:        "docker_cred_" + random.ShortGUID(),
		Password:    "",
		ProjectId:   t.projectID,
		Type:        "DOCKER",
		UserName:    "jenkins",
	}).AssertContainString("创建Docker凭证应成功", t.Resp, "success")

	// 获取Docker Credential ID
	t.Send(newapi.GetCredential{
		Type:      "DOCKER",
		ProjectId: t.projectID,
		Name:      "docker",
	}).DecodeJSON(&dockerCredentials).AssertEqualInt("获取Docker凭证ID应成功", len(dockerCredentials.Data.Data), 1)
	dockerID := dockerCredentials.Data.Data[0].ID
	t.Printf("Docker credential ID is: %v\n", dockerID)

	// 创建K8S环境
	envName := "env-" + random.ShortGUID()
	namespace := "e-" + random.ShortGUID()
	t.Send(newapi.CreateK8sEnv{
		Config:     []byte(conf.ReadEnvFile()),
		Kubectx:    conf.EnvContext,
		Name:       envName,
		Namespace:  namespace,
		ServerAddr: "https://" + conf.ServerAddr + ":6443",
		Type:       "TEST",
		ProjectId:  t.projectID,
		RoomId:     2,
		ZoneId:     3,
	}).AssertContainString("创建k8s环境应成功", t.Resp, "success")

	// 获取环境ID
	var envs API.Environments
	t.Send(newapi.GetEnvs{
		ProjectId: t.projectID,
	}).DecodeJSON(&envs).AssertEqualInt("获取环境ID应成功", len(envs.Data.Data), 1)
	t.envID = envs.Data.Data[0].ID
	t.Printf("Env ID is: %v\n", t.envID)

	// 创建测试
	image := "hub.iot.chinamobile.com/offline/ecp-test:latest"

	t.Send(newapi.CreateTest{
		CredentialId:  dockerID,
		EnvId:         t.envID,
		Image:         image,
		ImageAuthType: "CREDENTIAL",
		Name:          "test_" + random.ShortGUID(),
		ProjectId:     t.projectID,
		Timeout:       300,
	}).DecodeJSON(&gen).AssertContainString("创建测试应成功", t.Resp, "success")
	t.testID = gen.Data

	// 编辑测试k8s配置
	t.Send(newapi.EditTestConfig{
		TestID: t.testID,
		ContainerLog: newapi.ContainerLog{
			LogPath: "",
		},
		HostMap: newapi.HostMap{
			Host: "10.12.4.9",
		},
		KubeConfig: newapi.TestKubeConfig{
			Cpu:       "1000m",
			Mem:       "512Mi",
			Namespace: namespace,
			IsChanged: true,
			Duplicate: 1,
		},
		VarMap: newapi.VarMap{
			CONFIGFILE: "config-test.json",
			SUITE:      "artifact",
		},
	}).AssertContainString("编辑测试k8s配置，加入Host和Env配置应成功", t.Resp, "success")

	// 编辑测试基本配置，加入参数
	t.Send(
		newapi.EditTestBase{
			CredentialId: dockerID,
			EnvId:        t.envID,
			GroupAuth: []newapi.GroupAuth{
				{
					Auth:     true,
					Copy:     true,
					Delete:   true,
					Edit:     true,
					GroupId:  0,
					Readonly: true,
					Run:      true,
				},
			},
			Id:    t.testID,
			Image: image,
			Name:  "edited_test",
			ParameterReqList: []newapi.ParameterReq{
				{
					DefaultValue: "test_default",
					Description:  "this is a param",
					Name:         "test_param",
					ParamValues:  "test_value",
					Required:     newapi.NOT_REQUIRED,
					Type:         newapi.STRING_PARAM,
				},
			},
			ProjectId: t.projectID,
			Timeout:   600,
		},
	).AssertContainString("编辑测试基本配置，加入参数", t.Resp, "success")

	// 查询测试详情
	var testDetail API.TestDetail
	t.Send(
		newapi.GetTestDetailByID{
			TestID: t.testID,
		}).DecodeJSON(&testDetail).AssertEqualInt("查询测试详情获取到的测试ＩＤ与实际测试ＩＤ一致", t.testID, testDetail.Data.Id)

	// 查询测试k8s详情
	var testConf API.TestConfig
	t.Send(
		newapi.GetTestK8SConfByID{TestID: t.testID}).DecodeJSON(&testConf)

	t.AssertContainString("测试k8s config中的环境变量与编辑时保存的一致", testConf.Data.VarMap.CONFIGFILE, "config-test.json")
	t.AssertContainString("测试k8s config中的DNS解析(host)与编辑时保存的一致", testConf.Data.HostMap.Host, "10.12.4.9")

	// 查询执行测试之前的参数配置
	var testParam API.TestParams
	t.Send(
		newapi.TestBeforeRun{
			TestID: t.testID,
		}).DecodeJSON(&testParam).AssertEqualInt("查询执行测试之前的参数配置应成功", 1, len(testParam.Data))

	// 执行测试
	plist := make([]newapi.Parameter, 0, 0)
	t.Send(newapi.ExecuteTest{
		ParameterList: plist,
		TestId:        t.testID,
	}).DecodeJSON(&gen).AssertContainString("执行测试应成功", t.Resp, "success")
	jobID := gen.Data
	t.Printf("Job ID is: %v\n", jobID)

	// 持续检查测试状态
	sChan := make(chan string, 1)

	go func() {
		var sl API.TestStatusList
		for {
			t.Send(
				newapi.GetTestStatus{
					IDs: strconv.Itoa(t.testID),
				},
			).DecodeJSON(&sl)

			if len(sl.Data) < 1 {
				sChan <- "NO_RESP"
				return
			}

			sChan <- sl.Data[0]
			time.Sleep(2 * time.Second)
		}
	}()

Loop:
	for {
		select {
		case s := <-sChan:
			switch s {
			case "RUNNING":
				t.Printf("测试中, 测试状态为: %s\n", s)
			case "SUCCESS":
				t.Printf("测试运行成功！测试状态为: %s\n", s)
				break Loop
			case "FAILURE":
				t.Printf("测试运行失败！测试状态为: %s\n", s)
				t.ErrChan <- fmt.Errorf("测试运行失败！测试状态为: %s", s)
				break Loop
			case "NO_RESP":
				t.Println("状态接口未返回任何data！测试失败")
				t.ErrChan <- errors.New("状态接口未返回任何data！测试失败")
				break Loop
			default:
				err := fmt.Errorf("未知的测试任务状态: %s, 测试失败", s)
				t.Println(err)
				t.ErrChan <- err
				break Loop
			}
		case <-time.After(60 * time.Second):
			err := errors.New("测试超过1分钟未结束，属于异常状态，测试失败！")
			t.Println(err)
			t.ErrChan <- err
			break Loop
		}
	}

	// 检查执行记录列表有一条记录
	var testRecords API.TestRecords
	t.Send(
		newapi.GetTestRecords{
			TestID: t.testID,
		}).DecodeJSON(&testRecords).AssertEqualInt("检查执行记录列表有一条记录", 1, len(testRecords.Data.Data))

	// 检查执行记录状态是成功
	t.Send(newapi.GetTestDetail{
		JobID: jobID,
	}).AssertContainString(`测试执行状态应为"成功"`, t.Resp, "SUCCESS")

	// 检查调度日志中显示测试执行成功
	t.Send(newapi.GetTestSystemLog{
		LastLogId: "0",
		TestID:    jobID,
	}).AssertContainString("调度日志应显示测试成功", t.Resp, "所有实例运行完成 测试任务执行成功 ")

	// 查询执行记录的pod信息
	var pods API.TestPods
	t.Send(
		newapi.GetTestPods{
			TestID:    t.testID,
			HistoryID: jobID,
		}).DecodeJSON(&pods).AssertEqualInt("查询执行记录的pod信息可以获取到一个pod", 1, len(pods.Data))
	podName := pods.Data[0].PodName
	podNamespace := pods.Data[0].Namespace

	// 下载Pod日志
	t.Send(
		newapi.DownloadTestLog{
			HistoryID:     jobID,
			EnvId:         t.envID,
			Namespace:     podNamespace,
			PodName:       podName,
			ContainerName: "app",
			SinceSeconds:  0,
			TailLine:      0,
			Timestamps:    false,
		}).AssertContainString("能成功下载pod日志且日志内容正确", t.Resp, "测试通过")

	// 执行测试后停止测试

	// 执行测试
	t.Send(newapi.ExecuteTest{
		ParameterList: plist,
		TestId:        t.testID,
	}).DecodeJSON(&gen).AssertContainString("执行测试应成功", t.Resp, "success")
	jobID = gen.Data

	time.Sleep(3 * time.Second)

	// TODO 停止测试
	t.Send(
		newapi.StopTest{
			TestID:    t.testID,
			HistoryId: jobID,
		}).AssertContainString("停止测试接口应调用成功", t.Resp, "success")

	t.Println("等待5秒")
	time.Sleep(5 * time.Second)

	// 检查调度日志中显示测试执行成功
	t.Send(newapi.GetTestSystemLog{
		LastLogId: "0",
		TestID:    jobID,
	}).AssertContainString("调度日志应显示测试任务已停止", t.Resp, fmt.Sprintf("停止测试任务(%v)成功", t.testID))

}

func (t *TestFeatureTest) Teardown() {
	if t.SkipTeardown == true {
		t.Println("Skipping Teardown...")
		t.logFile.Close()
		t.EndTest()
		return
	}

	defer t.logFile.Close()
	defer t.EndTest()

	t.Send(newapi.DeleteK8sEnv{
		EnvID: t.envID,
	}).AssertContainString("删除k8s环境应成功", t.Resp, "success")

	t.Send(newapi.DeleteTest{
		TestID: t.testID,
	}).AssertContainString("删除测试应成功", t.Resp, "success")

	t.Send(newapi.DeleteProject{
		ProjectID: t.projectID,
	}).AssertContainString("删除项目应成功", t.Resp, "success")

}
