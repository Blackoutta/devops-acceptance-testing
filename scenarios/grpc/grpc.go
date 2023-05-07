package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"gitlab.blackoutta.com/devops-acceptance-testing/v1/API"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/grpc/greetpb"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/req"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/assertion"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/conf"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/errors"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/param"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/prep"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/random"
	"google.golang.org/grpc"
)

func RunGrpcTest(exitChan chan assertion.TestResult) {
	// 准备工作
	f, ast, sp, c := prep.SetupTest("GRPC测试套件")
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

	// 创建构建
	buildName := fmt.Sprintf(`b_%v`, random.ShortGUID())
	resp = API.CreateGrpcBuild(c, buildName, sp.ProjectID)
	bc := API.BuildCreated{}
	err := json.Unmarshal(resp.Response, &bc)
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

	// 编辑构建
	sp.ImageTag = "grpc-test"
	resp = API.EditGrpcBuild(c, sp.GitCredentialID, sp.DockerCredentialID, sp.ProjectID, sp.BuildID, buildName, sp.SourceID, sp.ImageTag)
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
	resp = API.UpdateEnvWithGRPC(c, sp.EnvID, namespace, conf.ServerAddr, env, conf.EnvContext, conf.ServerAddr)
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
	resp = API.CreateGrpcDeployment(c, sp.AppID, sp.EnvID, sp.ProjectID, deployName, sp.DockerCredentialID, sp.ImageTag)
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
	resp = API.EditGrpcDeployDetails(c, sp.DeployID)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("编辑部署，加入ingress服务暴露", gen.ErrorInfo, resp)

	// 获取GRPC域名
	resp = API.GetDeployConfig(c, sp.DeployID)
	dc := API.DeployConfig{}
	errors.UnmarshalAndHandleError(resp.Response, &dc)
	grpcDomain := dc.Data.Grpc.Grpc[0].GrpcName + ".grpc.local"
	fmt.Printf("grpc config is: (%s %s) \n", conf.ServerAddr, grpcDomain)

	// 写入GRPC Host
	hf, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		ast.Println(err)
	}

	if _, err := hf.WriteString(fmt.Sprintf("%v %v\n", conf.ServerAddr, grpcDomain)); err != nil {
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

	// 执行构建
	resp = API.RunBuild(c, sp.BuildID, sp.CommitID, "master")
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
		case "BUILDING":
			ast.Println("部署状态为：构建中，等待5秒...")
			time.Sleep(5 * time.Second)
			continue
		case "SUCCESS":
			ast.Println("部署成功!")
			time.Sleep(time.Second)
			deploySuccess = true
		case "FAILED":
			ast.Println("部署失败，请检查日志确认原因。")
		default:
			ast.Printf("部署出现异常，状态为: %v\n", currentStatus)
			panic("部署出现异常")
		}
	}

	// 等待部署充分完成
	ast.SleepWithCounter("等待部署充分完成", 60)

	// 启动GRPC客户端连接GRPC服务并检查接口应能调用
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	cc, err := grpc.DialContext(ctx, grpcDomain+":8080", grpc.WithInsecure())
	if err != nil {
		ast.Printf("could not connect: %v\n", err)
	}
	defer cancel()
	defer cc.Close()

	ccc := greetpb.NewGreetServiceClient(cc)

	// invoking grpc calls
	ast.Println("Starting to do Unary RPC!")
	greq := &greetpb.GreetingRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "yang",
			LastName:  "hu",
		},
	}
	res, err := ccc.Greet(context.Background(), greq)
	if err != nil {
		ast.Println(res.String())
		ast.Printf("error while calling Greet RPC: %v", err)
	}
	ast.AssertStringEqual("验证grpc客户端可调用grpc服务端接口", res.GetResult(), "Hello yang", req.Record{
		Method:     "gRPC Call",
		URL:        grpcDomain + ":8080",
		Body:       nil,
		Response:   []byte(res.GetResult()),
		StatusCode: 0,
	})

}

func tearDown(c http.Client, ast *assertion.Assertion, sp *param.SuiteParams, exitChan chan assertion.TestResult) {
	ast.PrintTearDownStart()

	//删除构建
	resp := API.DeleteBuild(c, sp.BuildID)
	d2 := API.ItemDeleted{}
	err := json.Unmarshal(resp.Response, &d2)
	errors.HandleError("err unmarshaling ItemDeleted 2 response", err)
	ast.AssertSuccess("删除构建", d2.ErrorInfo, resp)

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
