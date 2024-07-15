// Package zconf 该包负责相关配置文件定义及加载方式
package zconf

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
)

// Conf 全局变量，用于需要取配置的地方
var Conf *Config

type Config struct {
	/*
		Server
	*/
	Host    string // 当前服务器的监听IP
	Port    int    // 当前服务器的监听端口
	Name    string // 当前服务器的名称标识
	Network string

	/*
		myzinx
	*/
	Version          string // 当前myzinx框架的版本
	MaxConnection    int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前myzinx框架数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作worker的goroutine数量
	MaxWorkerTaskLen uint32 // myzinx框架的每个worker对应的消息队列的任务的数量最大值
}

// Reload 用于从配置文件zinx.json加载自定义的参数
func (c *Config) Reload(path string) {
	_, err := os.Stat(path)
	if err != nil {
		// 配置文件不存在的话，可以直接返回，用初始化值的默认值即可
		if os.IsNotExist(err) {
			fmt.Println("[RELOAD] The config file is not exist, the config file path is:", path)
			return
		}
		panic("[ERROR] Read config file is error")
	}

	fileData, err := os.ReadFile(path)
	// 将json文件解析到Config中
	err = json.Unmarshal(fileData, &Conf)
	if err != nil {
		panic("[ERROR] Config json unmarshal error")
	}
}

func (c *Config) Show() {
	objVal := reflect.ValueOf(c).Elem()
	objType := reflect.TypeOf(*c)

	fmt.Println("===== myzinx config show =====")
	for i := 0; i < objVal.NumField(); i++ {
		field := objVal.Field(i)
		typeField := objType.Field(i)

		fmt.Printf("%s: %v\n", typeField.Name, field.Interface())
	}
	fmt.Println("==============================")
}

// 初始化当前的Config对象
func init() {
	// 如果配置文件没有加载，赋默认值
	Conf = &Config{
		Name:             "MyServer",
		Version:          "v0.4",
		Network:          "tcp",
		Host:             "0.0.0.0",
		Port:             8888,
		MaxConnection:    1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   uint32(runtime.NumCPU()),
		MaxWorkerTaskLen: 1024,
	}
	// 从conf/zinx.json加载用户自定义的参数
	Conf.Reload("../test/myzinx.json")
}
