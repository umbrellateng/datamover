/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 4:02 下午
 */
package kafka

import (
	"fmt"
	"sync"

	"core.bank/datamover/log"
	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
)

var (
	from   string
	target string
)

var onlineCmd = &cobra.Command{
	Use: "online",
	Short: "move kafka data from source cluster target the target cluster",
	Args: cobra.NoArgs,
	Run: moveCommandFunc,
}

func init() {
	onlineCmd.Flags().StringVarP(&from, "from", "f", "", "source kafka cluster url")
	onlineCmd.Flags().StringVarP(&target, "target", "t", "", "target kafka cluster url")

	_ = onlineCmd.MarkFlagRequired("from")
	_ = onlineCmd.MarkFlagRequired("target")
}

func moveCommandFunc(cmd *cobra.Command, args []string) {

	err := kafkaOnlineMove(from, target)
	if err != nil {
		log.Logger.Error("redis data migration error: " + err.Error())
		return
	}

	fmt.Println()
	log.Logger.Info("redis data migration on success!")
}

func kafkaOnlineMove(from, target string) error {
	// 创建源集群客户端
	client, err := sarama.NewClient([]string{from}, nil)
	if err != nil {
		return fmt.Errorf("kafka sarama.NewClient error: " + err.Error())
	}
	defer client.Close()

	// 获取源集群所有topic和分区信息
	topics, err := client.Topics()
	if err != nil {
		return fmt.Errorf("kafka client.Topics error: " + err.Error())
	}

	// 创建一个等待组
	var wg sync.WaitGroup

	// 对每个topic和分区进行数据迁移
	for _, topic := range topics {
		partitions, err := client.Partitions(topic)
		if err != nil {
			return fmt.Errorf("kafka client.Partition error: " + err.Error())
		}
		for _, partition := range partitions {
			wg.Add(1) // 增加等待组计数
			go func(topic string, partition int32) {
				defer wg.Done() // 减少等待组计数
				migrate(from, target, topic, partition) // 调用迁移函数
			}(topic, partition)
		}
	}
	wg.Wait() // 等待所有迁移函数完成

	return nil

}

// 迁移函数，接收topic和partition作为参数
func migrate(from, target, topic string, partition int32) {
	log.Logger.Info("kafka migrating topic %s partition %d\n", topic, partition)

	// 创建生产者
	producer, err := sarama.NewSyncProducer([]string{target}, nil)
	if err != nil {
		log.Logger.Error("sarama.NewSyncProducer error: " + err.Error())
		return
	}
	defer producer.Close()

	// 创建消费者
	consumer, err := sarama.NewConsumer([]string{from}, nil)
	if err != nil {
		log.Logger.Error("sarama.NewConsumer error: " + err.Error())
		return
	}
	if consumer == nil {
		log.Logger.Error("kafka consumer is nil")
		return
	}
	defer consumer.Close()


	// 订阅源集群主题和分区
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		log.Logger.Error("kafka consumer.ConsumePartition error: " + err.Error())
		return
	}
	defer partitionConsumer.Close()
	// 从源集群读取消息并发送到目标集群
	for message := range partitionConsumer.Messages() {
		fmt.Printf("Consumed message: topic=%s, partition=%d, offset=%d, key=%s, value=%s\n", message.Topic,
			message.Partition, message.Offset, string(message.Key), string(message.Value))
		msg := &sarama.ProducerMessage{
			Topic:     topic,
			Key:       sarama.StringEncoder(message.Key),
			Value:     sarama.StringEncoder(message.Value),
			Partition: partition,
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Logger.Error("kafka producer.SendMessage error: " + err.Error())
		} else {
			log.Logger.Info("kafka produced message: topic=%s, partition=%d, offset=%d\n", topic, partition, offset)
		}
	}

	log.Logger.Info("kafka migrated topic %s partition %d\n", topic, partition)
}

