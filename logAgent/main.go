package main

import (
	"fmt"
	"logagent/kafka"
	"logagent/taillog"
	"time"
)

func run() {
	//1、读取日志
	for {
		select {
		case line := <-taillog.ReadChan():
			//2、发送kafka
			kafka.SendToKafka("web_log", line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}

//logAgent入口
func main() {
	//1、初始化kafka连接
	err := kafka.Init([]string{"192.168.66.10:9092"})
	if err != nil {
		fmt.Printf("Kafka initialization failed, err:%v\n", err)
		return
	}
	fmt.Println("Kafka initialization succeeded")
	//2、收集日志
	err = taillog.Init("./my.log")
	if err != nil {
		fmt.Printf("taillog initialization failed,err:%v\n", err)
		return
	}
	fmt.Println("tailLog initialization succeeded")
	run()
}
