# Project: Publish Ticket Refunded Event

One more event we miss is the `TicketRefunded` event.
We'll use this data to update the read model when a ticket is refunded.

## Exercise

Exercise path: ./project

1. Introduce a new event: `TicketRefunded` containing `ticket_id`.

2. Publish it while handling the `RefundTicket` command.

{{hints}}
{{hint 1}}

This is how `TicketRefunded` event can look like:

```go
type TicketRefunded struct {
	Header MessageHeader `json:"header"`

	TicketID string `json:"ticket_id"`
}
```

{{endhint}}
{{endhints}}