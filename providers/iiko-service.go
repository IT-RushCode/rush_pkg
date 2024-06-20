package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IikoService struct {
	APIKey   string
	APIURL   string
	Username string
	Password string
}

func NewIikoService(apiKey, apiURL, username, password string) *IikoService {
	return &IikoService{
		APIKey:   apiKey,
		APIURL:   apiURL,
		Username: username,
		Password: password,
	}
}

func (s *IikoService) SendOrder(orderData interface{}) (interface{}, error) {
	jsonData, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/orders", s.APIURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *IikoService) GetOrderStatus(orderID string) (interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/orders/%s", s.APIURL, orderID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
