package processor

import (
	"database/sql"
	"encoding/json"
	"log"
	"wb-project1/database"
	"wb-project1/model"
)

type OrderHandler struct {
	DB    *sql.DB
	Cache map[string]string
}

func (h *OrderHandler) HandleKafkaMessage(msg []byte) {
	var order model.Order
	if err := json.Unmarshal(msg, &order); err != nil {
		log.Println("Невалидный JSON:", err)
		return
	}

	if order.OrderUID == "" {
		log.Println("Пропущен заказ без UID")
		return
	}

	if err := database.SaveOrder(h.DB, order); err != nil {
		log.Println("Ошибка при сохранении заказа:", err)
		return
	}

	jsonStr, err := json.Marshal(order)
	if err != nil {
		log.Println("Ошибка маршалинга:", err)
		return
	}

	h.Cache[order.OrderUID] = string(jsonStr)
	log.Printf("Заказ %s обработан\n", order.OrderUID)
}
