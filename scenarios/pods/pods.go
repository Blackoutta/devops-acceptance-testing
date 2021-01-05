package pods

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/API"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/req"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/assertion"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/conf"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/errors"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/param"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/prep"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/random"
	"golang.org/x/net/websocket"
)

func RunPodsTest(exitChan chan assertion.TestResult) {
	// 准备工作
	f, ast, sp, c := prep.SetupTest("应用测试套件")
	defer f.Close()

	// Tear Down动作，会删除所有测试资源
	defer tearDown(c, ast, sp, exitChan)

	// 如果有panic发生，记录panic错误日志，将测试结果置为失败，然后通过recover()让Tear Down可以正常执行
	defer ast.RecoverFromPanic()

	// 创建项目
	projectName := fmt.Sprintf("p-%v", random.ShortGUID())
	projectIdentifier := fmt.Sprintf("mp-%v", random.ShortGUID())
	resp := API.CreateProject(c, projectName, projectIdentifier)
	pjc := API.ProjectCreated{}
	errors.UnmarshalAndHandleError(resp.Response, &pjc)
	ast.AssertSuccess("创建项目", pjc.ErrorInfo, resp)

	// 关联项目ID
	resp = API.GetProjectDetail(c, projectName)
	pjs := API.Projects{}
	errors.UnmarshalAndHandleError(resp.Response, &pjs)
	sp.ProjectID = pjs.Data.Data[0].ID
	fmt.Printf("Project ID is: %v\n", sp.ProjectID)
	ast.AssertSuccess("获取项目ID", pjs.ErrorInfo, resp)

	// 创建GIT凭证
	gitCredentialName := fmt.Sprintf("git_%v", random.ShortGUID())
	resp = API.CreateGitCredential(c, sp.ProjectID, gitCredentialName)
	gen := API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建Git凭证", gen.ErrorInfo, resp)

	// 关联GIT凭证
	resp = API.GetGitCredential(c, sp.ProjectID, gitCredentialName)
	gcs := API.GitCredentials{}
	errors.UnmarshalAndHandleError(resp.Response, &gcs)
	sp.GitCredentialID = gcs.Data.Data[0].ID
	ast.AssertSuccess("获取Git Credential ID", gcs.ErrorInfo, resp)

	//创建Docker凭证
	dockerCredentialName := fmt.Sprintf("d_%v", random.ShortGUID())
	resp = API.CreateDockerCredential(c, sp.ProjectID, dockerCredentialName)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建Docker凭证", gen.ErrorInfo, resp)

	// 关联Docker凭证
	resp = API.GetDockerCredential(c, sp.ProjectID, dockerCredentialName)
	dcs := API.GitCredentials{}
	errors.UnmarshalAndHandleError(resp.Response, &dcs)
	sp.DockerCredentialID = dcs.Data.Data[0].ID
	ast.AssertSuccess("获取Docker Credential ID", dcs.ErrorInfo, resp)

	// 创建用户组
	groupName := "mytestgroup"
	resp = API.CreateUserGroup(c, sp.ProjectID, groupName)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建用户组", gen.ErrorInfo, resp)

	// 关联用户组ID
	resp = API.GetUserGroups(c, sp.ProjectID, groupName)
	ugp := API.UserGroups{}
	errors.UnmarshalAndHandleError(resp.Response, &ugp)
	sp.UserGroupID = ugp.Data.Data[0].ID
	sp.UserID = ugp.Data.Data[0].Creator
	ast.AssertSuccess("关联用户组ID", ugp.ErrorInfo, resp)

	// 把自己加入用户组
	// resp = API.AddUserToGroup(c, sp.ProjectID, sp.UserID, sp.UserGroupID)
	// gen = API.GeneralResp{}
	// errors.UnmarshalAndHandleError(resp.Response, &gen)
	// ast.AssertSuccess("将用户加入用户组", gen.ErrorInfo, resp)

	// 创建制品库
	artifactName := fmt.Sprintf("art_%v", random.ShortGUID())
	resp = API.CreateArtifactLibrary(c, artifactName, sp.ProjectID)
	alc := API.ArtifactLibraryCreated{}
	errors.UnmarshalAndHandleError(resp.Response, &alc)
	ast.AssertSuccess("创建制品库", alc.ErrorInfo, resp)

	// 创建构建
	buildName := fmt.Sprintf(`b_%v`, random.ShortGUID())
	resp = API.CreateBuild(c, buildName, sp.ProjectID)
	bc := API.BuildCreated{}
	err := json.Unmarshal(resp.Response, &bc)
	errors.HandleError("err unmarshaling buildCreated response", err)
	ast.AssertSuccess("创建构建", bc.ErrorInfo, resp)

	// 关联构建ID
	sp.BuildID = bc.Data

	// 关联制品库ID
	resp = API.GetArtifactID(c, artifactName, sp.ProjectID)
	ao := API.ArtifactObtained{}
	err = json.Unmarshal(resp.Response, &ao)
	errors.HandleError("err unmarshaling ArtifactObtained response", err)
	ast.AssertSuccess("关联制品库ID", ao.ErrorInfo, resp)

	sp.ArtifactID = ao.Data.Data[0].ID

	// 关联Git Source ID
	resp = API.GetBuildDetail(c, sp.BuildID)
	bdo := API.BuildDetailObtained{}
	err = json.Unmarshal(resp.Response, &bdo)
	errors.HandleError("err unmarshaling BuildDetailObtained response", err)
	ast.AssertSuccess("关联Git Source ID", bdo.ErrorInfo, resp)

	sp.SourceID = bdo.Data.Source.ID

	// 编辑构建
	sp.ImageTag = "pods-test"
	resp = API.EditBuild(c, sp.GitCredentialID, sp.DockerCredentialID, sp.ProjectID, sp.BuildID, buildName, sp.ArtifactID, sp.SourceID, sp.ImageTag)
	bed := API.BuildEdited{}
	err = json.Unmarshal(resp.Response, &bed)
	errors.HandleError("err unmarshaling BuildDetailObtained response", err)
	ast.AssertSuccess("向构建中添加步骤", bed.ErrorInfo, resp)

	//关联Commit ID
	resp = API.GetCommitID(c, sp.BuildID)
	bb := API.BeforeBuild{}
	err = json.Unmarshal(resp.Response, &bb)
	errors.HandleError("err unmarshaling BeforeBuild response", err)
	ast.AssertSuccess("关联Commit ID", bb.ErrorInfo, resp)

	sp.CommitID = bb.Data.SpecificBranchOrTag.CommitID

	// 创建一个K8S环境
	env := conf.ReadEnvFile()
	envName := "我的测试环境"
	namespace := fmt.Sprintf("n-%v", random.ShortGUID())
	resp = API.CreateK8sEnv(c, sp.ProjectID, envName, namespace, conf.ServerAddr, env, conf.EnvContext)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建K8S环境", gen.ErrorInfo, resp)

	// 获取环境ID
	resp = API.GetEnvDetail(c, sp.ProjectID)
	envs := API.Environments{}
	errors.UnmarshalAndHandleError(resp.Response, &envs)
	sp.EnvID = envs.Data.Data[0].ID
	ast.AssertSuccess("获取环境ID", envs.ErrorInfo, resp)

	// 更新环境，加入高级选项
	resp = API.UpdateEnv(c, sp.EnvID, namespace, conf.ServerAddr, env, conf.EnvContext, conf.ServerAddr)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("更新环境，加入高级选项", gen.ErrorInfo, resp)

	// 创建一个应用
	sp.AppName = fmt.Sprintf("app-%v", random.ShortGUID())
	resp = API.CreateApp(c, sp.ProjectID, sp.UserID, sp.AppName)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建应用", gen.ErrorInfo, resp)

	// 获取应用ID
	resp = API.GetAppDetail(c, sp.ProjectID, sp.AppName)
	apps := API.Apps{}
	errors.UnmarshalAndHandleError(resp.Response, &apps)
	ast.AssertSuccess("关联应用ID", apps.ErrorInfo, resp)
	sp.AppID = apps.Data.Data[0].ID

	// 创建一个部署
	deployName := fmt.Sprintf("dp-%v", random.ShortGUID())
	resp = API.CreateK8sImageDeployment(c, sp.AppID, sp.EnvID, sp.ProjectID, deployName, sp.DockerCredentialID, sp.ImageTag)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建部署", gen.ErrorInfo, resp)

	// 关联DeployID
	resp = API.GetDeployDetail(c, sp.ProjectID, deployName)
	dps := API.Deploys{}
	errors.UnmarshalAndHandleError(resp.Response, &dps)
	ast.AssertSuccess("查询并关联Deploy ID", dps.ErrorInfo, resp)
	sp.DeployID = dps.Data.Data[0].ID

	// 编辑部署
	ingressHost := fmt.Sprintf("devops.%v.com", random.ShortGUID())
	resp = API.EditDeployDetails(c, sp.DeployID, ingressHost, sp.ProjectID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("编辑部署，加入ingress服务暴露", gen.ErrorInfo, resp)

	// 创建流水线
	pipelineName := fmt.Sprintf("pip-%v", random.ShortGUID())
	resp = API.CreatePipeline(c, sp.ProjectID, pipelineName)
	pic := API.PipelineCreated{}
	errors.UnmarshalAndHandleError(resp.Response, &pic)
	ast.AssertSuccess("创建流水线", pic.ErrorInfo, resp)
	sp.PipelineID = pic.Data.ID

	// 获取流水线TriggerMode ID
	resp = API.GetPipelineDetail(c, sp.PipelineID)
	pd := API.PipelineDetail{}
	errors.UnmarshalAndHandleError(resp.Response, &pd)
	ast.AssertSuccess("获取流水线TriggerMode ID", pd.ErrorInfo, resp)
	sp.TriggerID = pd.Data.TriggerMode.ID

	// 编辑流水线
	resp = API.EditPipelineWithUnitTestAndCheckpoint(c, sp.PipelineID, sp.ProjectID, pipelineName, sp.BuildID, sp.DeployID, sp.TriggerID, sp.GitCredentialID, sp.UserID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("编辑流水线，加入构建和部署步骤", gen.ErrorInfo, resp)

	// 拿到unit test 步骤的ID
	resp = API.GetPipelineDetail(c, sp.PipelineID)
	pld2 := API.PipelineDetail{}
	errors.UnmarshalAndHandleError(resp.Response, &pld2)
	sp.UnitTestID = pld2.Data.PipelineSteps[0].ID

	// 执行流水线
	resp = API.RunPipeline(c, sp.PipelineID, sp.BuildID, sp.CommitID, sp.UnitTestID, "smoke_build")
	pr := API.PipelineRan{}
	errors.UnmarshalAndHandleError(resp.Response, &pr)
	sp.PipelineJobID = pr.PipelineJobID
	ast.AssertSuccess("执行流水线", pr.ErrorInfo, resp)

	ast.SleepWithCounter("等待人工卡点被触发", 10)

	// 拿到卡点步骤的ID
	var checkpointCount int

	for {
		if checkpointCount >= 30 {
			ast.FailTest("一直未获取到卡点操作步骤信息，测试将继续，但测试结果将为失败")
			break
		}
		resp = API.GetPipelineRuntime(c, sp.PipelineJobID)
		prt := API.PipelineRuntime{}
		errors.UnmarshalAndHandleError(resp.Response, &prt)
		stepNum := len(prt.Data.StepNodeList)
		checkpointCount++
		if stepNum < 2 {
			ast.SleepWithCounter("暂未获取到卡点步骤信息", 3)
			continue
		}
		sp.CheckpointID = prt.Data.StepNodeList[1].HistoryID
		ast.Println("成功获取到卡点步骤ID：", sp.CheckpointID)
		break
	}

	// 审核人工卡点，让其通过
	resp = API.DoCheckpoint(c, sp.CheckpointID, "SUCCESS")
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("审核人工卡点，让其通过", gen.ErrorInfo, resp)

	// 检查流水线执行成功
	ast.Println("检查流水线运行状态...")
	pipelineSuccess := false
	var counter int
	for {
		if pipelineSuccess == true {
			ast.Println("流水线运行成功!")
			break
		}

		if counter >= 24 {
			ast.FailTest("流水线一直处于运行中状态，属于异常行为，测试将继续但测试结果将为失败。")
			break
		}

		resp = API.GetPipelineHistory(c, sp.PipelineJobID)
		ph := API.PipelineHistory{}
		errors.UnmarshalAndHandleError(resp.Response, &ph)
		currentStatus := ph.Data.Status
		counter++

		switch currentStatus {
		case "WAITING":
			ast.Println("流水线状态为：等待中，等待3秒...")
			time.Sleep(3 * time.Second)
		case "RUNNING":
			ast.Println("流水线状态为：运行中，等待3秒...")
			time.Sleep(3 * time.Second)
			continue
		case "SUCCESS":
			time.Sleep(time.Second)
			pipelineSuccess = true
		case "FAILED":
			panic("流水线运行失败，请检查日志确认原因。")
		default:
			errMsg := fmt.Sprintf("流水线运行出现异常，状态为: %v\n", currentStatus)
			panic(errMsg)
		}
	}

	// 等待部署充分完成
	ast.Println("等待8秒，让部署充分完成...")
	time.Sleep(4 * time.Second)
	ast.Println("还有4秒...")
	time.Sleep(4 * time.Second)

	// 获取pod名
	resp = API.GetPods(c, sp.DeployID, sp.EnvID, sp.AppID)
	pl1 := API.PodList{}
	errors.UnmarshalAndHandleError(resp.Response, &pl1)
	ast.AssertIntegerEqual("获取pod名称", pl1.AppId, sp.AppID, resp)
	podNameOne := pl1.KubernetesDeploy[0].PodInstances[0].PodName

	// 验证进入容器功能
	tep := testEnterPod(sp.DeployID, sp.EnvID, namespace, podNameOne, "app", strconv.Itoa(sp.ProjectID), sp.UserID, ast)
	ast.AssertBooleanEqual("验证进入容器后可以用ls命令获取当前目录列表", tep, true, req.Record{})

	// 验证控制台日志功能
	tcl := testConsoleLog(sp.DeployID, sp.EnvID, namespace, podNameOne, "app", strconv.Itoa(sp.ProjectID), sp.UserID, 100, 0, ast)
	ast.AssertBooleanEqual("验证控制台日志能打印出应用的日志", tcl, true, req.Record{})

	// 验证下载控制台日志
	resp = API.DownloadPodLog(c, sp.EnvID, sp.DeployID, namespace, podNameOne, "app", conf.UserID)
	ddlog := string(resp.Response)
	ast.AssertContainString("验证下载的日志能正常被读取", ddlog, "server started", resp)

	// 重启pod
	resp = API.RebootPod(c, sp.EnvID, sp.DeployID, namespace, podNameOne)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("重启pod", gen.ErrorInfo, resp)

	ast.Println("pod重启中，等待15秒...")
	time.Sleep(15 * time.Second)

	// 验证重启pod后，新pod会替换老pod
	resp = API.GetPods(c, sp.DeployID, sp.EnvID, sp.AppID)
	pl2 := API.PodList{}
	errors.UnmarshalAndHandleError(resp.Response, &pl2)
	ast.AssertIntegerEqual("获取新pod的名称", pl2.AppId, sp.AppID, resp)
	podNameTwo := pl2.KubernetesDeploy[0].PodInstances[0].PodName
	ast.AssertStringNotEqual("验证重启pod后，新pod会替换老pod", podNameOne, podNameTwo, resp)

	// 验证可以使用扩缩容功能将应用数量扩展到3个
	resp = API.RescalePods(c, 3, sp.DeployID, sp.AppID, sp.EnvID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("执行扩容，将replicasets设为3", gen.ErrorInfo, resp)

	ast.Println("扩容中，等待12秒...")
	time.Sleep(12 * time.Second)

	resp = API.GetPods(c, sp.DeployID, sp.EnvID, sp.AppID)
	pl3 := API.PodList{}
	errors.UnmarshalAndHandleError(resp.Response, &pl3)
	numOfPods := len(pl3.KubernetesDeploy[0].PodInstances)
	ast.AssertIntegerEqual("pod数量已变化为3", numOfPods, 3, resp)

	// 验证变更资源
	resp = API.ChangePodResource(c, sp.DeployID, sp.AppID, sp.EnvID, "2000m", "1024Mi")
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("变更应用资源至cpu:2000m, mem:1024Mi", gen.ErrorInfo, resp)

	ast.Println("资源变更中，等待5秒...")
	time.Sleep(5 * time.Second)

	// 通过获取k8s描述文件，验证资源变更功能
	resp = API.GetK8sDescription(c, sp.DeployID, sp.AppID, sp.EnvID)
	kd := API.K8sDescription{}
	errors.UnmarshalAndHandleError(resp.Response, &kd)
	ast.AssertContainString("通过获取k8s描述文件，验证limits资源被正常修改", kd.Data, `limits: {cpu: '2', memory: 1Gi}`, resp)
	ast.AssertContainString("通过获取k8s描述文件，验证requests资源被正常修改", kd.Data, `requests: {cpu: 500m, memory: 1Gi}`, resp)
}

func RunPrometheusTest(exitChan chan assertion.TestResult) {
	// 准备工作
	f, ast, sp, c := prep.SetupTest("Prometheus测试套件")
	defer f.Close()

	// Tear Down动作，会删除所有测试资源
	defer tearDown(c, ast, sp, exitChan)

	// 如果有panic发生，记录panic错误日志，将测试结果置为失败，然后通过recover()让Tear Down可以正常执行
	defer ast.RecoverFromPanic()

	// 创建项目
	projectName := fmt.Sprintf("p-%v", random.ShortGUID())
	projectIdentifier := fmt.Sprintf("mp-%v", random.ShortGUID())
	resp := API.CreateProject(c, projectName, projectIdentifier)
	pjc := API.ProjectCreated{}
	errors.UnmarshalAndHandleError(resp.Response, &pjc)
	ast.AssertSuccess("创建项目", pjc.ErrorInfo, resp)

	// 关联项目ID
	resp = API.GetProjectDetail(c, projectName)
	pjs := API.Projects{}
	errors.UnmarshalAndHandleError(resp.Response, &pjs)
	sp.ProjectID = pjs.Data.Data[0].ID
	ast.AssertSuccess("获取项目ID", pjs.ErrorInfo, resp)

	// 创建GIT凭证
	gitCredentialName := fmt.Sprintf("git_%v", random.ShortGUID())
	resp = API.CreateGitCredential(c, sp.ProjectID, gitCredentialName)
	gen := API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建Git凭证", gen.ErrorInfo, resp)

	// 关联GIT凭证
	resp = API.GetGitCredential(c, sp.ProjectID, gitCredentialName)
	gcs := API.GitCredentials{}
	errors.UnmarshalAndHandleError(resp.Response, &gcs)
	sp.GitCredentialID = gcs.Data.Data[0].ID
	ast.AssertSuccess("获取Git Credential ID", gcs.ErrorInfo, resp)

	//创建Docker凭证
	dockerCredentialName := fmt.Sprintf("d_%v", random.ShortGUID())
	resp = API.CreateDockerCredential(c, sp.ProjectID, dockerCredentialName)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建Docker凭证", gen.ErrorInfo, resp)

	// 关联Docker凭证
	resp = API.GetDockerCredential(c, sp.ProjectID, dockerCredentialName)
	dcs := API.GitCredentials{}
	errors.UnmarshalAndHandleError(resp.Response, &dcs)
	sp.DockerCredentialID = dcs.Data.Data[0].ID
	ast.AssertSuccess("获取Docker Credential ID", dcs.ErrorInfo, resp)

	// 创建用户组
	groupName := "mytestgroup"
	resp = API.CreateUserGroup(c, sp.ProjectID, groupName)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建用户组", gen.ErrorInfo, resp)

	// 关联用户组ID
	resp = API.GetUserGroups(c, sp.ProjectID, groupName)
	ugp := API.UserGroups{}
	errors.UnmarshalAndHandleError(resp.Response, &ugp)
	sp.UserGroupID = ugp.Data.Data[0].ID
	sp.UserID = ugp.Data.Data[0].Creator
	ast.AssertSuccess("关联用户组ID", ugp.ErrorInfo, resp)

	// 把自己加入用户组
	// resp = API.AddUserToGroup(c, sp.ProjectID, sp.UserID, sp.UserGroupID)
	// gen = API.GeneralResp{}
	// errors.UnmarshalAndHandleError(resp.Response, &gen)
	// ast.AssertSuccess("将用户加入用户组", gen.ErrorInfo, resp)

	// 创建制品库
	artifactName := fmt.Sprintf("art_%v", random.ShortGUID())
	resp = API.CreateArtifactLibrary(c, artifactName, sp.ProjectID)
	alc := API.ArtifactLibraryCreated{}
	errors.UnmarshalAndHandleError(resp.Response, &alc)
	ast.AssertSuccess("创建制品库", alc.ErrorInfo, resp)

	// 创建构建
	buildName := fmt.Sprintf(`b_%v`, random.ShortGUID())
	resp = API.CreateBuild(c, buildName, sp.ProjectID)
	bc := API.BuildCreated{}
	err := json.Unmarshal(resp.Response, &bc)
	errors.HandleError("err unmarshaling buildCreated response", err)
	ast.AssertSuccess("创建构建", bc.ErrorInfo, resp)

	// 关联构建ID
	sp.BuildID = bc.Data

	// 关联制品库ID
	resp = API.GetArtifactID(c, artifactName, sp.ProjectID)
	ao := API.ArtifactObtained{}
	err = json.Unmarshal(resp.Response, &ao)
	errors.HandleError("err unmarshaling ArtifactObtained response", err)
	ast.AssertSuccess("关联制品库ID", ao.ErrorInfo, resp)

	sp.ArtifactID = ao.Data.Data[0].ID

	// 关联Git Source ID
	resp = API.GetBuildDetail(c, sp.BuildID)
	bdo := API.BuildDetailObtained{}
	err = json.Unmarshal(resp.Response, &bdo)
	errors.HandleError("err unmarshaling BuildDetailObtained response", err)
	ast.AssertSuccess("关联Git Source ID", bdo.ErrorInfo, resp)

	sp.SourceID = bdo.Data.Source.ID

	// 编辑构建
	sp.ImageTag = "prometheus-test"
	resp = API.EditBuild(c, sp.GitCredentialID, sp.DockerCredentialID, sp.ProjectID, sp.BuildID, buildName, sp.ArtifactID, sp.SourceID, sp.ImageTag)
	bed := API.BuildEdited{}
	err = json.Unmarshal(resp.Response, &bed)
	errors.HandleError("err unmarshaling BuildDetailObtained response", err)
	ast.AssertSuccess("向构建中添加步骤", bed.ErrorInfo, resp)

	//关联Commit ID
	resp = API.GetCommitID(c, sp.BuildID)
	bb := API.BeforeBuild{}
	err = json.Unmarshal(resp.Response, &bb)
	errors.HandleError("err unmarshaling BeforeBuild response", err)
	ast.AssertSuccess("关联Commit ID", bb.ErrorInfo, resp)

	sp.CommitID = bb.Data.SpecificBranchOrTag.CommitID

	// 创建一个K8S环境
	env := conf.ReadEnvFile()
	envName := "我的测试环境"
	namespace := fmt.Sprintf("n-%v", random.ShortGUID())
	resp = API.CreateK8sEnv(c, sp.ProjectID, envName, namespace, conf.ServerAddr, env, conf.EnvContext)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建K8S环境", gen.ErrorInfo, resp)

	// 获取环境ID
	resp = API.GetEnvDetail(c, sp.ProjectID)
	envs := API.Environments{}
	errors.UnmarshalAndHandleError(resp.Response, &envs)
	sp.EnvID = envs.Data.Data[0].ID
	ast.AssertSuccess("获取环境ID", envs.ErrorInfo, resp)

	// 更新环境，加入高级选项
	resp = API.UpdateEnv(c, sp.EnvID, namespace, conf.ServerAddr, env, conf.EnvContext, conf.ServerAddr)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("更新环境，加入高级选项", gen.ErrorInfo, resp)

	// 创建一个应用
	sp.AppName = fmt.Sprintf("app-%v", random.ShortGUID())
	resp = API.CreateApp(c, sp.ProjectID, sp.UserID, sp.AppName)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建应用", gen.ErrorInfo, resp)

	// 获取应用ID
	resp = API.GetAppDetail(c, sp.ProjectID, sp.AppName)
	apps := API.Apps{}
	errors.UnmarshalAndHandleError(resp.Response, &apps)
	ast.AssertSuccess("关联应用ID", apps.ErrorInfo, resp)
	sp.AppID = apps.Data.Data[0].ID

	// 创建一个部署
	deployName := fmt.Sprintf("dp-%v", random.ShortGUID())
	resp = API.CreateK8sImageDeployment(c, sp.AppID, sp.EnvID, sp.ProjectID, deployName, sp.DockerCredentialID, sp.ImageTag)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建部署", gen.ErrorInfo, resp)

	// 关联DeployID
	resp = API.GetDeployDetail(c, sp.ProjectID, deployName)
	dps := API.Deploys{}
	errors.UnmarshalAndHandleError(resp.Response, &dps)
	ast.AssertSuccess("查询并关联Deploy ID", dps.ErrorInfo, resp)
	sp.DeployID = dps.Data.Data[0].ID

	// 编辑部署
	ingressHost := fmt.Sprintf("devops.%v.com", random.ShortGUID())
	resp = API.EditDeployDetails(c, sp.DeployID, ingressHost, sp.ProjectID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("编辑部署，加入ingress服务暴露", gen.ErrorInfo, resp)

	// 创建流水线
	pipelineName := fmt.Sprintf("pip-%v", random.ShortGUID())
	resp = API.CreatePipeline(c, sp.ProjectID, pipelineName)
	pic := API.PipelineCreated{}
	errors.UnmarshalAndHandleError(resp.Response, &pic)
	ast.AssertSuccess("创建流水线", pic.ErrorInfo, resp)
	sp.PipelineID = pic.Data.ID

	// 获取流水线TriggerMode ID
	resp = API.GetPipelineDetail(c, sp.PipelineID)
	pd := API.PipelineDetail{}
	errors.UnmarshalAndHandleError(resp.Response, &pd)
	ast.AssertSuccess("获取流水线TriggerMode ID", pd.ErrorInfo, resp)
	sp.TriggerID = pd.Data.TriggerMode.ID

	// 编辑流水线
	resp = API.EditPipelineWithUnitTestAndCheckpoint(c, sp.PipelineID, sp.ProjectID, pipelineName, sp.BuildID, sp.DeployID, sp.TriggerID, sp.GitCredentialID, sp.UserID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("编辑流水线，加入构建和部署步骤", gen.ErrorInfo, resp)

	// 拿到unit test 步骤的ID
	resp = API.GetPipelineDetail(c, sp.PipelineID)
	pld2 := API.PipelineDetail{}
	errors.UnmarshalAndHandleError(resp.Response, &pld2)
	sp.UnitTestID = pld2.Data.PipelineSteps[0].ID

	// 执行流水线
	resp = API.RunPipeline(c, sp.PipelineID, sp.BuildID, sp.CommitID, sp.UnitTestID, "smoke_build")
	pr := API.PipelineRan{}
	errors.UnmarshalAndHandleError(resp.Response, &pr)
	sp.PipelineJobID = pr.PipelineJobID
	ast.AssertSuccess("执行流水线", pr.ErrorInfo, resp)

	ast.SleepWithCounter("等待人工卡点被触发", 10)

	// 拿到卡点步骤的ID
	var checkpointCount int

	for {
		if checkpointCount >= 30 {
			ast.FailTest("一直未获取到卡点操作步骤信息，测试将继续，但测试结果将为失败")
			break
		}
		resp = API.GetPipelineRuntime(c, sp.PipelineJobID)
		prt := API.PipelineRuntime{}
		errors.UnmarshalAndHandleError(resp.Response, &prt)
		stepNum := len(prt.Data.StepNodeList)
		checkpointCount++
		if stepNum < 2 {
			ast.SleepWithCounter("暂未获取到卡点步骤信息", 3)
			continue
		}
		sp.CheckpointID = prt.Data.StepNodeList[1].HistoryID
		ast.Println("成功获取到卡点步骤ID：", sp.CheckpointID)
		break
	}

	// 审核人工卡点，让其通过
	resp = API.DoCheckpoint(c, sp.CheckpointID, "SUCCESS")
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("审核人工卡点，让其通过", gen.ErrorInfo, resp)

	// 检查流水线执行成功
	ast.Println("检查流水线运行状态...")
	pipelineSuccess := false
	var counter int
	for {
		if pipelineSuccess == true {
			ast.Println("流水线运行成功!")
			break
		}

		if counter >= 24 {
			ast.FailTest("流水线一直处于运行中状态，属于异常行为，测试将继续但测试结果将为失败。")
			break
		}

		resp = API.GetPipelineHistory(c, sp.PipelineJobID)
		ph := API.PipelineHistory{}
		errors.UnmarshalAndHandleError(resp.Response, &ph)
		currentStatus := ph.Data.Status
		counter++

		switch currentStatus {
		case "WAITING":
			ast.Println("流水线状态为：等待中，等待3秒...")
			time.Sleep(3 * time.Second)
		case "RUNNING":
			ast.Println("流水线状态为：运行中，等待3秒...")
			time.Sleep(3 * time.Second)
			continue
		case "SUCCESS":
			time.Sleep(time.Second)
			pipelineSuccess = true
		case "FAILED":
			panic("流水线运行失败，请检查日志确认原因。")
		default:
			errMsg := fmt.Sprintf("流水线运行出现异常，状态为: %v\n", currentStatus)
			panic(errMsg)
		}
	}

	// 等待部署充分完成
	ast.Println("等待8秒，让部署充分完成...")
	time.Sleep(4 * time.Second)
	ast.Println("还有4秒...")
	time.Sleep(4 * time.Second)

	// 获取pod名
	resp = API.GetPods(c, sp.DeployID, sp.EnvID, sp.AppID)
	pl1 := API.PodList{}
	errors.UnmarshalAndHandleError(resp.Response, &pl1)
	ast.AssertIntegerEqual("获取pod名称", pl1.AppId, sp.AppID, resp)
	podName := pl1.KubernetesDeploy[0].PodInstances[0].PodName
	nodeName := pl1.KubernetesDeploy[0].PodInstances[0].NodeName

	ast.SleepWithCounter("检查promethes数据点应可以被获取，请耐心等待...", 240)

	d, err := time.ParseDuration("30m")
	if err != nil {
		log.Fatalln(err)
	}
	end := time.Now().Unix()
	r := d.Seconds()
	start := end - int64(r)
	step := int64(30)

	resp = API.GetPrometheusData(c, sp.EnvID, nodeName, podName, start, end, step)
	prom := API.PrometheusData{}
	errors.UnmarshalAndHandleError(resp.Response, &prom)

	cpuData := prom.Data.Cpu.Result[0].Values
	memData := prom.Data.Mem.Result[0].Values
	networkReceiveData := prom.Data.NetworkReceive.Result[0].Values
	networkTransmitData := prom.Data.NetworkTransmit.Result[0].Values
	ast.AssertIntegerGreaterThan("应能获取到cpu数据点", len(cpuData), 0, resp)
	ast.AssertIntegerGreaterThan("应能获取到mem数据点", len(memData), 0, resp)
	ast.AssertIntegerGreaterThan("应能获取到network receive数据点", len(networkReceiveData), 0, resp)
	ast.AssertIntegerGreaterThan("应能获取到network transmit数据点", len(networkTransmitData), 0, resp)

}

func tearDown(c http.Client, ast *assertion.Assertion, sp *param.SuiteParams, exitChan chan assertion.TestResult) {
	ast.PrintTearDownStart()

	//删除流水线
	resp := API.DeletePipeline(c, sp.PipelineID)
	gen := API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除流水线", gen.ErrorInfo, resp)

	//删除构建
	resp = API.DeleteBuild(c, sp.BuildID)
	d2 := API.ItemDeleted{}
	err := json.Unmarshal(resp.Response, &d2)
	errors.HandleError("err unmarshaling ItemDeleted 2 response", err)
	ast.AssertSuccess("删除构建", d2.ErrorInfo, resp)

	//删除制品库
	resp = API.DeleteArtifact(c, sp.ArtifactID)
	d3 := API.ItemDeleted{}
	err = json.Unmarshal(resp.Response, &d3)
	errors.HandleError("err unmarshaling ItemDeleted 1 response", err)
	ast.AssertSuccess("删除制品库", d3.ErrorInfo, resp)

	// 删除部署
	resp = API.DeleteDeploy(c, sp.DeployID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除部署", gen.ErrorInfo, resp)

	// 删除应用
	resp = API.DeleteApp(c, sp.AppID, sp.AppName)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除应用", gen.ErrorInfo, resp)

	// 删除环境
	resp = API.DeleteEnv(c, sp.EnvID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除环境", gen.ErrorInfo, resp)

	// 删除项目
	resp = API.DeleteProject(c, sp.ProjectID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除项目", gen.ErrorInfo, resp)

	ast.PrintTearDownEnd()
	// 判断测试成功与否
	ast.CheckSuiteResult(exitChan)
}

func testEnterPod(deployID, envID int, namespace string, podName string, containerName string, projectID string, userID string, ast *assertion.Assertion) bool {
	wsurl := fmt.Sprintf(conf.WsHost+"ws/k8s-deploy/ws/exec?deployId=%v&env_id=%v&namespace=%v&pod_name=%v&container_name=%v&projectId=%v&userId=%v",
		deployID, envID, namespace, podName, containerName, projectID, userID)
	origin := retrieveOrigin(conf.Host)

	ast.Printf("连接websocket: %v\n", wsurl)
	ws, err := websocket.Dial(wsurl, "", origin)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	textChan := make(chan string)

	go func(chan string) {
		for {
			data := make([]byte, 1024)
			n, err := ws.Read(data)
			if err == io.EOF {
				ast.Println("ws: encocuntered EOF error while reading ws data, exiting go routine...")
				return
			}
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					ast.Println("ws connection was closed, exiting this goroutine...")
					return
				}
				ast.Printf("error while reading ws data: %v", err)
				return
			}
			result := string(data[:n])
			textChan <- result
		}
	}(textChan)

	time.Sleep(time.Second)

	go func() {
		_, writeErr := ws.Write([]byte("l"))
		if writeErr != nil {
			ast.Printf("error while writing ws data to server: %v", writeErr)
		}

		_, writeErr = ws.Write([]byte("s"))
		if writeErr != nil {
			ast.Printf("error while writing ws data to server: %v", writeErr)
		}

		_, writeErr = ws.Write([]byte("\r"))
		if writeErr != nil {
			ast.Printf("error while writing ws data to server: %v", writeErr)
		}
	}()

	for {
		select {
		case l := <-textChan:
			if strings.Contains(l, "httpserver") {
				return true
			}
		case <-time.After(5 * time.Second):
			return false
		}
	}
}

func testConsoleLog(deployID, envID int, namespace string, podName string, containerName string, projectID string, userID string, tail int, since int, ast *assertion.Assertion) bool {
	wsurl := fmt.Sprintf(conf.WsHost+"ws/k8s-deploy/ws/log?deployId=%v&env_id=%v&namespace=%v&pod_name=%v&container_name=%v&projectId=%v&userId=%v&tail_line=%v&since_time=%v",
		deployID, envID, namespace, podName, containerName, projectID, userID, tail, since)
	origin := retrieveOrigin(conf.Host)

	ast.Printf("连接websocket: %v\n", wsurl)

	ws, err := websocket.Dial(wsurl, "", origin)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	textChan := make(chan string)

	go func(chan string) {
		for {
			data := make([]byte, 1024)
			n, err := ws.Read(data)

			if err == io.EOF {
				ast.Println("encocuntered EOF error while reading ws data, exiting go routine")
				return
			}
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					ast.Println("ws connection was closed, exiting this goroutine...")
					return
				}
				ast.Printf("error reading ws data: %v", err)
				return
			}
			result := string(data[:n])
			textChan <- result
		}
	}(textChan)

	for {
		select {
		case l := <-textChan:
			if strings.Contains(l, "server started") {
				return true
			}
		case <-time.After(5 * time.Second):
			return false
		}
	}
}

func retrieveOrigin(s string) string {
	result := strings.SplitAfterN(s, ".com", 2)
	return result[0]
}
