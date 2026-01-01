# Printing Tickets

{{background}}

Our operations team is now generating and printing tickets by hand.
This was a good strategy to roll out the product quickly, but it's not a good long-term solution because
they already struggle with the number of tickets they need to print.

Let's help our ops team by generating tickets for them!

We will use the `files` service to store the tickets.
The `files` will be available via the `gateway`, like other services.

{{endbackground}}

## Exercise

Exercise path: ./project

1. Implement an event handler triggered by the `TicketBookingConfirmed` event.

2. Store the ticket content.

Use the Files client from the [common library](https://github.com/ThreeDotsLabs/go-event-driven/blob/main/common/clients/clients.go).

```go
clients.Files.PutFilesFileIdContentWithTextBodyWithResponse(ctx, fileID, fileContent)
```

{{tip}}

Remember, these clients are initialized in your `main` function. Look for this line in `main.go`:
```go
apiClients, err := clients.NewClients(
```

{{endtip}}

The file name should have the format `[Ticket ID]-ticket.html`.
The content doesn't matter, it's just important that it contain the ticket ID, price, and amount.

You don't need to do anything on `TicketBookingCanceled` — the volume is low, and it's not a problem for the ops team to handle manually.

3. Handle events re-delivery.

Do you remember the discussion of eventual consistency? The client will return 409 when the file already exists.
We will use a similar strategy as in the previous module. If this error happens, you need to handle it gracefully.
It's worth adding a log in that situation, so you will know what happened in case any issues arise.

```go
import (
    "tickets/adapters"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

if resp.StatusCode() == http.StatusConflict {
	log.FromContext(ctx).With("file", fileID).Info("file already exists")
	return nil
}
```

{{tip}}

Note that we can add this functionality without changing any existing code.
In real life, it could be even implemented by a different team that has access to the events.

{{endtip}}
