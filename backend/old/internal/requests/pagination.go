package requests

import "github.com/gofiber/fiber/v2"

type PaginationRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Parse pagination request
func ParsePaginationRequest(c *fiber.Ctx) (*PaginationRequest, error) {
	return ParseRequest[PaginationRequest](c)
}