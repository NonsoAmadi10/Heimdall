package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/NonsoAmadi10/p2p-analysis/services"
	"github.com/NonsoAmadi10/p2p-analysis/utils"
	"github.com/gofiber/fiber/v2"
)

func TestGetMetricsSuccess(t *testing.T) {
	originalBitcoin := getBitcoinMetrics
	originalLightning := getLightningMetrics
	defer func() {
		getBitcoinMetrics = originalBitcoin
		getLightningMetrics = originalLightning
	}()

	getBitcoinMetrics = func() (*services.NodeMetrics, error) {
		return &services.NodeMetrics{Chain: "testnet"}, nil
	}
	getLightningMetrics = func() (*services.LNodeMetrics, error) {
		return &services.LNodeMetrics{Alias: "sample"}, nil
	}

	app := fiber.New()
	app.Get("/node-info", GetMetrics)

	req := httptest.NewRequest("GET", "/node-info", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestGetMetricsFailure(t *testing.T) {
	originalBitcoin := getBitcoinMetrics
	originalLightning := getLightningMetrics
	defer func() {
		getBitcoinMetrics = originalBitcoin
		getLightningMetrics = originalLightning
	}()

	getBitcoinMetrics = func() (*services.NodeMetrics, error) {
		return nil, fiber.ErrServiceUnavailable
	}
	getLightningMetrics = func() (*services.LNodeMetrics, error) {
		return nil, fiber.ErrServiceUnavailable
	}

	app := fiber.New()
	app.Get("/node-info", GetMetrics)

	req := httptest.NewRequest("GET", "/node-info", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusServiceUnavailable {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestGetConnMetricsSuccess(t *testing.T) {
	originalFetch := getConnectionMetric
	defer func() {
		getConnectionMetric = originalFetch
	}()

	getConnectionMetric = func() ([]utils.ConnectionMetrics, error) {
		return []utils.ConnectionMetrics{{BlockHeight: 100}}, nil
	}

	app := fiber.New()
	app.Get("/conn-metrics", GetConnMetrics)

	req := httptest.NewRequest("GET", "/conn-metrics", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	var payload Response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}
	if !payload.Success {
		t.Fatalf("expected success response")
	}
}

func TestGetConnMetricsFailure(t *testing.T) {
	originalFetch := getConnectionMetric
	defer func() {
		getConnectionMetric = originalFetch
	}()

	getConnectionMetric = func() ([]utils.ConnectionMetrics, error) {
		return nil, fiber.ErrInternalServerError
	}

	app := fiber.New()
	app.Get("/conn-metrics", GetConnMetrics)

	req := httptest.NewRequest("GET", "/conn-metrics", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}
