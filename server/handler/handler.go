package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/trancecho/mundo-prd-manager/initialize"
	"github.com/trancecho/mundo-prd-manager/models"
	"github.com/trancecho/mundo-prd-manager/server"
	"github.com/trancecho/mundo-prd-manager/server/api"
	"github.com/trancecho/mundo-prd-manager/server/libx"
	"github.com/trancecho/ragnarok/fastgpt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetPersonalChatID(c *gin.Context) {
	username := libx.GetUsername(c)
	var charHistory models.ChatIDHistory
	result := initialize.GetDB().Where("username = ?", username).First(&charHistory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			libx.Ok(c, http.StatusOK, "没有ChatID记录", nil)
			return
		}
		libx.Err(c, 500, "查询ChatID记录失败", result.Error)
		return
	}
	libx.Ok(c, http.StatusOK, "查询ChatID记录成功", charHistory.ChatIDs)
}

func ProductGenerate(c *gin.Context) {
	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		libx.Err(c, 400, "请求参数错误", err)
		return
	}

	uid, username := libx.Uid(c), libx.GetUsername(c)
	completeChatID := server.GenerateChatID(req.ChatID, strconv.Itoa(int(uid)))

	//检查用户是否使用过该ID
	var charHistory models.ChatIDHistory
	result := initialize.GetDB().Where("username = ?", username).First(&charHistory)
	if result.Error != nil {
		// 记录不存在，创建ID
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			charHistory = models.ChatIDHistory{
				Username: username,
				ChatIDs:  []string{req.ChatID},
			}
			if err := initialize.GetDB().Create(&charHistory).Error; err != nil {
				libx.Err(c, 500, "创建ChatID记录失败", err)
				return
			}
		} else {
			libx.Err(c, 500, "查询ChatID记录失败", result.Error)
			return
		}
	} else {
		//用户记录已经存在，检查是否有该chatID
		hastChatID := false
		for _, id := range charHistory.ChatIDs {
			if id == req.ChatID {
				hastChatID = true
				break
			}
		}
		// 如果没有该ChatID，则添加
		if !hastChatID {
			charHistory.ChatIDs = append(charHistory.ChatIDs, req.ChatID)
			if err := initialize.GetDB().Save(&charHistory).Error; err != nil {
				libx.Err(c, 500, "更新ChatID记录失败", err)
				return
			}
		}
	}

	// gpt_request
	var gptReq fastgpt.Request
	gptReq.ChatID = completeChatID
	gptReq.Stream = req.Stream
	gptReq.Detail = req.Detail
	gptReq.Variables = req.Variables
	m2 := fastgpt.Message{Role: "user", Content: req.Description}
	gptReq.Messages = []fastgpt.Message{m2}

	// 调用飞书接口
	resp, err := api.Cli.Fcli.FastGPTChat(gptReq)
	if err != nil {
		libx.Err(c, 500, "gpt接口请求失败", err)
		return
	}
	// 解析content
	var gptresp models.GptResponse
	if err = json.Unmarshal(resp, &gptresp); err != nil {
		libx.Err(c, 500, "解析GPT响应失败", err)
		return
	}
	if len(gptresp.Choices) > 0 {
		content := gptresp.Choices[0].Message.Content
		libx.Ok(c, 200, content, nil)
		//response, e := api.Cli.Client.Docx.Document.Create(c, larkdocx.NewCreateDocumentReqBuilder().
		//	Body(larkdocx.NewCreateDocumentReqBodyBuilder().
		//		FolderToken(req.Folder). // 文件夹 token，传空表示在根目录创建文档
		//		Title(req.DocuTitle).    // 文档标题
		//		Build()).
		//	Build())
		//if e != nil {
		//	libx.Err(c, 500, "创建文档失败", e)
		//	return
		//}
		//docID := response.Data.Document.DocumentId
		//创建文档块
	} else {
		libx.Err(c, 500, "GPT响应中没有内容", nil)
	}
}

func DeleteChatID(c *gin.Context) {
	username := libx.GetUsername(c)
	var charHistory models.ChatIDHistory
	if err := initialize.GetDB().Where("username = ?", username).First(&charHistory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			libx.Err(c, 400, "没有ChatID记录", err)
			return
		}
		libx.Err(c, 500, "查询ChatID记录失败", err)
		return
	}
	chatID := c.Query("chat_id")
	if chatID == "" {
		libx.Err(c, 400, "缺少chat_id参数", nil)
		return
	}
	// 检查ChatID是否存在
	found := false
	for i, id := range charHistory.ChatIDs {
		if id == chatID {
			// 删除ChatID
			charHistory.ChatIDs = append(charHistory.ChatIDs[:i], charHistory.ChatIDs[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		libx.Err(c, 400, "ChatID不存在", nil)
		return
	}
	// 更新数据库
	if err := initialize.GetDB().Save(&charHistory).Error; err != nil {
		libx.Err(c, 500, "更新ChatID记录失败", err)
		return
	}
	libx.Ok(c, http.StatusOK, "ChatID删除成功", nil)
}
