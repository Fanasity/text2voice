package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"aiServer/model"

	"github.com/Shopify/sarama"
)

const TOPIC string = "aiTask"

// kafka consumer
func Consumer(cfg model.KafkaConfig) (consumers []sarama.PartitionConsumer, err error) {
	consumer, err := sarama.NewConsumer(cfg.Gateway, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(cfg.TaskTopic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		var pc sarama.PartitionConsumer
		pc, err = consumer.ConsumePartition(cfg.TaskTopic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		consumers = append(consumers, pc)
	}
	return
}

func ConsumerMeesage(ctx context.Context, pc sarama.PartitionConsumer, newTask func(ctx context.Context, msg, key []byte) (Task, error)) {
	for {
		select {
		case msg := <-pc.Messages():
			task, _ := newTask(ctx, msg.Value, msg.Key)
			if task != nil {
				task.Handle()
			}
		}
	}

}

func Produce(cfg model.KafkaConfig, topic string, info interface{}) error {
	if topic == "" {
		topic = cfg.ResultTopic
	}
	config := sarama.NewConfig()
	// config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	// config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true // 成功交付的消息将在success channel返回
	// 连接kafka
	client, err := sarama.NewSyncProducer(cfg.Gateway, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return err
	}
	defer client.Close()

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	data, _ := json.Marshal(info)
	msg.Value = sarama.ByteEncoder(data)
	// 发送消息
	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	return nil
}
