package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type ApiReqBody struct {
	Model string          `json:"model"`
	Input ApiReqBodyInput `json:"input"`
}

type ApiReqBodyInput struct {
	Prompt string `json:"prompt"`
}

type ApiRes struct {
	Output struct {
		FinishReason string `json:"finish_reason"`
		Text         string `json:"text"`
	} `json:"output"`
	RequestID string `json:"request_id"`
}

func aiReq(client *resty.Client, prompt *string) string {
	log.Info("请求大模型中......")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken("Bearer " + getAiAkConf()).
		SetBody(ApiReqBody{
			Model: getAiModelConf(),
			Input: ApiReqBodyInput{
				Prompt: *prompt,
			},
		}).
		Post(getAiUrlConf())

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode() != 200 {
		log.Fatal(resp)
	}

	log.Info("请求完成......")

	result := ApiRes{}
	jsonErr := json.Unmarshal(resp.Body(), &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result.Output.Text
}
