/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"management_backend/src/config"
	loggers "management_backend/src/logger"
	"management_backend/src/router"
	"management_backend/src/utils"
)

const (
	configParam = "config"
)

//go:embed resources/*
var staticContent embed.FS

func main() {
	// 解析命令行参数：config
	var confYml string
	confYml = configYml()
	if confYml == "" {
		fmt.Println("can not find param [--config], will use default")
		confYml = "configs"
		// 判断confYml下文件是否存在
		ok, err := utils.PathExists(confYml + string(os.PathSeparator) + "config.yml")
		if err != nil {
			panic(err)
		}
		if !ok {
			fmt.Println("can not load config.yml, exit")
			return
		}
	}
	config.InitConfig(confYml, config.GetConfigEnv())
	// 初始化日志配置
	loggers.SetLogConfig(config.GlobalConfig.LogConf)
	// 初始化数据库配置
	// connection.InitDbConn(config.GlobalConfig.DBConf)
	// http-server启动

	router.HttpServe(config.GlobalConfig.WebConf, staticContent)
}

func configYml() string {
	configPath := flag.String(configParam, "", "config.yml's path")
	flag.Parse()
	return *configPath
}
