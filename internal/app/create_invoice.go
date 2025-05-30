package app

import (
	"context"
	"fmt"

	"github.com/gelozr/go-dash/internal/customer"
	"github.com/gelozr/go-dash/internal/db"
	"github.com/gelozr/go-dash/internal/invoice"
	"github.com/gelozr/go-dash/internal/logger"
)

type CreateInvoice struct {
	custSvc *customer.Service
	invSvc  *invoice.Service
	txm     db.TxManager
	logger  logger.Logger
}

func NewCreateInvoice(
	custSvc *customer.Service,
	invSvc *invoice.Service,
	txm db.TxManager,
	logger logger.Logger,
) *CreateInvoice {
	return &CreateInvoice{custSvc, invSvc, txm, logger}
}

func (c *CreateInvoice) Execute(ctx context.Context, i invoice.Invoice) (*invoice.Invoice, error) {
	var inv *invoice.Invoice

	txErr := c.txm.Do(ctx, func(txCtx context.Context) error {
		exists, err := c.custSvc.Exists(txCtx, *i.CustomerID)
		if err != nil {
			return fmt.Errorf("exists customer: %w", err)
		}

		if !exists {
			return customer.ErrCustomerNotFound
		}

		inv, err = c.invSvc.Create(txCtx, i)
		if err != nil {
			return fmt.Errorf("create invoice: %w", err)
		}
		return nil
	})

	if txErr != nil {
		return nil, fmt.Errorf("create invoice tx: %w", txErr)
	}

	return inv, nil
}
