package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	dto "github.com/IT-RushCode/rush_pkg/dto/chat"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketHandler struct {
	repo    *repositories.Repositories
	clients map[string]*websocket.Conn // Хранение подключений клиентов
}

func NewWebSocketHandler(repo *repositories.Repositories) *WebSocketHandler {
	return &WebSocketHandler{
		repo:    repo,
		clients: make(map[string]*websocket.Conn),
	}
}

// WebSocketChat для клиентов
func (h *WebSocketHandler) WebSocketChat() func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		sessionID, _ := strconv.Atoi(c.Params("sessionID"))
		clientID, _ := strconv.Atoi(c.Params("clientID"))
		ctx := context.Background()

		clientConnID := fmt.Sprintf("client-%d", clientID)
		h.clients[clientConnID] = c

		defer func() {
			delete(h.clients, clientConnID)
			c.Close()
		}()

		history, err := h.repo.Chat.GetMessages(
			ctx,
			&dto.GetChatHistoryDTO{
				SessionID: uint(sessionID),
				Limit:     50,
				Offset:    0,
			},
		)
		if err != nil {
			log.Println("Ошибка загрузки истории чата:", err)
			return
		}

		// Отправляем историю сообщений клиенту
		for _, msg := range history {
			messageData := map[string]interface{}{
				"id":        msg.ID,
				"senderId":  msg.SenderID,
				"isSupport": msg.IsSupport,
				"body":      msg.Body,
				"timestamp": msg.Timestamp,
			}

			messageJSON, err := json.Marshal(messageData)
			if err != nil {
				log.Println("Ошибка преобразования сообщения в JSON:", err)
				return
			}

			if err := c.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				log.Println("Ошибка отправки истории чата:", err)
				return
			}
		}

		// Основной цикл получения новых сообщений
		for {
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Ошибка чтения сообщения:", err)
				break
			}

			// Создание нового сообщения
			message := &models.ChatMessage{
				SessionID: uint(sessionID),
				SenderID:  uint(clientID),
				IsSupport: false, // Сообщение от клиента
				Body:      string(msg),
				Timestamp: time.Now(),
			}

			// Сохранение сообщения в базе данных
			err = h.repo.Chat.CreateMessage(ctx, message)
			if err != nil {
				log.Println("Ошибка сохранения сообщения:", err)
				break
			}

			// Отправка сообщения поддержке (если подключена)
			for id, client := range h.clients {
				if id == fmt.Sprintf("support-%d", sessionID) { // Отправляем только поддержке
					if err := client.WriteMessage(msgType, msg); err != nil {
						log.Println("Ошибка отправки сообщения поддержке:", err)
					}
				}
			}
		}
	}
}

// WebSocket для "Поддержки"
func (h *WebSocketHandler) WebSocketSupport() func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		sessionID, _ := strconv.Atoi(c.Params("sessionID"))
		ctx := context.Background()

		supportConnID := fmt.Sprintf("support-%d", sessionID)
		h.clients[supportConnID] = c

		defer func() {
			delete(h.clients, supportConnID)
			c.Close()
		}()

		// Основной цикл получения новых сообщений от поддержки
		for {
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Ошибка чтения сообщения:", err)
				break
			}

			// Создание нового сообщения
			message := &models.ChatMessage{
				SessionID: uint(sessionID),
				SenderID:  0,    // Условный ID для "Поддержки"
				IsSupport: true, // Сообщение от поддержки
				Body:      string(msg),
				Timestamp: time.Now(),
			}

			// Сохранение сообщения в базе данных
			err = h.repo.Chat.CreateMessage(ctx, message)
			if err != nil {
				log.Println("Ошибка сохранения сообщения:", err)
				break
			}

			// Отправка сообщения клиенту (если подключен)
			for id, client := range h.clients {
				if id == fmt.Sprintf("client-%d", sessionID) { // Отправляем только клиенту
					if err := client.WriteMessage(msgType, msg); err != nil {
						log.Println("Ошибка отправки сообщения клиенту:", err)
					}
				}
			}
		}
	}
}

// func (h *WebSocketHandler) WebSocketChat() func(c *websocket.Conn) {
// 	return func(c *websocket.Conn) {
// 		sessionID, _ := strconv.Atoi(c.Params("sessionID"))
// 		senderID, _ := strconv.Atoi(c.Params("senderID"))
// 		receiverID, _ := strconv.Atoi(c.Params("receiverID"))
// 		limit, _ := strconv.Atoi(c.Query("limit"))
// 		offset, _ := strconv.Atoi(c.Query("offset"))
// 		ctx := context.Background()

// 		clientID := fmt.Sprintf("%d-%d", senderID, receiverID)
// 		h.clients[clientID] = c
// 		defer func() {
// 			delete(h.clients, clientID)
// 			c.Close() // Закрываем соединение после выхода из функции
// 		}()

// 		history := &[]models.ChatMessage{}

// 		// Загрузка истории сообщений
// 		err := h.ChatController.GetHistoryMessages(
// 			ctx,
// 			&dto.GetChatHistoryDTO{
// 				SessionID:  uint(sessionID),
// 				SenderID:   uint(senderID),
// 				ReceiverID: uint(receiverID),
// 				Limit:      uint(limit),
// 				Offset:     uint(offset),
// 			},
// 			history,
// 		)
// 		if err != nil {
// 			fmt.Println("ошибка загрузки истории чата:", err)
// 			return
// 		}

// 		// Отправка истории сообщений клиенту
// 		for _, msg := range *history {
// 			messageData := map[string]interface{}{
// 				"id":         msg.ID,
// 				"senderId":   msg.SenderID,
// 				"receiverId": msg.ReceiverID,
// 				"sessionId":  msg.SessionID,
// 				"type":       msg.Type,
// 				"body":       msg.Body,
// 				"timestamp":  msg.Timestamp,
// 			}

// 			// Преобразуем данные в JSON
// 			messageJSON, err := json.Marshal(messageData)
// 			if err != nil {
// 				fmt.Println("ошибка преобразования сообщения в JSON:", err)
// 				return
// 			}

// 			// Отправляем JSON сообщение клиенту
// 			if err := c.WriteMessage(msg.Type, messageJSON); err != nil {
// 				fmt.Println("ошибка отправки истории чата:", err)
// 				return
// 			}
// 		}

// 		// Основной цикл обработки входящих сообщений
// 		for {
// 			msgType, msg, err := c.ReadMessage()
// 			if err != nil {
// 				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
// 					fmt.Println("Клиент закрыл соединение:", clientID)
// 				} else {
// 					fmt.Println("Ошибка чтения сообщения:", err)
// 				}
// 				break
// 			}

// 			message := &models.ChatMessage{
// 				SessionID:  sessionID,
// 				SenderID:   uint(senderID),
// 				ReceiverID: uint(receiverID),
// 				Type:       msgType,
// 				Body:       string(msg),
// 				Timestamp:  time.Now(),
// 			}

// 			// Сохранение сообщения в базу данных
// 			err = h.ChatController.CreateMessage(ctx, message)
// 			if err != nil {
// 				fmt.Println("Ошибка сохранения сообщения:", err)
// 				break
// 			}

// 			// Отправка сообщения только клиентам, которые связаны с этой сессией
// 			for id, client := range h.clients {
// 				if id == clientID || id == fmt.Sprintf("%d-%d", receiverID, senderID) {
// 					if err := client.WriteMessage(msgType, msg); err != nil {
// 						fmt.Println("Ошибка отправки сообщения клиенту:", err)
// 					}
// 				}
// 			}
// 		}

// 	}
// }
