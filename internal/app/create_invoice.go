package app

import (
	"context"
	"fmt"

	"go-nextjs-dashboard/internal/customer"
	"go-nextjs-dashboard/internal/db"
	"go-nextjs-dashboard/internal/invoice"
	"go-nextjs-dashboard/internal/logger"
)

type CreateInvoice struct {
	custStore customer.Store
	invStore  invoice.Store
	txm       db.TxManager
	logger    logger.Logger
}

func NewCreateInvoice(custStore customer.Store, invStore invoice.Store, txm db.TxManager, logger logger.Logger) *CreateInvoice {
	return &CreateInvoice{
		custStore: custStore,
		invStore:  invStore,
		txm:       txm,
		logger:    logger,
	}
}

func (c *CreateInvoice) Execute(ctx context.Context, i invoice.Invoice) (inv *invoice.Invoice, err error) {
	err = c.txm.Do(ctx, func(txCtx context.Context) error {
		exists, txErr := c.custStore.Exists(txCtx, i.CustomerID)
		if txErr != nil {
			return fmt.Errorf("exists customer: %w", txErr)
		}

		if !exists {
			return customer.ErrCustomerNotFound
		}

		inv, txErr = c.invStore.Save(txCtx, i)
		if txErr != nil {
			return fmt.Errorf("save invoice: %w", txErr)
		}

		return nil
	})
	return
}
