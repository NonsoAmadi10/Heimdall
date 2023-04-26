package controllers

import (
	"github.com/NonsoAmadi10/p2p-analysis/services"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func GetMetrics(c *fiber.Ctx) error {

	type NodeResponse struct {
		Lightning interface{} `json:"lightning"`
		Bitcoin   interface{} `json:"bitcoin"`
	}

	bitcoin := services.GetInfo()
	lightning := services.GetNodeInfo()

	response := &NodeResponse{
		Bitcoin:   bitcoin,
		Lightning: lightning,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetConnMetrics(c *fiber.Ctx) error {

	metrics := services.FetchMetrics()

	response := &Response{
		Success: true,
		Data:    metrics,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
