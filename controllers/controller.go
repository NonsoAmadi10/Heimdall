package controllers

import (
	"github.com/NonsoAmadi10/p2p-analysis/services"
	"github.com/gofiber/fiber/v2"
)

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
