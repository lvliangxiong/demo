package conf

import (
	"github.com/lvliangxiong/demo/rabbitmq/publish_subscribe/util"
	"github.com/spf13/viper"
)

var v *viper.Viper = viper.New()

func init() {
	util.FailOnError(initConfig(), "fail to init config")
}

func initConfig() error {
	confFile := "conf/conf.yaml"
	localConfFile := "conf/conf.local.yaml"

	v.SetConfigFile(confFile)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if util.IsFile(localConfFile) {
		v.SetConfigFile(localConfFile)
		if err := v.MergeInConfig(); err != nil {
			return err
		}
	}

	return nil
}

func GetString(name string) string {
	return v.GetString(name)
}
