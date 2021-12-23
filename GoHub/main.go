package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

var (
	conf = NewDefaultConfig()
)

func NewDefaultConfig() *Config {
	return &Config{
		//这里填写上传到Bucket的名称
		BucketName: "dqzboy",
	}
}

type Config struct {
	Endpoint   string
	AK         string
	SK         string
	BucketName string
}

func (c *Config) Validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("OSS Endpoint 必须")
	}
	if c.AK == "" {
		return fmt.Errorf("AccessKey ID 必须")
	}
	if c.SK == "" {
		return fmt.Errorf("AccessKey Secret 必须")
	}
	return nil
}

func LoadConfigFromEnv() {
	conf.Endpoint = os.Getenv("ALI_OSS_ENDPOINT")
	conf.AK = os.Getenv("ALI_AK")
	conf.SK = os.Getenv("ALI_SK")
}

func UploadFile(filename string) (downloadURL string, err error) {
	client, err := oss.New(conf.Endpoint, conf.AK, conf.SK)
	if err != nil {
		err = fmt.Errorf("new client error, %s", err)
		return
	}

	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		err = fmt.Errorf("get bucket %s error, %s", conf.BucketName, err)
		return
	}

	err = bucket.PutObjectFromFile(filename, filename)
	if err != nil {
		err = fmt.Errorf("upload file %s error, %s", filename, err)
		return
	}

	// 生成下载链接
	return bucket.SignURL(filename, oss.HTTPGet, 60*60*24*3)
}

//定义文件变量
var (
	fileName 	string
	help		bool
)

//声明CLI的参数
func init() {
	flag.StringVar(&fileName, "f","","请输入需要上传的文件的路径")
	flag.BoolVar(&help,"help",false,"打印本工具的使用说明")
}

//命令执行 帮助提示
func usage() {
	fmt.Fprintf(os.Stderr, `cloud-station version: 0.0.1
Usage: cloud-station [-h] -f <uplaod_file_path>
Options:
`)
	//如果选项错误，返回默认提示信息
	flag.PrintDefaults()
}

// LoadArgsFromCLI 负责接收用户传入的参数
func LoadArgsFromCLI() {
	//通过标准库flag parse解析CLI参数
	flag.Parse()
	//打印说明信息
	if help {
		usage()
		os.Exit(0)
	}
}

func main() {
	LoadConfigFromEnv()

	//校验配置
	if err := conf.Validate();err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//接收用户参数
	LoadArgsFromCLI()
	//上传文件
	downloadURL, err := UploadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//正常退出
	fmt.Printf("文  件: %s 上传成功\n", fileName)

	// 打印下载链接
	fmt.Printf("下载链接: %s\n", downloadURL)
	fmt.Println("\n注意: 文件下载有效期为1天, 保存时间为3天, 请及时下载")
	os.Exit(0)
}
