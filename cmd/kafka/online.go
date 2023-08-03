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
		log.Logger.Error("kafka data migration error: " + err.Error())
		return
	}

	fmt.Println()
	log.Logger.Info("kafka data migration on success!")
}

func kafkaOnlineMove(from, target string) error {
	// 创建源集群客户端
	client, err := sarama.NewClient([]string{from}, nil)
	if err != nil {
		return fmt.Errorf("kafka sarama.NewClient error: " + err.Error())
	}
	defer client.Close()


	// 创建生产者
	producer, err := sarama.NewSyncProducer([]string{target}, nil)
	if err != nil {
		return fmt.Errorf("sarama.NewSyncProducer error: " + err.Error())
	}
	defer producer.Close()

	// 创建消费者
	consumer, err := sarama.NewConsumer([]string{from}, nil)
	if err != nil {
		return fmt.Errorf("sarama.NewConsumer error: " + err.Error())
	}
	if consumer == nil {
		return fmt.Errorf("kafka consumer is nil")
	}
	defer consumer.Close()

	// 创建一个等待组
	var wg sync.WaitGroup

	// 获取源集群所有topic和分区信息
	topics, err := client.Topics()
	if err != nil {
		return fmt.Errorf("kafka client.Topics error: " + err.Error())
	}
	// 对每个topic和分区进行数据迁移
	for _, topic := range topics {
		if topic == "__consumer_offsets" {
			continue
		}
		partitions, err := client.Partitions(topic)
		if err != nil {
			return fmt.Errorf("kafka client.Partition error: " + err.Error())
		}
		for _, partition := range partitions {
			wg.Add(1) // 增加等待组计数
			go func(topic string, partition int32) {
				defer wg.Done() // 减少等待组计数
				migrate(producer, consumer, topic, partition) // 调用迁移函数
			}(topic, partition)
		}
	}
	wg.Wait() // 等待所有迁移函数完成

	return nil

}

// 迁移函数，接收topic和partition作为参数
func migrate(producer sarama.SyncProducer, consumer sarama.Consumer, topic string, partition int32) {
	// 订阅源集群主题和分区
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		log.Logger.Error("kafka consumer.ConsumePartition error: " + err.Error())
		return
	}
	defer partitionConsumer.Close()
	// 从源集群读取消息并发送到目标集群
	var sum int64 = 0
	waterOffset := partitionConsumer.HighWaterMarkOffset()
	log.Logger.Info("topic: %s, partition: %d, high water offset: %d", topic, partition, waterOffset)
	if waterOffset == 0 {
		log.Logger.Warning("topic: %s, partition: %d, water offset: %d, do not need to read for loop", topic, partition, waterOffset)
		return
	}

	for message := range partitionConsumer.Messages() {
		//log.Logger.Info("Consumed message: topic=%s, partition=%d, offset=%d, key=%s, value=%s", message.Topic,
		//	message.Partition, message.Offset, string(message.Key), string(message.Value))
		msg := &sarama.ProducerMessage{
			Topic:     topic,
			Key:       sarama.StringEncoder(message.Key),
			Value:     sarama.StringEncoder(message.Value),
			Partition: partition,
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Logger.Error("kafka producer.SendMessage error: " + err.Error() + ", topic=%s, partition=%d, offset=%d", topic, partition, offset)
		} else {
			log.Logger.Info("kafka produced message: topic=%s, partition=%d, offset=%d, sum=%d, water offset=%d", topic, partition, offset, sum, waterOffset)
		}
		sum++
		if sum == partitionConsumer.HighWaterMarkOffset() {
			break
		}
	}

	log.Logger.Info("kafka migrated topic %s partition %d on success!", topic, partition)
}

