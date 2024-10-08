package openapi

/*
 * DKAM Service API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.1.0-dev
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

// package openapi

import (
	// "encoding/json"
	"encoding/json"
	// "fmt"
	"net/http"
	// "reflect"
	"strings"
	// "github.com/gorilla/mux"
	// "github.com/GIT_USER_ID/GIT_REPO_ID/dbpackage"
)

// HealthCheckAPIController binds http requests to an api service and writes the service results to the http response
type HealthCheckAPIController struct {
	service      HealthCheckAPIServicer
	errorHandler ErrorHandler
}

// HealthCheckAPIOption for how the controller is set up.
type HealthCheckAPIOption func(*HealthCheckAPIController)

// WithHealthCheckAPIErrorHandler inject ErrorHandler into controller
func WithHealthCheckAPIErrorHandler(h ErrorHandler) HealthCheckAPIOption {
	return func(c *HealthCheckAPIController) {
		c.errorHandler = h
	}
}

// NewHealthCheckAPIController creates a default api controller
func NewHealthCheckAPIController(s HealthCheckAPIServicer, opts ...HealthCheckAPIOption) Router {
	controller := &HealthCheckAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the HealthCheckAPIController
func (c *HealthCheckAPIController) Routes() Routes {
	return Routes{
		"GetHealthCheck": Route{
			strings.ToUpper("Get"),
			"/health_check",
			c.GetHealthCheck,
		},
	}
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GetHealthCheck - Retrieves all health_check
func (c *HealthCheckAPIController) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetHealthCheck(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}

	if result.Code == 404 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HealthResponse{
			Status:  "404",
			Message: "Table Does not exist",
		})
	} else if result.Code == 200 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HealthResponse{
			Status:  "200",
			Message: "Data Found!",
		})
		EncodeJSONResponse(result.Body, &result.Code, w)
	} else {
		EncodeJSONResponse(result.Body, &result.Code, w)
	}
}
