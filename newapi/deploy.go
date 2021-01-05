package newapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Blackoutta/profari"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/req"
)

const (
	DeployTypeKubernetesDeploy = "KUBERNETES_DEPLOY"
	K8sTimeoutFiveMinutes      = 300
	ImageAuthEnabled           = 1
	SourceTypeImage            = "IMAGE"
	SourceTypeArtifact         = "ARTIFACT"
)

type CreateDeploy struct {
	AppId               int    `json:"appId,omitempty"`
	DeployType          string `json:"deployType,omitempty"`
	Description         string `json:"description,omitempty"`
	EnvId               int    `json:"envId,omitempty"`
	K8sTimeout          int    `json:"k8sTimeout,omitempty"`
	Name                string `json:"name,omitempty"`
	ProjectId           int    `json:"projectId,omitempty"`
	VmGroupId           int    `json:"vmGroupId,omitempty"`
	SourceCreateRequest `json:"sourceCreateRequest,omitempty"`
}

type SourceCreateRequest struct {
	ArtifactId      int    `json:"artifactId,omitempty"`
	ArtifactVersion string `json:"artifactVersion"`
	CredentialId    int    `json:"credentialId,omitempty"`
	ImageAddr       string `json:"imageAddr,omitempty"`
	ImageAuthType   int    `json:"imageAuthType,omitempty"`
	SourceType      string `json:"sourceType,omitempty"`
}

func (r CreateDeploy) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPost, "deploy/base", nil, r)
}

type DeleteDeploy struct {
	DeployID int
}

func (r DeleteDeploy) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodDelete, fmt.Sprintf("deploy/base/%v", r.DeployID), nil, nil)
}

// 编辑k8s部署
func NewEditK8sDeployConfig(deployID int, ingressHost string) *EditK8sDeployConfig {
	return &EditK8sDeployConfig{
		Apply: false,
		KubeConfig: KubeConfig{
			Namespace:      "",
			Cpu:            "1000m",
			Mem:            "512Mi",
			Duplicate:      1,
			PodSuccessFlag: true,
			IsChanged:      true,
			IsEnvNamespace: true,
		},
		Readiness: Readiness{
			Port:                "9005",
			Path:                "/",
			Scheme:              SchemeHTTP,
			PeriodSeconds:       3,
			TimeoutSeconds:      10,
			InitialDelaySeconds: 0,
			ProbeType:           ProbeTypeHttpGet,
			HttpHeaders: struct {
				Platform int `json:"platform"`
			}{Platform: 1},
		},
		Liveness: LiveNess{
			Port:                "9005",
			Path:                "/",
			Scheme:              SchemeHTTP,
			PeriodSeconds:       3,
			TimeoutSeconds:      10,
			InitialDelaySeconds: 0,
			ProbeType:           ProbeTypeHttpGet,
			HttpHeaders: struct {
				Platform int `json:"platform"`
			}{Platform: 1},
		},
		DeployId: deployID,
		ContainerLog: ContainerLog{
			LogPath:   "",
			IsChanged: false,
		},
		ConfigMap: ConfigMap{
			ConfigMap: []ConfigMapConf{
				{
					MountPath: "/data/myconfig/config.json",
					Data:      "{\"hello\", \"world\"}",
				},
				{
					MountPath: "/data/myconfig/config.yml",
					Data:      "hello: world",
				},
				{
					MountPath: "/data/myconfig/config.properties",
					Data:      "hello=world",
				},
			},
		},
		EnvVar: EnvVar{
			EnvVar: struct {
				TestEnv string `json:"TESTENV"`
			}{
				TestEnv: "helloenv",
			},
			IsChanged: true,
		},
		Grpc: Grpc{
			Grpc: []GrpcConf{},
		},
		HostAlias: HostAlias{
			HostAlias: struct {
				Host string `json:"devops.testhost.com"`
			}{
				Host: "6.6.6.6",
			},
			IsChanged: true,
		},
		Ingress: Ingress{
			Ingress: []IngressConf{
				{
					Type:         IngressTypeHTTP,
					Path:         "/",
					Host:         ingressHost,
					InternalPort: "9005",
					Annotations: AnnotationConf{
						WEBSOCKET: "false",
					},
				},
			},
			IsChanged: true,
		},
		IngressTemplate: IngressTemplate{
			Template: "",
		},
		KubeConfigTemplate: KubeConfigTemplate{
			Template: "",
		},
		Service: Service{
			Service: []ServiceConf{
				{
					ProtocolType: ProtocolTypeTCP,
					Port:         "9005",
					ServiceType:  ServiceTypeClusterIP,
				},
				{
					ProtocolType: ProtocolTypeTCP,
					Port:         "9005",
					ServiceType:  ServiceTypeNodePort,
				},
			},
			IsChanged: true,
		},
		ServiceTemplate: ServiceTemplate{
			Template: "",
		},
	}
}

type EditK8sDeployConfig struct {
	Apply              bool               `json:"apply"`
	KubeConfig         KubeConfig         `json:"kubeConfig"`
	Readiness          Readiness          `json:"readiness"`
	Liveness           LiveNess           `json:"liveness"`
	DeployId           int                `json:"deployId"`
	ContainerLog       ContainerLog       `json:"containerLog"`
	ConfigMap          ConfigMap          `json:"configMap"`
	EnvVar             EnvVar             `json:"envVar"`
	Grpc               Grpc               `json:"grpc"`
	HostAlias          HostAlias          `json:"hostAlias"`
	Ingress            Ingress            `json:"ingress"`
	IngressTemplate    IngressTemplate    `json:"ingressTemplate"`
	KubeConfigTemplate KubeConfigTemplate `json:"kubeConfigTemplate"`
	Service            Service            `json:"service"`
	ServiceTemplate    ServiceTemplate    `json:"serviceTemplate"`
}

func (r EditK8sDeployConfig) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPut, fmt.Sprintf("k8s-deploy/config/"), nil, r)
}

type ServiceTemplate struct {
	Template  string `json:"template"`
	IsChanged bool   `json:"isChanged"`
}

type Service struct {
	Service   []ServiceConf `json:"service"`
	IsChanged bool          `json:"isChanged"`
}

const (
	ProtocolTypeTCP      = "TCP"
	ProtocolTypeUDP      = "UDP"
	ServiceTypeClusterIP = "ClusterIP"
	ServiceTypeNodePort  = "NodePort"
)

type ServiceConf struct {
	ProtocolType string `json:"protocolType"`
	Port         string `json:"port"`
	ServiceType  string `json:"serviceType"`
}

type KubeConfigTemplate struct {
	Template  string `json:"template"`
	IsChanged bool   `json:"isChanged"`
}

type IngressTemplate struct {
	Template  string `json:"template"`
	IsChanged bool   `json:"isChanged"`
}

type Ingress struct {
	Ingress   []IngressConf `json:"ingress"`
	IsChanged bool          `json:"isChanged"`
}

const (
	IngressTypeHTTP = "HTTP"
)

type IngressConf struct {
	Type         string         `json:"type"`
	Path         string         `json:"path"`
	Host         string         `json:"host"`
	InternalPort string         `json:"internalPort"`
	Annotations  AnnotationConf `json:"annotations"`
}

type AnnotationConf struct {
	WEBSOCKET string `json:"WEBSOCKET"`
}

type HostAlias struct {
	HostAlias interface{} `json:"hostAlias"`
	IsChanged bool        `json:"isChanged"`
}

type Grpc struct {
	Grpc      []GrpcConf `json:"grpc"`
	IsChanged bool       `json:"isChanged"`
}

type GrpcConf struct {
}

type EnvVar struct {
	EnvVar    interface{} `json:"envVar"`
	IsChanged bool        `json:"isChanged"`
}

type ConfigMap struct {
	ConfigMap []ConfigMapConf `json:"configMap"`
	IsChanged bool            `json:"isChanged"`
}

type ConfigMapConf struct {
	MountPath     string      `json:"mountPath"`
	Data          interface{} `json:"data"`
	ConfigMapName string      `json:"configMapName"`
}

type KubeConfig struct {
	Namespace      string `json:"namespace"`
	Cpu            string `json:"cpu"`
	Mem            string `json:"mem"`
	Duplicate      int    `json:"duplicate"`
	PodSuccessFlag bool   `json:"podSuccessFlag"`
	IsChanged      bool   `json:"isChanged"`
	IsEnvNamespace bool   `json:"isEnvNamespace"`
}

const (
	ProbeTypeHttpGet = "HTTP_GET"
)

type Readiness struct {
	Port                string      `json:"port"`
	Path                string      `json:"path"`
	Scheme              string      `json:"scheme"`
	PeriodSeconds       int         `json:"periodSeconds"`
	TimeoutSeconds      int         `json:"timeoutSeconds"`
	InitialDelaySeconds int         `json:"initialDelaySeconds"`
	ProbeType           string      `json:"probeType"`
	HttpHeaders         interface{} `json:"httpHeaders"`
}

const (
	SchemeHTTP = "http"
)

type LiveNess struct {
	Port                string      `json:"port"`
	Path                string      `json:"path"`
	Scheme              string      `json:"scheme"`
	PeriodSeconds       int         `json:"periodSeconds"`
	TimeoutSeconds      int         `json:"timeoutSeconds"`
	InitialDelaySeconds int         `json:"initialDelaySeconds"`
	ProbeType           string      `json:"probeType"`
	HttpHeaders         interface{} `json:"httpHeaders"`
}

// 执行部署
type ExecuteDeploy struct {
	DeployId int         `json:"deployId"`
	Params   interface{} `json:"params,omitempty"`
}

func (r ExecuteDeploy) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodPatch, "deploy/base/execute", nil, r)
}

// 查询部署详情
type GetDeployHistory struct {
	DeployHistoryID int
}

func (r GetDeployHistory) Compose() (*http.Request, *profari.Record, error) {
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("deploy/history/%d", r.DeployHistoryID), nil, nil)
}

// 分页查询部署列表
type GetDeployList struct {
	ProjectID int
}

func (r GetDeployList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("pageNum", "1")
	q.Add("pageSize", "100")
	q.Add("projectId", strconv.Itoa(r.ProjectID))
	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("deploy/base/page"), q, nil)
}

// 获取实例列表
type GetPodList struct {
	DeployID int
	AppID    int
	EnvID    int
}

func (r GetPodList) Compose() (*http.Request, *profari.Record, error) {
	q := make(url.Values)
	q.Add("appId", strconv.Itoa(r.AppID))
	q.Add("deployEnvIds", fmt.Sprintf("%d:%d", r.DeployID, r.EnvID))

	return req.ComposeNewProfariRequest(http.MethodGet, fmt.Sprintf("deploy/pod/%d/pods", r.AppID), q, nil)
}
