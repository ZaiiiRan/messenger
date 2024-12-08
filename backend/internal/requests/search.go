package requests

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type SearchRequest struct {
	PaginationRequest
	Search string `json:"search"`
}

// trim spaces in search request
func (s *SearchRequest) TrimSpaces() {
	s.Search = strings.TrimSpace(s.Search)
}

// Parse search request
func ParseSearchRequest(c *fiber.Ctx) (*SearchRequest, error) {
	return ParseRequest[SearchRequest](c)
}
