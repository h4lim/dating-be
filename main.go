package main

import (
	"dating-be/common"
	"dating-be/docs"
	"dating-be/infra"
	"fmt"
	"os"
)

//
// @title DATING APPS API
// @version 1.0
// @description Api swagger collection for DATING APPS API.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Abdul Halim
// @license.url https://github.com/h4lim
func main() {

	// toml config
	configModel := common.ConfigModel{
		FileName: "config.toml",
	}

	config := common.NewConfig(configModel)
	if err := config.Open(); err != nil {
		fmt.Println("error config setup", *err)
		os.Exit(1)
	}

	// message config
	messageModel := common.MessageModel{
		Path:     "",
		FileName: common.ConfigString["message_json"],
	}
	messageConfig := common.NewMessageConfig(messageModel)
	if err := messageConfig.Setup(); err != nil {
		fmt.Println("error message setup", fmt.Sprintf("%v", *err))
		os.Exit(1)
	}

	// swagger config
	docs.SwaggerInfo.Host = common.ConfigString["swagger_host"]
	docs.SwaggerInfo.BasePath = common.ConfigString["swagger_basePath"]
	docs.SwaggerInfo.Schemes = []string{common.ConfigString["swagger_schemes"]}

	// zap config
	zapModel := common.ZapModel{
		ServiceName: common.ConfigString["service_name"],
		Mode:        common.ConfigString["mode"],
		OutputPath:  "log",
	}
	zapConfig := common.NewZapConfig(zapModel)
	if err := zapConfig.ZapSetup(); err != nil {
		fmt.Println("error zap setup", fmt.Sprintf("%v", *err))
		os.Exit(1)
	}

	// db setup
	dbModel := common.GormContext{
		Driver:   common.ConfigString["db_driver"],
		Port:     common.ConfigString["db_port"],
		Host:     common.ConfigString["db_host"],
		Username: common.ConfigString["db_username"],
		Password: common.ConfigString["db_password"],
		DBName:   common.ConfigString["db_name"],
	}
	gormDB := common.NewGormDB(dbModel)
	if _, err := gormDB.Open(); err != nil {
		fmt.Println("error gorm setup", fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	if err := gormDB.Migrate(); err != nil {
		fmt.Println("error gorm auto migrate", fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	if err := gormDB.Seeder(); err != nil {
		fmt.Println("error gorm auto seeder", fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	redisModel := common.RedisModel{
		Host:     common.ConfigString["redis_host"],
		Port:     common.ConfigString["redis_port"],
		Password: common.ConfigString["redis_password"],
	}

	if err := common.NewRedisConfig(redisModel).Open(); err != nil {
		fmt.Println("error open redis", fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	// gin config
	serverModel := infra.ServerModel{
		Port: common.ConfigString["port"],
	}
	if err := infra.NewServerConfig(serverModel).Run(); err != nil {
		fmt.Println("error zap setup", fmt.Sprintf("%v", *err))
		os.Exit(1)
	}

}
