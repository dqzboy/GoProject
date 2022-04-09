package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 专门往kafka写日志的模块
func main() {
	config := sarama.NewConfig()
	//tailf 包使用
	config.Producer.RequiredAcks = sarama.WaitForAll //WaitForAll等待所有节点响应
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	msg := &sarama.ProducerMessage{}
	msg.Topic = "web_log"                                  //topic名称
	msg.Value = sarama.StringEncoder("this ia a test log") //使用sarama.StringEncoder序列化
	//连接kafka
	clinet, err := sarama.NewSyncProducer([]string{"192.168.66.10:9092"}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	defer clinet.Close()
	//发送消息
	pid, offset, err := clinet.SendMessage(msg)
	if err != nil {
		fmt.Println("sed msg failed, err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}
