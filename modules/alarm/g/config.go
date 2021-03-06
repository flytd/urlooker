package g

import (
	"fmt"
	"log"
	"sync"

	"github.com/toolkits/file"
	"gopkg.in/yaml.v2"
)

type GlobalConfig struct {
	Debug  bool          `yaml:"debug"`
	Remain int           `yaml:"remain"`
	Rpc    *RpcConfig    `yaml:"rpc"`
	Web    *WebConfig    `yaml:"web"`
	Worker *WorkerConfig `yaml:"worker"`
	Smtp   *SmtpConfig   `yaml:"smtp"`
	WeChat *WeChatConfig `yaml:"wechat"`
}

type MysqlConfig struct {
	Addr string `yaml:"addr"`
	Idle int    `yaml:"idle"`
	Max  int    `yaml:"max"`
}

type RpcConfig struct {
	Listen string `yaml:"listen"`
}

type WebConfig struct {
	Addrs    []string `yaml:"addrs"`
	Timeout  int      `yaml:"timeout"`
	Interval int      `yaml:"interval"`
}

type WorkerConfig struct {
	Sms  int `yaml:"sms"`
	Mail int `yaml:"mail"`
	WeChat int `yaml:"wechat"`
}

type SmtpConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	Tls      bool   `yaml:"tls"`
}

type WeChatConfig struct {
	ToParty    string `json:"toparty"`
	AgentId    int    `json:agentid`
	CorpId     string `json:corpid`
	CorpSecret string `json:corpsecret`
}

var (
	Config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is not exists", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c GlobalConfig
	err = yaml.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	configLock.Lock()
	defer configLock.Unlock()
	Config = &c

	if Config.Remain < 10 {
		Config.Remain = 30
	}

	log.Println("load configuration file", cfg, "successfully")
	return nil
}
