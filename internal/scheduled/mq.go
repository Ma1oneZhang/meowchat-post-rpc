package scheduled

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/xh-polaris/meowchat-post-rpc/internal/config"
	"log"
	"net/url"
	"sync"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	prod *rocketmq.Producer
	pmu  sync.Mutex
)

func checkMqInstance(c config.Config) {
	if prod == nil {
		pmu.Lock()
		if prod == nil {
			produce, err := rocketmq.NewProducer(
				producer.WithNsResolver(primitive.NewPassthroughResolver(c.RocketMq.URL)),
				producer.WithRetry(c.RocketMq.Retry),
				producer.WithGroupName(c.RocketMq.GroupName),
			)
			if err != nil {
				log.Fatal(err)
			}
			err = produce.Start()
			if err != nil {
				log.Fatal(err)
			}
			prod = &produce
		}
		pmu.Unlock()
	}
}

// SendUrlUsedMessageToSts we need delivery it immediately
func SendUrlUsedMessageToSts(c *config.Config, msg *[]string) {
	var urls []url.URL
	for _, v := range *msg {
		u, _ := url.Parse(v)
		urls = append(urls, *u)
	}
	messageToMq(c, urls, "sts_used_url", 0)
}

// messageToMq set message delay time to consume.
// reference delay level definition: 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
// delay level starts from 1. for example, if we set param level=1, then the delay time is 1s.
// When delay level equal to 0, the message delivery immediately
func messageToMq(c *config.Config, message interface{}, topic string, level int) {
	checkMqInstance(*c)
	json, err := jsonx.Marshal(message)
	if err != nil {
		logx.Alert(err.Error())
		return
	}
	msg := &primitive.Message{
		Topic: topic,
		Body:  json,
	}
	msg.WithDelayTimeLevel(level)
	res, err := (*prod).SendSync(context.Background(), msg)
	if err != nil || res.Status != primitive.SendOK {
		for i := 0; i < 2; i++ {
			res, err := (*prod).SendSync(context.Background(), msg)
			if err == nil && res.Status == primitive.SendOK {
				break
			}
		}
	}
}
