package pipeline

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"gitlab.blackoutta.com/devops-acceptance-testing/v1/API"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/assertion"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/conf"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/errors"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/param"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/prep"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/random"
)

func RunPipelineTest(exitChan chan assertion.TestResult) {
	// 准备工作
	f, ast, sp, c := prep.SetupTest("流水线测试套件")
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
	sp.ImageTag = "pipeline-test"
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

	//echo 10.12.6.12 devops.build.com >> /etc/hosts
	hf, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		ast.Println(err)
	}

	if _, err := hf.WriteString(fmt.Sprintf("%v %v\n", conf.AccessNode, ingressHost)); err != nil {
		ast.Println(err)
	}

	if closeErr := hf.Close(); closeErr != nil {
		ast.Printf("error closing the file: %v\n", closeErr)
	}

	cmd2 := exec.Command("cat", "/etc/hosts")
	stdout, cmdErr := cmd2.Output()
	if cmdErr != nil {
		ast.Printf("error executing command: %v\n", cmdErr)
	}
	ast.Println("command result: " + string(stdout))

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

	time.Sleep(time.Second)

	// 复制流水线
	pipelineDuplicate := fmt.Sprintf("pip-%v", random.ShortGUID())
	resp = API.DuplicatePipeline(c, sp.PipelineID, pipelineDuplicate)
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("复制流水线", gen.ErrorInfo, resp)

	// 获取副本流水线ID
	var pips API.Pipelines
	resp = API.GetPipelieList(c)
	errors.UnmarshalAndHandleError(resp.Response, &pips)
	ast.AssertSuccess("获取副本流水线ID", pips.ErrorInfo, resp)
	sp.DuplicatePipID = pips.Data.Data[0].ID

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

}

func tearDown(c http.Client, ast *assertion.Assertion, sp *param.SuiteParams, exitChan chan assertion.TestResult) {
	ast.PrintTearDownStart()

	//删除流水线
	resp := API.DeletePipeline(c, sp.PipelineID)
	gen := API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除流水线", gen.ErrorInfo, resp)

	//删除流水线2
	resp = API.DeletePipeline(c, sp.DuplicatePipID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除流水线2", gen.ErrorInfo, resp)

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
