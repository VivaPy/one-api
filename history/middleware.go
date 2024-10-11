package history

import (
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/ctxkey"
	"github.com/songquanpeng/one-api/common/logger"
	"github.com/songquanpeng/one-api/relay/adaptor/openai"
	"github.com/songquanpeng/one-api/relay/model"
)

var recorder Recorder

func init() {
	recorder = NewRecorder("rbmq")
}

func SaveHistory(c *gin.Context) {

	var chatMessages []model.Message

	// 通过中间件传递
	reqContent, exists := c.Get("requestBody")
	if !exists {
		logger.Error(c.Request.Context(), "fail to get request body from context!")
		return
	}
	messages := reqContent.([]model.Message)
	chatMessages = append(chatMessages, messages...)

	respContent, exists := c.Get("responseBody")
	if !exists {
		logger.Error(c.Request.Context(), "fail to get response body from context!")
		return
	}

	for _, message := range respContent.([]openai.TextResponseChoice) {
		chatMessages = append(chatMessages, message.Message)
	}

	userID := c.GetInt(ctxkey.Id)
	usageVal, _ := c.Get("usage")
	usage, _ := usageVal.(model.Usage)

	recorder.Push(userID, chatMessages, &usage)
}
