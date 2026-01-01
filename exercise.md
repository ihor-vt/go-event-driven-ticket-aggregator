# Printing Tickets - emit event

{{background}}

One team in our company wants to integrate with the printing system to automate the printing of tickets.
They want to integrate with our system by subscribing to the `TicketPrinted` event.

They need information about the ticket ID and file name.

{{endbackground}}

## Exercise

Exercise path: ./project

Use the Event Bus to emit a `TicketPrinted` event after the ticket is printed.
You need to inject the Event Bus to your handler.

The emitted event should have the following format:

```go
type TicketPrinted struct {
	Header MessageHeader `json:"header"`

	TicketID string `json:"ticket_id"`
	FileName string `json:"file_name"`
}
```

{{tip}}

If you feel tempted to add the entire ticket model to the event, don't do it by default!

Remember that events become a contract between systems.
If you add an entire ticket model to the event, you will need to always keep adding this data to the event.

It's especially painful if you are refactoring in the future, and you want to split services or modules to smaller ones.
You may no longer have access to all the data that you emitted in the event in the past.

As an alternative, you can deprecate the old event and introduce a new one. 
However, it's always painful (as it may require a cross-team initiative).

[YAGNI!](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it)

{{endtip}}
