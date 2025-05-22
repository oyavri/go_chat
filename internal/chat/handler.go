package chat

import (
	"log/slog"
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
	var req CreateChatRequest

	if err := ctx.BindUri(&req); err != nil {
		slog.Error("[ChatHandler-CreateChatHandler]", "Error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameter"})
		return
	}

	// Parse JSON body
	if err := ctx.BindJSON(&req); err != nil {
		slog.Error("[ChatHandler-CreateChatHandler]", "Error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.Members) < 1 {
		slog.Error("[ChatHandler-CreateChatHandler]", "Error", "A list of members must be provided to create chat")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "A list of members must be provided to create chat"})
		return
	}

	chat, err := h.service.CreateChat(ctx.Request.Context(), req)

	if err != nil {
		slog.Error("[ChatHandler-CreateChatHandler]", "Error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	// Needs a proper JSON response
	ctx.JSON(http.StatusOK, chat)
}

// POST /chats/:chat_id/messages
func (h *ChatHandler) SendMessageHandler(ctx *gin.Context) {
	var req SendMessageRequest

	if err := ctx.BindUri(&req); err != nil {
		slog.Error("[ChatHandler-SendMessageHandler]", "Error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameter"})
		return
	}

	// Parse JSON body
	if err := ctx.BindJSON(&req); err != nil {
		slog.Error("[ChatHandler-SendMessageHandler]", "Error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	message, err := h.service.SendMessage(ctx.Request.Context(), req)

	if err != nil {
		slog.Error("[ChatHandler-SendMessageHandler]", "Error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	ctx.JSON(http.StatusOK, message)
}

// GET /chats/:chat_id/messages
func (h *ChatHandler) GetMessagesHandler(ctx *gin.Context) {
	var req GetMessagesRequest

	if err := ctx.BindUri(&req); err != nil {
		slog.Error("[ChatHandler-GetMessagesHandler]", "Error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path parameter"})
		return
	}

	messages, err := h.service.GetMessages(ctx.Request.Context(), req.ChatId, req.MessageCount, req.Offset)

	if err != nil {
		slog.Error("[ChatHandler-GetMessagesHandler]", "Error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
