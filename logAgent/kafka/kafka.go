package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 专门往kafka写日志的模块
var (
	client sarama.SyncProducer //声明一个全局的连接kafka的生产者

)

// Init 初始化
func Init(addrs []string) (err error) {
	config := sarama.NewConfig()
	//tailf 包使用
	config.Producer.RequiredAcks = sarama.WaitForAll        //WaitForAll等待所有节点响应
	config.Producer.Partitioner = sarama.NewHashPartitioner //新选出一个partition
	config.Producer.Return.Successes = true                 //成功发送的消息将在succcess channel返回

	//连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	return
}

func SendToKafka(topic, data string) {
	//构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic                      //topic名称
	msg.Value = sarama.StringEncoder(data) //使用sarama.StringEncoder序列化

	//发送到kafka
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("sed msg failed, err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}
