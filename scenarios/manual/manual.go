package manual

import (
	"encoding/json"
	"fmt"

	"gitlab.blackoutta.com/devops-acceptance-testing/v1/API"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/assertion"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/conf"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/errors"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/prep"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/random"
)

func RunManualTest(exitChan chan assertion.TestResult) {
	// 准备工作
	f, ast, sp, c := prep.SetupTest("手工测试")
	defer f.Close()

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

	//关联Commit ID
	resp = API.GetCommitID(c, sp.BuildID)
	bb := API.BeforeBuild{}
	err = json.Unmarshal(resp.Response, &bb)
	errors.HandleError("err unmarshaling BeforeBuild response", err)
	ast.AssertSuccess("关联Commit ID", bb.ErrorInfo, resp)

	sp.CommitID = bb.Data.SpecificBranchOrTag.CommitID

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
}
