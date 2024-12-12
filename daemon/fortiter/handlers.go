package fortiter

import (
	"fmt"
	"github.com/algorandfoundation/algorun-tui/api"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handlers struct {
	PrometheusHandler echo.HandlerFunc
	client            api.ClientWithResponses
	db                *sqlx.DB
}

var (
	upgrader = websocket.Upgrader{}
)

func (h Handlers) GetAgreementEvents(c echo.Context, hash string) error {
	//rows, err := h.db.NamedQuery(`SELECT * FROM agreements WHERE hash=:first_name`)
	//if err != nil {
	//	return err
	//}
	return c.String(http.StatusOK, "Hello")
}

func (h Handlers) GetStatus(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type LogWssMessage struct {
	method string
}

func (h Handlers) GetWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			break
			//c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
			//c.Logger().Error(err)
		}
		fmt.Printf("{ \"message\": \"%s\" }\n", msg)
	}
	return nil
}

func (h Handlers) GetMetrics(c echo.Context) error {
	return h.PrometheusHandler(c)
}
