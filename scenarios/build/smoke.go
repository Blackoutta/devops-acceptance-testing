package build

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/conf"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/param"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/prep"

	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/API"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/req"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/assertion"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/errors"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/random"
)

func RunBuildAndDeployTest(exitChan chan assertion.TestResult) {
	// 准备工作
	f, ast, sp, c := prep.SetupTest("构建+部署测试套件")
	defer f.Close()

	// Tear Down动作，会删除所有测试资源
	defer tearDown(c, ast, sp, exitChan)

	// 如果有panic发生，记录panic错误日志，将测试结果置为失败，然后通过recover()让Tear Down可以正常执行
	defer ast.RecoverFromPanic()

	// 创建项目
	projectName := fmt.Sprintf("p_%v", random.ShortGUID())
	projectIdentifier := fmt.Sprintf("pi-%v", random.ShortGUID())
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
	artifactName := fmt.Sprintf("art-%v", random.ShortGUID())
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
	sp.ImageTag = "build-test"
	resp = API.EditBuild(c, sp.GitCredentialID, sp.DockerCredentialID, sp.ProjectID, sp.BuildID, buildName, sp.ArtifactID, sp.SourceID, sp.ImageTag)
	bed := API.BuildEdited{}
	err = json.Unmarshal(resp.Response, &bed)
	errors.HandleError("err unmarshaling BuildDetailObtained response", err)
	ast.AssertSuccess("向构建中添加步骤", bed.ErrorInfo, resp)

	time.Sleep(time.Second)

	// 复制构建
	buildDuplicate := fmt.Sprintf(`b_%v`, random.ShortGUID())
	resp = API.DuplicateBuild(c, buildDuplicate, sp.BuildID)
	err = json.Unmarshal(resp.Response, &gen)
	errors.HandleError("err unmarshaling buildCreated response", err)
	ast.AssertSuccess("复制构建", gen.ErrorInfo, resp)

	time.Sleep(2 * time.Second)

	// 查询构建列表
	var builds API.Builds
	resp = API.GetBuildList(c, sp.ProjectID)
	err = json.Unmarshal(resp.Response, &builds)
	errors.HandleError("err unmarshaling buildCreated response", err)
	ast.AssertSuccess("查询构建列表", builds.ErrorInfo, resp)
	fmt.Println(string(resp.Response))

	sp.DuplicateBuildID = builds.Data.Data[0].ID

	//关联Commit ID
	resp = API.GetCommitID(c, sp.BuildID)
	bb := API.BeforeBuild{}
	err = json.Unmarshal(resp.Response, &bb)
	errors.HandleError("err unmarshaling BeforeBuild response", err)
	ast.AssertSuccess("关联Commit ID", bb.ErrorInfo, resp)

	sp.CommitID = bb.Data.SpecificBranchOrTag.CommitID

	//执行构建
	resp = API.RunBuild(c, sp.BuildID, sp.CommitID, "smoke_build")
	br := API.BuildRan{}
	err = json.Unmarshal(resp.Response, &br)
	errors.HandleError("err unmarshaling BuildRan response", err)
	sp.BuildJobID = br.Data
	ast.AssertSuccess("执行构建", br.ErrorInfo, resp)

	// 用JOB ID 查询构建记录，断言其status字段为SUCCESS
	buildSuccess := false
	for {
		if buildSuccess == true {
			break
		}
		resp = API.GetBuildRecord(c, sp.BuildJobID)
		bro := API.BuildRecordObtained{}
		err = json.Unmarshal(resp.Response, &bro)
		errors.HandleError("err unmarshaling BuildRecord Obtained response", err)
		currentStatus := bro.Data.Status

		switch currentStatus {
		case "WAITING":
			ast.Println("构建状态为：等待中，等待3秒...")
			time.Sleep(3 * time.Second)
		case "BUILDING":
			ast.Println("构建状态为：构建中，等待3秒...")
			time.Sleep(3 * time.Second)
			continue
		case "SUCCESS":
			ast.Println("构建成功!")
			time.Sleep(time.Second)
			buildSuccess = true
		case "FAILURE":
			panic("构建失败，请检查日志确认原因。")
		default:
			err := fmt.Sprintf("构建异常，状态为： %v\n", currentStatus)
			panic(err)
		}
	}

	// 检查发布记录
	resp = API.GetArtifactUploadRecord(c, sp.ArtifactID)
	ulr := API.UploadRecords{}
	errors.UnmarshalAndHandleError(resp.Response, &ulr)
	ast.AssertIntegerEqual("检查发布记录的发布方式是：构建任务", ulr.Data.Data[0].UploadType, 2, resp)

	// 检查下载记录
	resp = API.GetArtifactDownloadRecord(c, sp.ArtifactID)
	dlr := API.DownloadRecords{}
	err = json.Unmarshal(resp.Response, &dlr)
	errors.HandleError("err unmarshaling DownloadRecords response", err)
	ast.AssertIntegerEqual("检查应有一条下载记录", len(dlr.Data.Data), 1, resp)

	// 创建一个K8S环境
	env := conf.ReadEnvFile()
	envName := fmt.Sprintf("env-%v", random.ShortGUID())
	namespace := fmt.Sprintf("e-%v", random.ShortGUID())
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

	// echo 10.12.6.12 devops.build.com >> /etc/hosts
	ast.Printf("写入ingress域名：%v\n", ingressHost)
	hf, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		ast.Println(err)
	}
	defer hf.Close()
	if _, err := hf.WriteString(fmt.Sprintf("%v %v\n", conf.AccessNode, ingressHost)); err != nil {
		ast.Println(err)
	}

	ast.Println("读取/etc/hosts文件：")
	cmd2 := exec.Command("cat", "/etc/hosts")
	stdout, cmdErr := cmd2.Output()
	if cmdErr != nil {
		ast.Printf("failed to execute echo command: %v\n", cmdErr)
	}
	ast.Println("command result: " + string(stdout))

	// 执行部署
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
		case "FAILED":
			ast.Println("部署失败，请检查日志确认原因。")
			break
		default:
			ast.Printf("部署出现异常，状态为: %v\n", currentStatus)
			panic("部署出现异常")
		}
	}

	//********由于4a域和水土测试域不同，以下验证只能在测试环境运行*********//

	if !strings.Contains(conf.KcHost, "172") {
		// 检查：用户可以通过ingress访问应用
		ast.Println("暂停8秒等待部署充分完成")
		time.Sleep(8 * time.Second)
		ingressURL := "http://" + ingressHost
		checkSiteReq, err := http.NewRequest(http.MethodGet, ingressURL+"/", nil)
		errors.HandleError("生成新请求：访问主页", err)
		checkSiteRes := req.SendRequestAndGetResponse(c, checkSiteReq)
		ast.AssertContainString("应用的/路径应返回正确响应", string(checkSiteRes.Response), "看到这行字，证明你的应用已被成功部署！", checkSiteRes)

		// 检查：用户可以通过访问/env来获取配置好的环境变量，并检查其值应是"helloenv"
		checkEnvReq, err := http.NewRequest(http.MethodGet, ingressURL+"/env", nil)
		errors.HandleError("生成新请求：访问/env", err)
		checkEnvRes := req.SendRequestAndGetResponse(c, checkEnvReq)
		ast.AssertContainString("应用的/env路径应返回设置的环境变量NEWTESTENV的值：helloenv", string(checkEnvRes.Response), "helloenv", checkEnvRes)

		// 检查：用户可以通过访问/host来获取配置好的自定义DNS解析
		checkHostReq, err := http.NewRequest(http.MethodGet, ingressURL+"/host", nil)
		errors.HandleError("生成新请求：访问/host", err)
		checkHostRes := req.SendRequestAndGetResponse(c, checkHostReq)
		ast.AssertContainString("应用的/host路径应返回设置好的host: devops.testhost.com", string(checkHostRes.Response), "devops.testhost.com", checkHostRes)

		// 关联nodePort
		resp = API.GetDeployConfig(c, sp.DeployID)
		dc := API.DeployConfig{}
		errors.UnmarshalAndHandleError(resp.Response, &dc)
		var nodePort int
		for _, v := range dc.Data.Service.Service {
			if v.ServiceType != "NodePort" {
				continue
			}

			if v.NodePort == 0 {
				ast.Println("nodePort关联失败！")
			}

			nodePort = v.NodePort
		}

		// 检查：用户可以通过nodePort访问应用
		checkNodePortReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v:%v/", conf.AccessNode, nodePort), nil)
		errors.HandleError("生成新请求：通过nodePort访问应用首页", err)
		checkNodePortRes := req.SendRequestAndGetResponse(c, checkNodePortReq)
		ast.AssertContainString("用户可以通过nodePort访问应用", string(checkNodePortRes.Response), "看到这行字，证明你的应用已被成功部署！", checkNodePortRes)

		// 检查 JSON configmap
		checkJSONCmReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf(ingressURL+"/configmap/json"), nil)
		errors.HandleError("生成新请求：访问/configmap/json路径", err)
		checkJSONCmRes := req.SendRequestAndGetResponse(c, checkJSONCmReq)
		ast.AssertContainString("JSON configmap应被挂载进容器", string(checkJSONCmRes.Response), "world", checkJSONCmRes)

		// 检查 YML configmap
		checkYMLCmReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf(ingressURL+"/configmap/yml"), nil)
		errors.HandleError("生成新请求：访问/configmap/yml路径", err)
		checkYMLCmRes := req.SendRequestAndGetResponse(c, checkYMLCmReq)
		ast.AssertContainString("YML configmap应被挂载进容器", string(checkYMLCmRes.Response), "world", checkYMLCmRes)

		// 检查 Properties configmap
		checkProCmReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf(ingressURL+"/configmap/properties"), nil)
		errors.HandleError("生成新请求：访问/configmap/properties路径", err)
		checkProCmRes := req.SendRequestAndGetResponse(c, checkProCmReq)
		ast.AssertContainString("Properties configmap应被挂载进容器", string(checkProCmRes.Response), "world", checkProCmRes)
	}

}

func tearDown(c http.Client, ast *assertion.Assertion, sp *param.SuiteParams, exitChan chan assertion.TestResult) {
	ast.PrintTearDownStart()
	// 应无法删除已被构建关联的制品库
	resp := API.DeleteArtifact(c, sp.ArtifactID)
	d1 := API.ItemDeleted{}
	err := json.Unmarshal(resp.Response, &d1)
	errors.HandleError("err unmarshaling ItemDeleted 1 response", err)
	ast.AssertContainString("应无法删除被构建关联的制品库", d1.ErrorInfo, "该制品关联4条", resp)

	//删除构建1
	fmt.Println("删除构建1:", sp.BuildID)
	resp = API.DeleteBuild(c, sp.BuildID)
	d2 := API.ItemDeleted{}
	err = json.Unmarshal(resp.Response, &d2)
	errors.HandleError("err unmarshaling ItemDeleted 2 response", err)
	ast.AssertSuccess("删除构建1", d2.ErrorInfo, resp)

	//删除构建2
	fmt.Println("删除构建2:", sp.DuplicateBuildID)
	resp = API.DeleteBuild(c, sp.DuplicateBuildID)
	err = json.Unmarshal(resp.Response, &d2)
	errors.HandleError("err unmarshaling ItemDeleted 2 response", err)
	ast.AssertSuccess("删除构建2", d2.ErrorInfo, resp)

	//删除制品库
	resp = API.DeleteArtifact(c, sp.ArtifactID)
	d3 := API.ItemDeleted{}
	err = json.Unmarshal(resp.Response, &d3)
	errors.HandleError("err unmarshaling ItemDeleted 1 response", err)
	ast.AssertSuccess("删除制品库", d3.ErrorInfo, resp)

	// 删除部署
	resp = API.DeleteDeploy(c, sp.DeployID)
	gen := API.GeneralResp{}
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
