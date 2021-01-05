package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v4"
	"gitlab.onenet.com/huyangyi/devops-acceptance-testing/v1/util/errors"
)

var (
	Host, WsHost, Username, Password, KcHost, KcRealm, KcClientID, AccessNode, UserID, ServerAddr, EnvFile, EnvContext string
	cf                                                                                                                 = flag.String("cf", "config-test.json", "config file")

	// 缓存一个token
	t   time.Time
	jwt *gocloak.JWT
	err error
	mu  *sync.Mutex
)

func init() {
	configFile, ok := os.LookupEnv("CONFIGFILE")
	if !ok {
		log.Printf("config file not found or env var CONFIGFILE not set! using flag value: %s(defaults to config-test.json)\n", *cf)
		configFile = *cf
	}

	conf := make(map[string]string)
	bs, err := ioutil.ReadFile(configFile)
	errors.HandleError("err reading config file", err)
	if err := json.Unmarshal(bs, &conf); err != nil {
		log.Fatalf("err unmarshaling config: %v\n", err)
	}
	Host = conf["host"]
	WsHost = conf["wsHost"]
	Username = conf["username"]
	Password = conf["password"]
	KcHost = conf["kcHost"]
	KcRealm = conf["kcRealm"]
	KcClientID = conf["kcClientID"]
	AccessNode = conf["accessNode"]
	UserID = conf["userId"]
	ServerAddr = conf["serverAddr"]
	EnvFile = conf["envFile"]
	EnvContext = conf["envContext"]

	mu = &sync.Mutex{}
	getToken()
}

func GetToken() string {
	if jwt == nil {
		getToken()
	}

	if isTokenExpire(jwt) {
		refreshToken()
	}

	return jwt.AccessToken
}

// getToken not lock for get jwt
// 多次重复覆盖jwt没关系 - 不加锁
func getToken() {
	// 初始化 token
	t = time.Now()

	kc := gocloak.NewClient(KcHost)
	jwt, err = kc.Login(KcClientID, "", KcRealm, Username, Password)
	if err != nil {
		log.Panicf("get token from kc err: %v", err)
	}
}

func refreshToken() {
	mu.Lock()
	defer mu.Unlock()
	if !isTokenExpire(jwt) {
		return
	}

	getToken()
}

// isTokenExpire token 是否超时 - 提前30s刷新 (现在token超时时间设置比较长,基本不会出现超时)
func isTokenExpire(jwt *gocloak.JWT) bool {
	return int64(jwt.ExpiresIn) > time.Now().Unix()-t.Unix()+30
}

func ReadEnvFile() string {
	bs, err := ioutil.ReadFile(EnvFile)
	if err != nil {
		log.Fatalf("error while reading file: %s", EnvFile)
	}
	return string(bs)
}
