# Using Event Bus in the project

You should now understand how to use the Event Bus and Event Processor in your project.
It's time to apply them! Let's begin with the Event Bus.

## Exercise

Exercise path: ./project

1. Update your project to use Event Bus instead of Publisher.

Here are some tips on how to do this:

* Publish messages using the EventBus, not the Publisher directly.
* You should not do any JSON marshaling yourself. Just use the `JSONMarshaler` from Watermill in the EventBus config.
* Similarly, you don't need to create Watermill's `*message.Message` manually to publish it. Just pass the event struct to EventBus's `Publish`.

Don't forget to use the JSON marshaler with the custom `GenerateName` option:

```go
var marshaler = cqrs.JSONMarshaler{
	GenerateName: cqrs.StructName,
}
```

2. Add the Correlation ID decorator to the publisher.
   Without it, you will have a hard time debugging your application if something goes wrong.

You can use the decorator you implemented yourself, or use the `log.CorrelationPublisherDecorator` from
[`github.com/ThreeDotsLabs/go-event-driven/v2/common`](https://github.com/ThreeDotsLabs/go-event-driven/tree/main/common/log).

The decorator uses middleware from the [`github.com/ThreeDotsLabs/go-event-driven/v2/common/middleware`](https://github.com/ThreeDotsLabs/go-event-driven/blob/main/common/http/middlewares.go#L55) package.
This middleware extracts the correlation ID from the HTTP header and adds it to the context.
The decorator can then extract it using `log.CorrelationIDFromContext`.

After these changes, you should have much less boilerplate code in your project.
It'll also be much easier to add new events and handlers in the future.
