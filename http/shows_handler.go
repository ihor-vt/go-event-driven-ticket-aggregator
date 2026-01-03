package http

import (
	"net/http"
	"tickets/entities"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ShowsRequest struct {
	DeadNationID    string    `json:"dead_nation_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	StartTime       time.Time `json:"start_time"`
	Title           string    `json:"title"`
	Venue           string    `json:"venue"`
}

type ShowResponse struct {
	ShowId string `json:"show_id"`
}

func (h Handler) PostShows(c echo.Context) error {
	var request ShowsRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	showID := uuid.New()
	if err := h.showRepo.AddShow(c.Request().Context(), entities.Show{
		ShowID:          showID.String(),
		DeadNationID:    request.DeadNationID,
		NumberOfTickets: request.NumberOfTickets,
		StartTime:       request.StartTime,
		Title:           request.Title,
		Venue:           request.Venue,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	return c.JSON(http.StatusCreated, ShowResponse{
		ShowId: showID.String(),
	})
}
