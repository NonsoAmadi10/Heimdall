package controllers

import (
	"errors"
	"log"
	"strconv"

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
	getBitcoinMetrics   = services.GetInfo
	getLightningMetrics = services.GetNodeInfo
	getConnectionMetric = services.FetchMetrics
	getAlerts           = services.FetchAlerts
	ackAlert            = services.AcknowledgeAlert
	resolveAlert        = services.ResolveAlert
)

func GetMetrics(c *fiber.Ctx) error {

	type NodeResponse struct {
		Lightning interface{} `json:"lightning"`
		Bitcoin   interface{} `json:"bitcoin"`
	}

	bitcoin, bitcoinErr := getBitcoinMetrics()
	lightning, lightningErr := getLightningMetrics()
	if bitcoinErr != nil || lightningErr != nil {
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
