package request

import (
	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/validator"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateInvoiceDTO struct {
	CustomerID uuid.UUID
	Amount     float32
	Status     string
	Date       time.Time
}

func NewCreateInvoiceRequest(c *fiber.Ctx) (*CreateInvoiceDTO, error) {
	rules := map[string]any{
		"customer_id": "required,uuid4",
		"amount":      "required,numeric",
		"status":      "required,oneof=paid pending",
		"date":        "required,datetime=2006-01-02T15:04:05Z07:00",
	}

	payload := make(map[string]any)

	// var request struct {
	// 	CustomerID interface{} `json:"customer_id" validate:"required,uuid4"`
	// 	Amount     interface{} `json:"amount" validate:"required,number"`
	// 	Status     interface{} `json:"status" validate:"required,oneof=paid pending"`
	// 	Date       interface{} `json:"date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	// }

	if err := c.BodyParser(&payload); err != nil {
		return nil, err
	}

	// trans, _ := config.Uni.GetTranslator(c.Get("Accept-Language"))

	if err := validator.ValidateMap(payload, rules); err != nil {
		// return nil, err
		// return validation.ValidationErrors
		// for _, fieldErrs := range errs {
		// 	// if fieldErr, ok := fieldErrs.(validation.FieldError); ok {
		// 	// 	log.Printf("%T -- %v -- %v", fieldErr, fieldErr.Translate(trans), fieldErr.Namespace())
		// 	// }
		// 	for _, err := range fieldErrs.(validation.ValidationErrors) {
		// 		log.Printf("%T -- %v -- %v", err, err.Translate(trans), err.Namespace())
		// 		// return nil, err.(validation.FieldError)
		// 	}
		// }

		return nil, err
	}

	dto := CreateInvoiceDTO{}

	customerId, err := uuid.Parse(payload["customer_id"].(string))
	if err != nil {
		return nil, err
	}

	invoiceDate, err := time.Parse("2006-01-02T15:04:05Z", payload["date"].(string))
	if err != nil {
		return nil, err
	}

	var amount float64

	switch payloadAmount := payload["amount"].(type) {
	case string:
		amount, err = strconv.ParseFloat(payloadAmount, 64)
		if err != nil {
			return nil, err
		}
	case float64:
		amount = payloadAmount
	}

	dto.CustomerID = customerId
	dto.Amount = float32(amount)
	dto.Status = payload["status"].(string)
	dto.Date = invoiceDate

	return &dto, nil
}

type UpdateInvoiceDTO struct {
	InvoiceID  uuid.UUID
	CustomerID uuid.UUID
	Amount     float32
	Status     string
}

func NewUpdateInvoiceRequest(c *fiber.Ctx) (*UpdateInvoiceDTO, error) {
	if err := config.Validate.Var(c.Params("id"), "required,uuid4"); err != nil {
		return nil, fiber.NewError(400, "Invalid Invoice ID in path parameter")
	}

	rules := map[string]any{
		"customer_id": "required,uuid4",
		"amount":      "required,numeric",
		"status":      "required,oneof=paid pending",
	}

	payload := make(map[string]any)

	if err := c.BodyParser(&payload); err != nil {
		return nil, err
	}

	if err := validator.ValidateMap(payload, rules); err != nil {
		return nil, err
	}

	invoiceID, _ := uuid.Parse(c.Params("id"))
	customerID, _ := uuid.Parse(payload["customer_id"].(string))

	var amount float64
	var err error

	switch payloadAmount := payload["amount"].(type) {
	case string:
		amount, err = strconv.ParseFloat(payloadAmount, 64)
		if err != nil {
			return nil, err
		}
	case float64:
		amount = payloadAmount
	}

	return &UpdateInvoiceDTO{
		InvoiceID:  invoiceID,
		CustomerID: customerID,
		Amount:     float32(amount),
		Status:     payload["status"].(string),
	}, nil
}

type GetInvoicePagesDTO struct {
	Search string
	Limit  int
}

func NewGetInvoicePagesRequest(c *fiber.Ctx) (*GetInvoicePagesDTO, error) {
	rules := map[string]any{
		"limit": "numeric",
	}

	payload := map[string]string{
		"search": c.Query("search"),
		"limit":  c.Query("limit", "6"),
	}

	if err := validator.ValidateMapQueries(payload, rules); err != nil {
		return nil, err
	}

	search := payload["search"]
	limit, _ := strconv.Atoi(payload["limit"])

	return &GetInvoicePagesDTO{
		Search: search,
		Limit:  limit,
	}, nil
}

type GetFilteredInvoicesDTO struct {
	Search string
	Page   int
	Limit  int
}

func NewGetFilteredInvoicesRequest(c *fiber.Ctx) (*GetFilteredInvoicesDTO, error) {
	rules := map[string]any{
		"page":  "numeric",
		"limit": "numeric",
	}

	payload := map[string]string{
		"search": c.Query("search"),
		"page":   c.Query("page", "1"),
		"limit":  c.Query("limit", "6"),
	}

	if err := validator.ValidateMapQueries(payload, rules); err != nil {
		return nil, err
	}

	search := payload["search"]
	page, _ := strconv.Atoi(payload["page"])
	limit, _ := strconv.Atoi(payload["limit"])

	return &GetFilteredInvoicesDTO{
		Search: search,
		Page:   page,
		Limit:  limit,
	}, nil
}
