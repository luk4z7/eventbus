package api

import (
	"eventhandler/entities"
	"eventhandler/worker"
	"net/http"

	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(w *worker.Worker) *echo.Echo {
	e := commonHTTP.NewEcho()

	e.POST("/transaction", func(c echo.Context) error {
		var request entities.TransactionRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		tracingID := uuid.New().String()
		request.Payload.Header = entities.NewHeader()

		w.Send(entities.Message{
			TracingID: tracingID,
			Data:      request.Payload,
		})

		return c.String(http.StatusCreated, tracingID)
	})

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	return e
}
