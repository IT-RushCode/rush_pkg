package services

import (
	"context"
	"fmt"

	dto "github.com/IT-RushCode/rush_pkg/dto/payment"
	"github.com/IT-RushCode/rush_pkg/models"
	"github.com/IT-RushCode/rush_pkg/repositories"

	yookassa "github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type PaymentService struct {
	repo *repositories.Repositories
}

func NewPaymentService(repo *repositories.Repositories) *PaymentService {
	return &PaymentService{repo: repo}
}

// // PaymentRequest структура для создания платежа
// type PaymentRequest struct {
// 	Amount struct {
// 		Value    string `json:"value"`
// 		Currency string `json:"currency"`
// 	} `json:"amount"`
// 	PaymentMethodData struct {
// 		Type string `json:"type"`
// 	} `json:"payment_method_data"`
// 	Confirmation struct {
// 		Type      string `json:"type"`
// 		ReturnURL string `json:"return_url"`
// 	} `json:"confirmation"`
// 	Description string `json:"description"`
// 	Capture     bool   `json:"capture"`
// }

// // PaymentResponse структура для получения ответа от YooKassa
// type PaymentResponse struct {
// 	ID     string `json:"id"`
// 	Status string `json:"status"`
// 	Amount struct {
// 		Value    string `json:"value"`
// 		Currency string `json:"currency"`
// 	} `json:"amount"`
// 	Confirmation struct {
// 		ConfirmationURL string `json:"confirmation_url"`
// 	} `json:"confirmation"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// // CreatePayment создает новый платеж через YooKassa API
// func (s *PaymentService) CreatePayment(ctx context.Context, dto *dto.PaymentRequest) (*PaymentResponse, error) {
// 	store := &models.YooKassaSetting{}
// 	if err := s.repo.YooKassaSetting.Filter(
// 		ctx,
// 		map[string]interface{}{"point_id": dto.PointID},
// 		store,
// 	); err != nil {
// 		return nil, fmt.Errorf("настройки YooKassa не найдены")
// 	}

// 	// Создание тела запроса
// 	paymentReq := &PaymentRequest{
// 		Description: dto.Description,
// 	}
// 	paymentReq.Capture = true
// 	paymentReq.Amount.Value = dto.Amount
// 	paymentReq.Amount.Currency = dto.Currency
// 	paymentReq.PaymentMethodData.Type = dto.PaymentMethod
// 	paymentReq.Confirmation.Type = "redirect"
// 	paymentReq.Confirmation.ReturnURL = dto.ReturnURL

// 	// Сериализация тела запроса в JSON
// 	reqBody, err := json.Marshal(paymentReq)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при сериализации тела запроса: %v", err)
// 	}

// 	// Создание HTTP-запроса
// 	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при создании запроса: %v", err)
// 	}

// 	// Установка заголовков
// 	req.SetBasicAuth(store.StoreID, store.SecretKey)
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Idempotence-Key", generateIdempotenceKey())

// 	// Выполнение HTTP-запроса
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Чтение тела ответа для вывода при ошибке
// 	bodyBytes, _ := io.ReadAll(resp.Body)
// 	bodyString := string(bodyBytes)

// 	fmt.Println(bodyString)

// 	// Проверка на успешный статус ответа
// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
// 		return nil, fmt.Errorf("неуспешный статус ответа: %s, тело ответа: %s", resp.Status, bodyString)
// 	}

// 	// Десериализация ответа
// 	var paymentResp PaymentResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
// 		return nil, fmt.Errorf("ошибка при десериализации ответа: %v", err)
// 	}

// 	return &paymentResp, nil
// }

// // generateIdempotenceKey генерирует уникальный ключ для обеспечения идемпотентности запросов
// func generateIdempotenceKey() string {
// 	return fmt.Sprintf("%d", time.Now().UnixNano())
// }

func (s *PaymentService) CreatePayment(ctx context.Context, req *dto.PaymentRequest) (*yoopayment.Payment, error) {
	store := &models.YooKassaSetting{}

	if err := s.repo.YooKassaSetting.Filter(
		ctx,
		map[string]interface{}{"point_id": req.PointID},
		store,
	); err != nil {
		return nil, fmt.Errorf("настройки YooKassa не найдены")
	}

	client := yookassa.NewClient(store.StoreID, store.SecretKey)
	paymentKassa := yookassa.NewPaymentHandler(client)

	var payment *yoopayment.Payment
	var err error

	switch yoopayment.PaymentMethodType(req.PaymentMethod) {
	case yoopayment.PaymentTypeBankCard:
		payment, err = paymentKassa.CreatePayment(&yoopayment.Payment{
			Metadata: map[string]interface{}{
				"orderID": 1,
			},
			MerchantCustomerID: "felixkot00@gmail.com",
			Capture:            true,
			Amount: &yoocommon.Amount{
				Value:    req.Amount,
				Currency: req.Currency,
			},
			PaymentMethod: yoopayment.BankCard{
				Card: yoopayment.Card{
					First6:        "220000",
					Last4:         "0004",
					ExpiryYear:    "05",
					ExpiryMonth:   "2030",
					CardType:      "MIR",
					IssuerCountry: "RU",
					IssuerName:    "Sberbank",
					Source:        "sber_pay",
				},
			},
			Description: req.Description,
		})
	case yoopayment.PaymentTypeSBP:
		payment, err = paymentKassa.CreatePayment(&yoopayment.Payment{
			Amount: &yoocommon.Amount{
				Value:    req.Amount,
				Currency: req.Currency,
			},
			PaymentMethod: yoopayment.PaymentMethodType(req.PaymentMethod),
			Description:   req.Description,
			Confirmation: yoopayment.Redirect{
				Type:      "redirect",
				ReturnURL: "https://westerdam.rushcode.ru",
			},
		})
	case yoopayment.PaymentTypeCash:
		return &yoopayment.Payment{}, nil
	default:
		return nil, fmt.Errorf("не поддерживаемый метод оплаты: %s", req.PaymentMethod)
	}

	if err != nil {
		return nil, err
	}

	return payment, nil
}
