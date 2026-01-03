# Calling Dead Nation

Remember that our {{exerciseLink "goal" "11-outbox" "01-project-replacing-dead-nation"}} is to migrate from the Dead Nation API to our own implementation.
To do that, we'll capture and forward requests to their API.

Previously, you exposed an endpoint that mimics the Dead Nation API.
Now you need to call the Dead Nation's endpoint.

We'll do it asynchronously, so we don't need to worry about their API being down or slow.
We'll use the `BookingMade` event that you emitted in the {{exerciseLink "previous exercise" "11-outbox" "10-project-emit-place-booked"}}.

## Exercise

Exercise path: ./project

1. Add a new event handler for the `BookingMade` event.

2. Map our `ShowID` to Dead Nation's `EventId`. We are calling the Dead Nation API, so we need to call it with their event ID, not our internal one. 
You can get Dead Nation's `EventId` from the database (it's stored in the `POST /shows` call, see {{exerciseLink "the previous exercise" "11-outbox" "02-project-store-show"}}).

This is intentionally a different name: EventID is a term used by Dead Nation, while we prefer the name `ShowID` so it's not confused with our internal ID.

You may need a new repository method to get the Show by its ID.

```go
func (s ShowsRepository) ShowByID(ctx context.Context, showID string) (entities.Show, error) {
	var show entities.Show
	err := s.db.GetContext(ctx, &show, `SELECT * FROM shows WHERE show_id = $1`, showID)
	if err != nil {
		return entities.Show{}, err
	}

	return show, nil
}
```

3. In the `BookingMade` handler, call the Dead Nation API to book a ticket. Remember to use `EventId`, not our internal `ShowID`.

There's a ready-to-use Dead Nation client in [the common library](https://github.com/ThreeDotsLabs/go-event-driven/blob/main/common/clients/clients.go).

{{tip}}

Remember, these clients are initialized in your `main` function. Look for this line in `main.go`:

```go
apiClients, err := clients.NewClients(
```

{{endtip}}

```go
resp, err := clients.DeadNation.PostTicketBookingWithResponse(
    ctx,
    dead_nation.PostTicketBookingRequest{
        BookingId:       booking.BookingID,
        EventId:         booking.DeadNationEventID,
        NumberOfTickets: booking.NumberOfTickets,
        CustomerAddress: booking.CustomerEmail,
    },
)
```

As it usually happens, the names of the external API don't match 1:1 with our names.
For example, `CustomerAddress` is `CustomerEmail` in our codebase.


{{tip}}

Adapters (like repositories or clients) are usually good spots to make the translation from external to internal language.
Like in this case, the Dead Nation API uses `CustomerAddress`, but we use `CustomerEmail`.

It allows us to keep the language inside our application consistent and free of external influences.

{{endtip}}

If everything goes fine, Dead Nation should call your `POST /ticket-status` endpoint.
