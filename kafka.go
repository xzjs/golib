package lib

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

//Producer .
func Producer(topic string, value string) (err error) {
	//配置
	conf := Conf()
	connStr := fmt.Sprintf("%s:%s", conf.GetConf("kafka", "host"), conf.GetConf("kafka", "port"))
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	p, err := sarama.NewSyncProducer([]string{connStr}, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	//defer p.Close()

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: 1,
		Value:     sarama.ByteEncoder(value),
	}
	_, _, err = p.SendMessage(msg)
	if err != nil {
		log.Printf("send message(%s) err=%s \n", value, err)
	}
	return
}
