# Project: Add Booking ID to TicketBookingConfirmed

We'll want to use `TicketBookingConfirmed` event for our read model.
But it misses the `BookingID` field which we'll use as the primary ID for the read model.

Thankfully, Dead Nation's API will pass through the `BookingID` field if we send it in the API request.

## Exercise

Exercise path: ./project

1. Extend the `TicketBookingConfirmed` event to include the `booking_id` field.
2. Update the `POST /tickets-status` HTTP handler. Pass the `booking_id` field from the request's payload to the `TicketBookingConfirmed` event when publishing it.
