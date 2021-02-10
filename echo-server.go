package main

import (
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/labstack/echo/v4"
)

type Resp struct {
	Method string
	Path string
	Header http.Header
	QueryString string
	Data interface{}
	Error string
}

// ServerHeader middleware adds a `Server` header to the response.
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}

func main() {
	e := echo.New()

	// Debug mode
	e.Debug = true

	// Server header
	e.Use(ServerHeader)

	// Handler
	e.Any("/*", func(c echo.Context) error {
		r := new(Resp)
		r.Method = c.Request().Method
		r.Path = c.Request().RequestURI
		r.Header = c.Request().Header
		r.QueryString = c.QueryString()
		params, err := c.FormParams()
		if err == nil { 
			r.Data = params 
		}
		if r.Header["Content-Type"] != nil {
			ct := r.Header["Content-Type"][0]
			if strings.Contains(ct, "application/json") {
				if body := c.Request().Body; body != nil {
					bodyBytes, err := ioutil.ReadAll(body)
					if err == nil {
						var data interface{}
						err = json.Unmarshal(bodyBytes, &data)
						if err != nil {
							r.Error = "Error parsing body: " + err.Error()
						} else {
							r.Data = data
						}
					} else {
						r.Error = "Error reading body: " + err.Error()
					}
				}
			}
		}
		
		return c.JSON(http.StatusOK, r)
	})

	// Start server
	e.Logger.Fatal(e.Start(":1234"))
}
