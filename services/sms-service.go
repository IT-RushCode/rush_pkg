package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/IT-RushCode/rush_pkg/config"
	dto "github.com/IT-RushCode/rush_pkg/dto/sms"
	"github.com/IT-RushCode/rush_pkg/utils"
)

// Промежуточная структура для получения статуса
type intermediateResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

type SmsService struct {
	cfg *config.SmsConfig
}

func NewSmsService(cfg *config.SmsConfig) *SmsService {
	return &SmsService{cfg: cfg}
}

// SendSMS делает POST-запрос к указанному URL с предоставленным payload.
func (s SmsService) SendSMS(cfg *config.Config, data dto.SMSRequestDTO) (*dto.SmsSenderResponse, error) {
	url := cfg.SMS.URL
	method := cfg.SMS.METHOD

	var (
		phone   = strings.TrimSpace(data.Messages[0].Phone)
		otpCode string
	)

	data.Messages[0].Sender = cfg.SMS.SENDER

	// Если IsOTP true, то отправляем OTP код на номер тел.
	if data.IsOTP {
		otpCode = utils.GenerateOTP()
		data.Messages[0].Text = fmt.Sprintf(cfg.SMS.TEMPLATE, otpCode)
	}

	// Если Channel пустой, то для всех сообщений устанавливаем буквенный канал сообщений
	for i := range data.Messages {
		if data.Messages[i].Channel == "" {
			data.Messages[i].Channel = "char"
		}

		if data.Messages[i].Channel == "char" {
			data.Messages[0].Sender = cfg.SMS.SENDER
		}
	}

	payload := struct {
		APIKey string           `json:"apiKey"`
		SMS    []dto.SMSMessage `json:"sms"`
	}{
		APIKey: cfg.SMS.TOKEN,
		SMS:    data.Messages,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации payload: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}
	log.Println(string(body))

	// Промежуточная десериализация для получения статуса
	var intermRes intermediateResponse
	err = json.Unmarshal(body, &intermRes)
	if err != nil {
		return nil, fmt.Errorf("ошибка при десериализации промежуточного ответа: %v", err)
	}

	// В зависимости от статуса десериализуем data
	var jsonRes dto.SMSRes
	if intermRes.Status == "success" {
		var dataArray []dto.Data
		err = json.Unmarshal(intermRes.Data, &dataArray)
		if err != nil {
			return nil, fmt.Errorf("ошибка при десериализации data как массива: %v", err)
		}
		jsonRes = dto.SMSRes{
			Status: intermRes.Status,
			Data:   dataArray,
		}
	} else {
		var dataObject dto.Data
		err = json.Unmarshal(intermRes.Data, &dataObject)
		if err != nil {
			return nil, fmt.Errorf("ошибка при десериализации data как объекта: %v", err)
		}
		jsonRes = dto.SMSRes{
			Status: intermRes.Status,
			Data:   []dto.Data{dataObject},
		}
	}

	return &dto.SmsSenderResponse{
		Message: jsonRes,
		Phone:   phone,
		OTPCode: otpCode,
	}, nil
}
