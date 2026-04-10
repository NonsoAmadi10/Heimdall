package controllers

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NonsoAmadi10/p2p-analysis/services"
	"github.com/NonsoAmadi10/p2p-analysis/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

func TestGetConnMetricsAnalyticsSuccess(t *testing.T) {
	originalAnalytics := getConnectionAnalytics
	defer func() {
		getConnectionAnalytics = originalAnalytics
	}()

	getConnectionAnalytics = func(from, to time.Time, interval time.Duration) (*services.MetricsAnalyticsResponse, error) {
		return &services.MetricsAnalyticsResponse{
			From:            from,
			To:              to,
			IntervalMinutes: int(interval.Minutes()),
			Points: []services.MetricsAnalyticsPoint{
				{Samples: 2},
			},
		}, nil
	}

	app := fiber.New()
	app.Get("/conn-metrics/analytics", GetConnMetricsAnalytics)

	req := httptest.NewRequest("GET", "/conn-metrics/analytics?interval_minutes=30", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestGetConnMetricsAnalyticsBadInterval(t *testing.T) {
	app := fiber.New()
	app.Get("/conn-metrics/analytics", GetConnMetricsAnalytics)

	req := httptest.NewRequest("GET", "/conn-metrics/analytics?interval_minutes=0", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestGetAlertsSuccess(t *testing.T) {
	originalGetAlerts := getAlerts
	defer func() {
		getAlerts = originalGetAlerts
	}()

	getAlerts = func(status string) ([]utils.Alert, error) {
		return []utils.Alert{{Type: "sync_stalled", Status: utils.AlertStatusOpen}}, nil
	}

	app := fiber.New()
	app.Get("/alerts", GetAlerts)

	req := httptest.NewRequest("GET", "/alerts?status=open", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestAcknowledgeAlertNotFound(t *testing.T) {
	originalAck := ackAlert
	defer func() {
		ackAlert = originalAck
	}()

	ackAlert = func(id uint) (*utils.Alert, error) {
		return nil, gorm.ErrRecordNotFound
	}

	app := fiber.New()
	app.Patch("/alerts/:id/ack", AcknowledgeAlert)

	req := httptest.NewRequest("PATCH", "/alerts/44/ack", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestResolveAlertSuccess(t *testing.T) {
	originalResolve := resolveAlert
	defer func() {
		resolveAlert = originalResolve
	}()

	resolveAlert = func(id uint) (*utils.Alert, error) {
		return &utils.Alert{ID: id, Status: utils.AlertStatusResolved}, nil
	}

	app := fiber.New()
	app.Patch("/alerts/:id/resolve", ResolveAlert)

	req := httptest.NewRequest("PATCH", "/alerts/1/resolve", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestAcknowledgeAlertBadRequest(t *testing.T) {
	originalAck := ackAlert
	defer func() {
		ackAlert = originalAck
	}()

	ackAlert = func(id uint) (*utils.Alert, error) {
		return nil, errors.New("bad state")
	}

	app := fiber.New()
	app.Patch("/alerts/:id/ack", AcknowledgeAlert)

	req := httptest.NewRequest("PATCH", "/alerts/2/ack", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}
