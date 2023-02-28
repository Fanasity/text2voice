package main

import (
	"context"
	"encoding/json"
	"fmt"

	"aiServer/model"
	"aiServer/queue"
	"aiServer/service"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("conf/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var config model.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, model.CONFIG, config)
	consumers, err := queue.Consumer(config.Kafka)
	if err != nil || len(consumers) == 0 {
		panic(fmt.Errorf("get consumers err: %w", err))
	}

	for _, pc := range consumers {
		go queue.ConsumerMeesage(ctx, pc, func(ctx context.Context, msg []byte, key []byte) (queue.Task, error) {
			var task model.Task
			err := json.Unmarshal(msg, &task)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			switch config.System.Type {
			case model.TaskTypePDFToWord:
			case model.TaskTypeTextToVoice:
				return service.NewSpeech(ctx, task), nil
			case model.TaskTypePictureCluster:
			default:

			}
			return nil, nil
		})
	}

	select {}

}
