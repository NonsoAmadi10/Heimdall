package controllers

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/NonsoAmadi10/p2p-analysis/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

var (
	getBitcoinMetrics      = services.GetInfo
	getLightningMetrics    = services.GetNodeInfo
	getConnectionMetric    = services.FetchMetrics
	getConnectionAnalytics = services.FetchMetricsAnalytics
	getAlerts              = services.FetchAlerts
	ackAlert               = services.AcknowledgeAlert
	resolveAlert           = services.ResolveAlert
	defaultAnalyticsWindow = 24 * time.Hour
	defaultAnalyticsBucket = 60
	minAnalyticsBucket     = 1
	maxAnalyticsBucket     = 24 * 60
)

func GetMetrics(c *fiber.Ctx) error {

	type NodeResponse struct {
		Lightning interface{} `json:"lightning"`
		Bitcoin   interface{} `json:"bitcoin"`
	}

	bitcoin, bitcoinErr := getBitcoinMetrics()
	lightning, lightningErr := getLightningMetrics()
	if bitcoinErr != nil {
		log.Printf("Failed to fetch bitcoin node metrics: %v", bitcoinErr)
	}
	if lightningErr != nil {
		log.Printf("Failed to fetch lightning node metrics: %v", lightningErr)
	}
	if bitcoinErr != nil && lightningErr != nil {
		log.Printf("Failed to fetch node metrics. bitcoin_error=%v lightning_error=%v", bitcoinErr, lightningErr)
		return c.Status(fiber.StatusServiceUnavailable).JSON(&Response{
			Success: false,
			Error:   "unable to fetch node information",
		})
	}

	response := &NodeResponse{
		Bitcoin:   bitcoin,
		Lightning: lightning,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetConnMetrics(c *fiber.Ctx) error {

	metrics, err := getConnectionMetric()
	if err != nil {
		log.Printf("Failed to fetch connection metrics: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&Response{
			Success: false,
			Error:   "unable to fetch connection metrics",
		})
	}

	response := &Response{
		Success: true,
		Data:    metrics,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetConnMetricsAnalytics(c *fiber.Ctx) error {
	now := time.Now().UTC()
	from := now.Add(-defaultAnalyticsWindow)
	to := now

	if fromQuery := c.Query("from"); fromQuery != "" {
		parsedFrom, err := time.Parse(time.RFC3339, fromQuery)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&Response{
				Success: false,
				Error:   "invalid 'from' timestamp, expected RFC3339 format",
			})
		}
		from = parsedFrom.UTC()
	}

	if toQuery := c.Query("to"); toQuery != "" {
		parsedTo, err := time.Parse(time.RFC3339, toQuery)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&Response{
				Success: false,
				Error:   "invalid 'to' timestamp, expected RFC3339 format",
			})
		}
		to = parsedTo.UTC()
	}

	intervalMinutes := defaultAnalyticsBucket
	if intervalQuery := c.Query("interval_minutes"); intervalQuery != "" {
		parsedInterval, err := strconv.Atoi(intervalQuery)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&Response{
				Success: false,
				Error:   "invalid 'interval_minutes', expected integer",
			})
		}
		intervalMinutes = parsedInterval
	}

	if intervalMinutes < minAnalyticsBucket || intervalMinutes > maxAnalyticsBucket {
		return c.Status(fiber.StatusBadRequest).JSON(&Response{
			Success: false,
			Error:   "interval_minutes must be between 1 and 1440",
		})
	}

	analytics, err := getConnectionAnalytics(from, to, time.Duration(intervalMinutes)*time.Minute)
	if err != nil {
		log.Printf("Failed to fetch connection analytics: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&Response{
			Success: false,
			Error:   "unable to fetch connection analytics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&Response{
		Success: true,
		Data:    analytics,
	})
}

func GetAlerts(c *fiber.Ctx) error {
	status := c.Query("status")
	alerts, err := getAlerts(status)
	if err != nil {
		log.Printf("Failed to fetch alerts: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(&Response{
			Success: false,
			Error:   "unable to fetch alerts",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&Response{
		Success: true,
		Data:    alerts,
	})
}

func AcknowledgeAlert(c *fiber.Ctx) error {
	alertID, err := parseAlertID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&Response{
			Success: false,
			Error:   "invalid alert id",
		})
	}

	alert, err := ackAlert(alertID)
	if err != nil {
		log.Printf("Failed to acknowledge alert: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&Response{
				Success: false,
				Error:   "alert not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(&Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&Response{
		Success: true,
		Data:    alert,
	})
}

func ResolveAlert(c *fiber.Ctx) error {
	alertID, err := parseAlertID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&Response{
			Success: false,
			Error:   "invalid alert id",
		})
	}

	alert, err := resolveAlert(alertID)
	if err != nil {
		log.Printf("Failed to resolve alert: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&Response{
				Success: false,
				Error:   "alert not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&Response{
			Success: false,
			Error:   "unable to resolve alert",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&Response{
		Success: true,
		Data:    alert,
	})
}

func parseAlertID(c *fiber.Ctx) (uint, error) {
	alertID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(alertID), nil
}
