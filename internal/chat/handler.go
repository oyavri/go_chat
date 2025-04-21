package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service *ChatService
}

func NewChatHandler(service *ChatService) *ChatHandler {
	return &ChatHandler{
		service: service,
	}
}

func (h *ChatHandler) CreateChatHandler(ctx *gin.Context) {
	chat, err := h.service.CreateChat()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	ctx.JSON(http.StatusOK, chat)
}

// POST /chats/:chat_id/messages
func (h *ChatHandler) SendMessageHandler(ctx *gin.Context) {
	var req SendMessageRequest

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameter"})
		return
	}

	// Parse JSON body
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	message, err := h.service.SendMessage(req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	ctx.JSON(http.StatusOK, message)
}

// GET /chats/:chat_id/messages
func (h *ChatHandler) GetMessagesHandler(ctx *gin.Context) {
	var req GetMessagesRequest

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameter"})
		return
	}

	messages, err := h.service.GetMessages(req.ChatId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
