# Malformed messages

As much as retrying is useful most of the time, it doesn't work in all cases.

One example is a *malformed message*. This is a message that cannot be processed not because of an error on the handler side
but because the handler doesn't understand it. This could be a broken JSON payload or a message sent to the wrong topic.

In this scenario, **you need to remove the message from the queue.**

We will show a few more advanced tools for this in future modules.
But sometimes simple approaches are more than enough.

For now, let's consider a dead-simple approach to handling a malformed message: acknowledge and ignore it.
Your handler can remove the message by returning `nil` early instead of an error.

For example, if it's an invalid message schema, you can check the metadata for it:

```go
if msg.Metadata.Get("type") != "booking.created" {
	slog.Error("Invalid message type")
	return nil
}
```

It's always worth logging the message payload so you can investigate it later.
Otherwise, you won't know what went wrong.

### Removing a particular message

If there's a particular message that got published by mistake and can't be unmarshalled,
you can check its UUID.

```go
if msg.UUID == "5f810ce3-222b-4626-bc04-cbfb460c98c7" {
	return nil
}
```

It's a primitive way of doing this, but it works and might be good enough for your use case.
It helps if you have a healthy CI/CD pipeline and can quickly deploy a new version of the service.
Sometimes that's a pragmatic choice if you were to spend too much time on this.

{{tip}}

If you use this method, you need to make sure all your messages have a unique UUID!
This is why you shouldn't publish the same message twice.

{{endtip}}

### Handling permanent errors

It may happen that errors on the business domain level are not retryable.
If you know this beforehand, you can create a dedicated error type and middleware that handles it.

```go
type PermanentError interface {
	IsPermanent() bool
}

func SkipPermanentErrorsMiddleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		msgs, err := h(msg)

		var permErr PermanentError
		if errors.As(err, &permErr) && permErr.IsPermanent() {
			return nil, nil
		}

		return msgs, err
	}
}
```

You can then use it in your application logic.
For example, if the message misses a critical field, there's no point in retrying it.
It's a good idea to raise some kind of alert when this error occurs.

```go
type MissingInvoiceNumber struct {}

func (m MissingInvoiceNumber) Error() string {
	return "missing the invoice number - can't continue"
}

func (m MissingInvoiceNumber) IsPermanent() bool {
	return true
}
```

## Exercise

Exercise path: ./project

Handle two types of malformed messages in your project:

1. **A `TicketBookingConfirmed` event with invalid JSON payload and ID `2beaf5bc-d5e4-4653-b075-2b36bbf28949` has been published.**

Update your handlers to detect, ignore, and acknowledge this specific message.

2. **An event was published to the `TicketBookingConfirmed` topic, but its `type` metadata was incorrectly set to `TicketBooking` instead of `TicketBookingConfirmed`.**

When publishing `TicketBookingConfirmed` and `TicketBookingCanceled` events from your service, always set the correct `type` metadata so the event type can be verified.  
**Later, ensure that all events on the topic have the correct `type` metadata and that it matches the expected value.**
Make your handlers ignore and acknowledge invalid messages.
