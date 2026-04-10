package controllers

import (
	"log"

	"github.com/NonsoAmadi10/p2p-analysis/services"
	"github.com/gofiber/fiber/v2"
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
