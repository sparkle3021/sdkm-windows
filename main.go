package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sdkm/jdk"
	"sdkm/model"
)

const (
	configFileName = "\\conf\\config.json"
	jdkVersionData = "\\conf\\jdk_version.json"
	helpContent    = "sdkm ls # 列出已安装的JDK版本\r" +
		"sdkm use <version> # 切换JDK版本\r" +
		"sdkm install <version> # 安装指定版本的JDK\r" +
		"sdkm remove <version> # 卸载指定版本的JDK\r" +
		"sdkm available <jdk type> # 列出可用的JDK版本 "
)

func main() {
	// 读取配置文件
	config, err := readConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	// 处理命令行参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: mj <command>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "ls":
		jdk.ListLocalJdkVersions(config)
	case "use":
		if len(os.Args) < 3 {
			fmt.Println("Usage: mj use <version>")
			os.Exit(1)
		}
		version := os.Args[2]
		jdk.UseVersion(config, version)
	case "available":
		jdkTypes, _ := readJdkVersionConfig()
		jdk.ListJdkVersions(jdkTypes, os.Args[2])
	case "install":
		jdkTypes, _ := readJdkVersionConfig()
		version := os.Args[2]
		jdk.Install(jdkTypes, version, config.JDKDir)
	case "remove":
		jdk.RemoveJdk(os.Args[2], config.JDKDir)
	case "help":
		fmt.Println(helpContent)
	default:
		fmt.Println(helpContent)
		os.Exit(1)
	}
}

func exePath() string {
	exePath, _ := os.Executable()
	realExePath, _ := filepath.EvalSymlinks(exePath)
	exeDir := filepath.Dir(realExePath)
	return exeDir
}

// readConfig 从配置文件中读取配置信息
func readConfig() (model.Config, error) {
	config := model.Config{}
	// 读取配置文件
	file, err := os.Open(exePath() + configFileName)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// 解析 JSON 配置
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// readConfig 从配置文件中读取配置信息
func readJdkVersionConfig() ([]model.JdkType, error) {
	var jdkTypes []model.JdkType
	// 读取配置文件
	file, err := os.Open(exePath() + jdkVersionData)
	if err != nil {
		return jdkTypes, err
	}
	defer file.Close()

	// 解析 JSON 配置
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jdkTypes)
	if err != nil {
		return jdkTypes, err
	}

	return jdkTypes, nil
}
