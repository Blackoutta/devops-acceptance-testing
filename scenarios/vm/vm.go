package vm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/API"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/req"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/assertion"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/errors"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/param"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/prep"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/random"
)

func RunVmTest(exitChan chan assertion.TestResult) {
	// 准备工作
	f, ast, sp, c := prep.SetupTest("虚拟机部署测试套件")
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

	// 创建用户组
	groupName := "mytestgroup"
	resp = API.CreateUserGroup(c, sp.ProjectID, groupName)
	gen := API.GeneralResp{}
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

	// 创建GIT凭证
	gitCredentialName := fmt.Sprintf("git_%v", random.ShortGUID())
	resp = API.CreateGitCredential(c, sp.ProjectID, gitCredentialName)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("创建Git凭证", gen.ErrorInfo, resp)

	// 关联GIT凭证
	resp = API.GetGitCredential(c, sp.ProjectID, gitCredentialName)
	gcs := API.GitCredentials{}
	errors.UnmarshalAndHandleError(resp.Response, &gcs)
	sp.GitCredentialID = gcs.Data.Data[0].ID
	ast.AssertSuccess("获取Git Credential ID", gcs.ErrorInfo, resp)

	//创建Docker凭证
	dockerCredentialName := fmt.Sprintf("dk_%v", random.ShortGUID())
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

	// TODO 创建制品库
	// 创建制品库
	artifactName := fmt.Sprintf("art_%v", random.ShortGUID())
	resp = API.CreateArtifactLibrary(c, artifactName, sp.ProjectID)
	alc := API.ArtifactLibraryCreated{}
	errors.UnmarshalAndHandleError(resp.Response, &alc)
	ast.AssertSuccess("创建制品库", alc.ErrorInfo, resp)

	// 关联制品库ID
	resp = API.GetArtifactID(c, artifactName, sp.ProjectID)
	ao := API.ArtifactObtained{}
	err := json.Unmarshal(resp.Response, &ao)
	errors.HandleError("err unmarshaling ArtifactObtained response", err)
	ast.AssertSuccess("关联制品库ID", ao.ErrorInfo, resp)
	sp.ArtifactID = ao.Data.Data[0].ID

	// TODO 创建构建
	buildName := fmt.Sprintf(`b_%v`, random.ShortGUID())
	resp = API.CreateVMBuild(c, buildName, sp.ProjectID)
	bc := API.BuildCreated{}
	err = json.Unmarshal(resp.Response, &bc)
	errors.HandleError("err unmarshaling buildCreated response", err)
	ast.AssertSuccess("创建构建", bc.ErrorInfo, resp)

	// 关联构建ID
	sp.BuildID = bc.Data

	// 关联Git Source ID
	resp = API.GetBuildDetail(c, sp.BuildID)
	bdo := API.BuildDetailObtained{}
	err = json.Unmarshal(resp.Response, &bdo)
	errors.HandleError("err unmarshaling BuildDetailObtained response", err)
	ast.AssertSuccess("关联Git Source ID", bdo.ErrorInfo, resp)

	sp.SourceID = bdo.Data.Source.ID

	// TODO 编辑构建
	// 编辑构建
	resp = API.EditVMBuild(c, sp.ProjectID, sp.BuildID, buildName, sp.SourceID, sp.ArtifactID)
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
	// 创建主机组
	vmGroupName := fmt.Sprintf("vmg-%v", random.ShortGUID())
	gen, resp = API.CreateVMGroup(c, sp.ProjectID, vmGroupName)
	ast.AssertSuccess("创建主机组", gen.ErrorInfo, resp)

	// 关联主机组ID
	sp.VMGroupID = gen.Data

	// TODO 创建两个主机
	ip1 := "10.12.6.12"
	ip2 := "10.12.6.13"
	gen, resp = API.CreateVM(c, ip1, 22, "app", "rX8O60lGgSQK!%ey", sp.VMGroupID, "PASSWORD")
	ast.AssertSuccess("创建虚拟机"+ip1, gen.ErrorInfo, resp)

	// 关联主机1的ID
	//vmOneID := gen.Data

	gen, resp = API.CreateVM(c, ip2, 22, "app", "rX8O60lGgSQK!%ey", sp.VMGroupID, "PASSWORD")
	ast.AssertSuccess("创建虚拟机"+ip2, gen.ErrorInfo, resp)

	// 关联主机2的ID
	//vmTwoID := gen.Data

	// TODO 检查主机连通性
	gen, resp = API.CheckVMStatus(c, sp.VMGroupID)
	ast.AssertSuccess("检查主机连通性", gen.ErrorInfo, resp)

	// TODO 创建应用
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

	// TODO 创建部署
	vmDeployName := fmt.Sprintf("vmd-%v", random.ShortGUID())
	gen, resp = API.CreateVMDeploy(c, sp.AppID, vmDeployName, sp.ProjectID, sp.ArtifactID, sp.VMGroupID)
	ast.AssertSuccess("创建虚拟机部署", gen.ErrorInfo, resp)
	// 关联部署ID
	sp.DeployID = gen.Data

	// TODO 编辑部署
	gen, resp = API.EditVMDeploy(c, sp.DeployID)
	ast.AssertSuccess("编辑虚拟机部署，加入部署步骤", gen.ErrorInfo, resp)

	// TODO 创建流水线
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
	gen, resp = API.EditPipeline(c, sp.PipelineID, sp.ProjectID, pipelineName, sp.BuildID, sp.DeployID, sp.TriggerID)
	ast.AssertSuccess("编辑流水线，加入构建和部署步骤", gen.ErrorInfo, resp)

	// TODO 执行流水线
	resp = API.RunPipeline(c, sp.PipelineID, sp.BuildID, sp.CommitID, sp.UnitTestID, "smoke_build")
	pr := API.PipelineRan{}
	errors.UnmarshalAndHandleError(resp.Response, &pr)
	sp.PipelineJobID = pr.PipelineJobID
	ast.AssertSuccess("执行流水线", pr.ErrorInfo, resp)

	// 检查流水线执行成功
	ast.Println("检查流水线运行状态...")
	pipelineSuccess := false
	var counter int
	for {
		if pipelineSuccess == true {
			ast.Println("流水线运行成功!")
			break
		}

		if counter >= 35 {
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
			panic("流水线运行失败，请检查日志定位问题")
		default:
			errMsg := fmt.Sprintf("流水线运行出现异常，状态为: %v\n", currentStatus)
			panic(errMsg)
		}
	}

	// TODO 访问应用
	res, err := visitApp(ip1)
	if err != nil {
		panic(err)
	}
	ast.AssertContainString("访问虚拟机1的应用首页，应可以看到首页有正确响应", res, "成功部署", req.Record{})

	res, err = visitApp(ip2)
	if err != nil {
		panic(err)
	}
	ast.AssertContainString("访问虚拟机2的应用首页，应可以看到首页有正确响应", res, "成功部署", req.Record{})

	// TODO 停止并清理虚拟机中的应用
	gen, resp = API.EditVMDeployToTeardown(c, sp.DeployID)
	ast.AssertSuccess("将虚拟机部署编辑为清理部署", gen.ErrorInfo, resp)

	// 执行清理部署
	resp = API.ExecuteDeploy(c, sp.DeployID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("执行部署", gen.ErrorInfo, resp)
	sp.DeployJobID = gen.Data

	// 用JOB ID 查询部署记录，断言其status为SUCCESS
	deploySuccess := false
	for {
		if deploySuccess == true {
			break
		}
		resp = API.GetDeployHistory(c, sp.DeployJobID)
		dh := API.DeployHistory{}
		err := json.Unmarshal(resp.Response, &dh)
		errors.HandleError("err unmarshaling BuildRecord Obtained response", err)
		currentStatus := dh.Data.Status

		switch currentStatus {
		case "WAITING":
			ast.Println("部署状态为：等待中，等待5秒...")
			time.Sleep(5 * time.Second)
		case "RUNNING":
			ast.Println("部署状态为：运行中，等待5秒...")
			time.Sleep(5 * time.Second)
			continue
		case "SUCCESS":
			ast.Println("部署成功!")
			time.Sleep(time.Second)
			deploySuccess = true
		case "FAILURE":
			ast.Println("部署失败，请检查日志确认原因。")
		default:
			ast.Printf("部署出现异常，状态为: %v\n", currentStatus)
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

	//删除构建
	resp = API.DeleteBuild(c, sp.BuildID)
	d2 := API.ItemDeleted{}
	err := json.Unmarshal(resp.Response, &d2)
	errors.HandleError("err unmarshaling ItemDeleted 2 response", err)
	ast.AssertSuccess("删除构建", d2.ErrorInfo, resp)

	// 删除部署
	resp = API.DeleteDeploy(c, sp.DeployID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除部署", gen.ErrorInfo, resp)

	//删除制品库
	resp = API.DeleteArtifact(c, sp.ArtifactID)
	d3 := API.ItemDeleted{}
	err = json.Unmarshal(resp.Response, &d3)
	errors.HandleError("err unmarshaling ItemDeleted 1 response", err)
	ast.AssertSuccess("删除制品库", d3.ErrorInfo, resp)

	// 删除应用
	resp = API.DeleteApp(c, sp.AppID, sp.AppName)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除应用", gen.ErrorInfo, resp)

	//TODO 删除主机组
	gen, resp = API.DeleteVMGroup(c, sp.VMGroupID)
	ast.AssertSuccess("删除主机组", gen.ErrorInfo, resp)

	// 删除项目
	resp = API.DeleteProject(c, sp.ProjectID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除项目", gen.ErrorInfo, resp)

	ast.PrintTearDownEnd()
	// 判断测试成功与否
	ast.CheckSuiteResult(exitChan)
}

func visitApp(IP string) (string, error) {
	res, err := http.Get("http://" + IP + ":9005/")
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	return string(body), nil
}
