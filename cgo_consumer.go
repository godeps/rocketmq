// +build linux darwin,amd64,cgo
// +build rocket_cgo

package rocketmq

import (
	"context"
	cmq "github.com/apache/rocketmq-client-go/core"
	mqc "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"log"
	"sync"
)

type CRocketMQConsumer struct {
	md         *Metadata
	client     cmq.PushConsumer
	contextMap sync.Map
	config     *cmq.PushConsumerConfig
}

func init() {
	Consumers[TCPCGOProto] = &CRocketMQConsumer{}
}

// NewRocketMQConsumer
func NewRocketMQConsumer(md *Metadata) (*CRocketMQConsumer, error) {
	mq := &CRocketMQConsumer{md: md}
	return mq, mq.Init(md)
}

func (mq *CRocketMQConsumer) Init(md *Metadata) error {
	mq.md = md
	if md.ConsumerBatchSize < 1 {
		mq.md.ConsumerBatchSize = defaultConsumerNumOfMessages
	}

	mq.config = &cmq.PushConsumerConfig{
		ClientConfig: cmq.ClientConfig{
			GroupID:      md.ConsumerGroup,
			NameServer:   md.NameServer,
			InstanceName: md.Namespace,
			Credentials: &cmq.SessionCredentials{
				AccessKey: md.AccessKey,
				SecretKey: md.AccessSecret,
				Channel:   defaultRocketMQPublicChannel,
			},
		},
		ThreadCount:         md.ConsumerThreadNums,
		MessageBatchMaxSize: md.ConsumerBatchSize,
		Model:               cmq.Clustering,
		ConsumerModel:       cmq.CoCurrently,
	}
	if md.NameServer == defaultRocketMQChannel || (len(md.NameServer) == 0 && len(md.Endpoint) == 0) {
		mq.config.Credentials.Channel = defaultRocketMQChannel
	}
	var err error
	if mq.client, err = cmq.NewPushConsumer(mq.config); err != nil {
		return err
	}
	return mq.client.Start()
}

// Start the PullConsumer for consuming message
func (mq *CRocketMQConsumer) Start() error {
	return nil
}

// Shutdown the PullConsumer
func (mq *CRocketMQConsumer) Shutdown() error {
	mq.contextMap.Range(func(key, value interface{}) bool {
		if err := mq.Unsubscribe(key.(string)); err != nil {
			log.Println(err)
		}
		return true
	})
	return nil
}

// Subscribe a topic for consuming
func (mq *CRocketMQConsumer) Subscribe(topic string, selector mqc.MessageSelector, f func(context.Context, ...*primitive.MessageExt) (mqc.ConsumeResult, error)) error {
	consumer, err := cmq.NewPushConsumer(mq.config)
	if err != nil {
		return err
	}

	if err := consumer.Subscribe(topic, selector.Expression, func(msgEntry *cmq.MessageExt) cmq.ConsumeStatus {
		msg := &primitive.MessageExt{}
		msg.MsgId = msgEntry.MessageID
		msg.BornTimestamp = msgEntry.BornTimestamp
		msg.Topic = topic
		msg.Body = []byte(msgEntry.Body)
		msg.WithProperties(msgEntry.Property)
		status, err := f(context.TODO(), msg)
		if err == nil && status == mqc.ConsumeSuccess {
			return cmq.ConsumeSuccess
		}
		if err != nil {
			log.Printf("consume message failed. topic:%s MessageID:%s status:%v", topic, msgEntry.MessageID, status)
		}
		return cmq.ReConsumeLater
	}); err != nil {
		return err
	}

	if err = consumer.Start(); err != nil {
		return err
	}
	mq.contextMap.Store(topic, consumer)
	return nil
}

// Unsubscribe a topic
func (mq *CRocketMQConsumer) Unsubscribe(topic string) error {
	v, ok := mq.contextMap.LoadAndDelete(topic)
	if ok {
		mqConsumer := v.(cmq.PushConsumer)
		return mqConsumer.Shutdown()
	}
	return nil
}
