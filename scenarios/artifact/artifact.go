package artifact

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gitlab.blackoutta.com/devops-acceptance-testing/v1/API"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/assertion"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/errors"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/param"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/prep"
	"gitlab.blackoutta.com/devops-acceptance-testing/v1/util/random"
)

func RunArtifactTest(exitChan chan assertion.TestResult) {
	// 准备工作
	f, ast, sp, c := prep.SetupTest("制品库测试套件")
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
	fmt.Println("artifact ID is:", sp.ArtifactID)

	// 获取制品库下载地址
	var ad API.ArtifactDetail
	resp = API.GetArtifactByID(c, sp.ArtifactID)
	err = json.Unmarshal(resp.Response, &ad)
	errors.HandleError("err unmarshaling ArtifactObtained response", err)
	ast.AssertSuccess("获取制品库下载地址", ad.ErrorInfo, resp)
	downloadAddress := ad.Data.DownloadAddress
	fields := strings.Split(downloadAddress, "/")
	downloadID := fields[len(fields)-1]
	fmt.Println(fields)
	fmt.Println("downloadID is:", downloadID)
	// 制品库测试

	// 上传小于100M的ZIP包可以成功
	version := random.ShortGUID()
	resp = API.UploadArtifact(c, version, sp.ArtifactID, "data/test_upload.zip", "application/x-zip-compressed")
	gen := API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("传小于100M的ZIP包可以成功", gen.ErrorInfo, resp)

	// 检查发布记录
	resp = API.GetArtifactUploadRecord(c, sp.ArtifactID)
	ulr := API.UploadRecords{}
	errors.UnmarshalAndHandleError(resp.Response, &ulr)
	ast.AssertStringEqual("检查应有一条该版本制品的发布记录", ulr.Data.Data[0].Version, version, resp)
	ast.AssertIntegerEqual("检查发布记录的发布方式是：上传", ulr.Data.Data[0].UploadType, 1, resp)

	// 验证下载并解压制品可以成功
	resp = API.DownloadArtifact(c, downloadID, "", "")
	dd := string(resp.Response)
	ast.AssertContainString("验证下载的制品中的数据完好（包含upload success）字符", dd, "upload success", resp)

	// 检查下载记录
	resp = API.GetArtifactDownloadRecord(c, sp.ArtifactID)
	dlr := API.DownloadRecords{}
	errors.UnmarshalAndHandleError(resp.Response, &dlr)
	ast.AssertStringEqual("检查应有一条该版本制品的下载记录", dlr.Data.Data[0].Version, version, resp)

	// 修改制品库，加入鉴权
	artifactToken := "12345678"
	resp = API.EditArtifact(c, sp.ArtifactID, artifactName, "1", artifactToken, false)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("修改制品库，加入鉴权", gen.ErrorInfo, resp)

	// 传入错误的token
	resp = API.DownloadArtifact(c, downloadID, "wrong token", "")
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertIntegerEqual("检查手动下载制品时，在header传入错误token会导致无法下载制品", gen.ErrorCode, 403, resp)
	fmt.Println(string(resp.Response))

	// 不传token
	resp = API.DownloadArtifact(c, downloadID, "", "")
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertIntegerEqual("检查手动下载制品时，不在header传入token会导致无法下载制品", gen.ErrorCode, 401, resp)
	fmt.Println(string(resp.Response))

	// 传入正确token
	resp = API.DownloadArtifact(c, downloadID, artifactToken, "")
	ast.AssertIntegerNotEqual("检查手动下载制品时，在header传入正确token后，可以下载制品", len(resp.Response), 0, resp)

	// 上传一个新版本制品
	newVersion := "new_version"
	resp = API.UploadArtifact(c, newVersion, sp.ArtifactID, "data/test_upload.zip", "application/x-zip-compressed")
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("上传一个新版本制品", gen.ErrorInfo, resp)

	// 验证可以下载指定版本制品
	resp = API.DownloadArtifact(c, downloadID, artifactToken, newVersion)
	ddd := string(resp.Response)
	ast.AssertContainString("验证在有鉴权的情况下，通过指定版本下载的历史版本制品中的数据完好（包含upload success）字符", ddd, "upload success", resp)

	// 禁用制品库
	resp = API.EditArtifact(c, sp.ArtifactID, artifactName, "1", artifactToken, true)
	gen = API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("禁用制品库", gen.ErrorInfo, resp)

	// 禁用制品库后无法上传
	version = random.ShortGUID()
	resp = API.UploadArtifact(c, version, sp.ArtifactID, "data/test_upload.zip", "application/x-zip-compressed")
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertIntegerEqual("禁用制品库后禁止上传", gen.ErrorCode, 2009, resp)

	// 禁用制品库后无法下载
	resp = API.GetArtifactDownloadRecord(c, sp.ArtifactID)
	errors.UnmarshalAndHandleError(resp.Response, &dlr)
	ast.AssertIntegerEqual("禁用制品库后禁止下载", gen.ErrorCode, 2009, resp)

}

func tearDown(c http.Client, ast *assertion.Assertion, sp *param.SuiteParams, exitChan chan assertion.TestResult) {
	ast.PrintTearDownStart()

	//删除制品库
	resp := API.DeleteArtifact(c, sp.ArtifactID)
	d3 := API.ItemDeleted{}
	err := json.Unmarshal(resp.Response, &d3)
	errors.HandleError("err unmarshaling ItemDeleted 1 response", err)
	ast.AssertSuccess("删除制品库", d3.ErrorInfo, resp)

	// 删除项目
	resp = API.DeleteProject(c, sp.ProjectID)
	gen := API.GeneralResp{}
	errors.UnmarshalAndHandleError(resp.Response, &gen)
	ast.AssertSuccess("删除项目", gen.ErrorInfo, resp)

	ast.PrintTearDownEnd()

	// 判断测试成功与否
	ast.CheckSuiteResult(exitChan)
}
