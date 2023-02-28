package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Fanasity/text2voice/model"
	"github.com/Fanasity/text2voice/queue"
	"github.com/Fanasity/text2voice/storage"
)

// Speech struct
type Speech struct {
	ctx  context.Context
	Task model.Task
}

func NewSpeech(ctx context.Context, task model.Task) *Speech {
	return &Speech{
		ctx:  ctx,
		Task: task,
	}
}
func (speech *Speech) Handle() (err error) {
	config, ok := speech.ctx.Value(model.CONFIG).(model.Config)
	if !ok {
		panic("can not get config")
	}
	timeLoc, _ := time.LoadLocation("Asia/Shanghai")
	var res = model.TaskResult{
		NodeID:          config.System.NodeId,
		ID:              speech.Task.ID,
		CalType:         config.System.Type,
		Status:          0,
		HandleBeginDate: time.Now().In(timeLoc).Format("2006-01-02 15:04:05"),
	}
	defer func() {
		res.HandleCompleteDate = time.Now().In(timeLoc).Format("2006-01-02 15:04:05")
		err = queue.Produce(config.Kafka, config.Kafka.ResultTopic, res)
	}()
	var buffer = new(bytes.Buffer)
	gateway := "https://dict.youdao.com/dictvoice?audio=%s&le=%s&product=pc&type=null"
	// gateway := "http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s"
	url := fmt.Sprintf(gateway, url.QueryEscape(speech.Task.FileInput), "zh")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	response, err := httpClient.Get(url)
	if err != nil {
		res.Status = 1
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()
	n, err := io.Copy(buffer, response.Body)
	if err != err {
		res.Status = 1
		fmt.Println(err)
		return err
	}
	client, err := storage.NewClient(config.MinIO)
	if err != nil {
		res.Status = 1
		fmt.Println(err)
		return err
	}
	res.HandleResult, err = storage.UploadObject(speech.ctx, client, "demo", generateObjectName(speech.Task.FileInput, "mp3"), "audio/mp3", buffer, n)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func generateObjectName(name string, fileType string) string {
	hash := md5.Sum([]byte(name))
	return fmt.Sprintf("voice/%s.%s", hex.EncodeToString(hash[:]), fileType)
}
