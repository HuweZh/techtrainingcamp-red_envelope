package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"red_envelope/configure"
	"red_envelope/models"
	"time"
)

var MyProducer rocketmq.Producer
var MyConsumer rocketmq.PushConsumer

// InitRMQ (如果是init函数就会自动执行)
func InitRMQ() {
	//初始化生产者
	var err error
	MyProducer, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{configure.RmqNameseverAddr})),
		producer.WithRetry(2),
	)
	if err != nil {
		if configure.GIN_MODE == "debug" {
			fmt.Printf("create producer error: %s\n", err.Error())
		}
		MyLog.Error("create producer error: ", err.Error())
		os.Exit(1)
	}
	err = MyProducer.Start()
	if err != nil {
		if configure.GIN_MODE == "debug" {
			fmt.Printf("start producer error: %s\n", err.Error())
		}
		MyLog.Error("start producer error: ", err.Error())
		os.Exit(1)
	}

	//初始化topic
	rmqWalletTopic := "wallet"
	myAdmin, _ := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{configure.RmqNameseverAddr})))

	err = myAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate(rmqWalletTopic),
		admin.WithBrokerAddrCreate(configure.RmqBrokerAddr),
	)
	if err != nil {
		if configure.GIN_MODE == "debug" {
			fmt.Printf("create topic error: %s\n", err.Error())
		}
		MyLog.Error("create topic error: ", err.Error())
		os.Exit(1)
	}

	//初始化消费者
	MyConsumer, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(configure.GROUP_NAME),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{configure.RmqNameseverAddr})),
	)
	err = MyConsumer.Subscribe(rmqWalletTopic, consumer.MessageSelector{}, callbackWallet)
	err = MyConsumer.Start()
	if err != nil {
		if configure.GIN_MODE == "debug" {
			fmt.Printf("start consumer error: %s\n", err.Error())
		}
		MyLog.Error("start consumer error: ", err.Error())
		os.Exit(1)
	}
}

type RocketMqMessage struct {
	MessageId    int64
	MessageBytes []byte
	Tag          string
	Topic        string
}

func SendToRMQ(msg RocketMqMessage) {
	ms := primitive.NewMessage(msg.Topic, msg.MessageBytes)
	_, err := MyProducer.SendSync(context.Background(), ms)
	//err := MyProducer.SendAsync(context.Background(),
	//	func(ctx context.Context, result *primitive.SendResult, e error) {
	//		if e != nil {
	//			if configure.GIN_MODE == "debug" {
	//				fmt.Printf("receive message error: %s\n", e)
	//			}
	//			MyLog.Error("receive message error: ", e)
	//		}
	//	}, primitive.NewMessage(msg.Topic, msg.MessageBytes))
	if err != nil {
		if configure.GIN_MODE == "debug" {
			fmt.Printf("send message error: %s\n", err)
		}
		MyLog.Error("send message error: ", err)
	}
}

func callbackWallet(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		time.Sleep(time.Duration(1) * time.Microsecond)//每1ms执行一次，这是个异步过程，不影响外面的request
		wallet := models.Wallet{}
		err := json.Unmarshal(msgs[i].Body, &wallet)
		if err != nil {
			MyLog.Error("callbackWallet error: ", err)
			continue
		}
		err = models.InsertWallet(DB, &wallet)
		if err != nil { //可能存在该记录，则再试试更新
			if err = models.UpdateWalletByUid(DB, wallet.Uid, &map[string]interface{}{"money": wallet.Money}); err != nil {
				MyLog.Error("update wallet err:", err)
			}
		}
	}
	return consumer.ConsumeSuccess, nil
}
