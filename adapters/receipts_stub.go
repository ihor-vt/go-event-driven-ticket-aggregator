package adapters

import (
	"context"
	"sync"

	"tickets/entities"
)

type ReceiptsServiceStub struct {
	lock sync.Mutex

	IssuedReceipts map[string]entities.IssueReceiptRequest
	VoidedReceipts []entities.VoidReceipt
}

func (c *ReceiptsServiceStub) IssueReceipt(
	ctx context.Context,
	request entities.IssueReceiptRequest,
) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.IssuedReceipts[request.TicketID] = request

	return nil
}

func (c *ReceiptsServiceStub) VoidReceipt(ctx context.Context, request entities.VoidReceipt) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.VoidedReceipts = append(c.VoidedReceipts, request)

	return nil
}
