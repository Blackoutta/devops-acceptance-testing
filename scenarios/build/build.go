package build

import (
	"fmt"
	"os"

	"github.com/Blackoutta/profari"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/API"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/newapi"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/random"
)

type BuildTest struct {
	BuildId      int
	Name         string
	ErrChan      chan error
	SkipTeardown bool
	*profari.Client
	logFile *os.File
	*newapi.SetupParams
}

func (t *BuildTest) GetName() string {
	return t.Name
}

func (t *BuildTest) GetErrChan() chan error {
	return t.ErrChan
}

func (t *BuildTest) Run() {
	var err error
	// You must initialize the profari client before test starts
	t.Client, t.logFile, err = profari.NewClient(t.Name, t.ErrChan)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 准备工作
	t.SetupParams = newapi.BasicSetup(t.Client)
	fmt.Printf("%+v\n", t.SetupParams)

	// create build
	var build API.BuildCreated
	t.Send(&newapi.CreateBuild{
		Name:      "build_" + random.ShortGUID(),
		ProjectId: t.ProjectID,
		BuildType: newapi.BuildTypeGolang,
		Source: newapi.BuildSource{
			ProjectId:      t.ProjectID,
			AuthPolicy:     newapi.AuthPolicyCREDENTIAL,
			CredentialId:   t.GitID,
			Type:           "GIT",
			Address:        "http://gitlab.blackoutta.com/devops-test-httpserver.git",
			BranchPolicy:   newapi.BranchPolicySpecific,
			SpecificBranch: "smoke_build",
		},
	}).DecodeJSON(&build).AssertContainString("创建构建应成功", build.ErrorInfo, "success")
	t.BuildId = build.Data

	// get build
	var builds API.Builds
	t.Send(&newapi.GetBuildList{
		ProjectID: t.ProjectID,
	}).DecodeJSON(&builds).AssertEqualInt("查询构建列表应成功并有1条数据", len(builds.Data.Data), 1)

	// add collection
	t.Send(&newapi.AddBuildCollection{
		BuildId: t.BuildId,
	}).AssertContainString("收藏构建应成功", build.ErrorInfo, "success")

	// get collection build list
	var collectionBuilds API.Builds
	t.Send(&newapi.GetBuildCollectionList{
		ProjectId: t.ProjectID,
	}).DecodeJSON(&collectionBuilds).AssertEqualInt("收藏的构建应为1个", len(collectionBuilds.Data.Data), 1)

	// delete collection
	t.Send(&newapi.DeleteBuildCollection{
		BuildId: t.BuildId,
	}).AssertContainString("删除收藏构建应成功", t.Resp, "success")

	// get collection build list again
	t.Send(&newapi.GetBuildCollectionList{
		ProjectId: t.ProjectID,
	}).DecodeJSON(&collectionBuilds).AssertEqualInt("删除之后收藏的构建应为0个", len(collectionBuilds.Data.Data), 0)

	// get type image
	var imagePage API.ImagePage
	pageSize := 5
	t.Send(&newapi.GetImageList{
		Type:     newapi.GolangImage,
		PageSize: pageSize,
	}).DecodeJSON(&imagePage).AssertContainString("查询镜像列表应成功", imagePage.ErrorInfo, "success")
	if imagePage.Data.Total >= pageSize {
		// 总数超过一页: 查询pageSize
		t.AssertEqualInt("查询的镜像数量应一致", len(imagePage.Data.Data), pageSize)
	} else {
		// 总数少于一页: 查询total
		t.AssertEqualInt("查询的镜像数量应一致", len(imagePage.Data.Data), imagePage.Data.Total)
	}

	// token such as:2df8131e-be1e-4a64-a22a-ae96cb4183d4
	var token API.StringData
	t.Send(&newapi.GenHookToken{}).DecodeJSON(&token).AssertEqualInt("生成的token应为36位", len(token.Data), 36)
}

func (t *BuildTest) Teardown() {
	defer t.logFile.Close()
	defer t.EndTest()

	if t.SkipTeardown == true {
		t.Println("Skipping Teardown...")
		return
	}

	// delete build
	t.Send(newapi.DeleteBuild{
		BuildID: t.BuildId,
	}).AssertContainString("删除构建应成功", t.Resp, "success")

	// clean setup data
	t.SetupParams.Teardown(t.Client)
}
