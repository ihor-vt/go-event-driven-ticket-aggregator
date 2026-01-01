package adapters

import (
	"context"
	"sync"

	"tickets/entities"
)

type ReceiptsServiceStub struct {
	lock sync.Mutex

	IssuedReceipts map[string]entities.IssueReceiptRequest
}

func (c *ReceiptsServiceStub) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.IssuedReceipts[request.IdempotencyKey] = request

	return nil
}
