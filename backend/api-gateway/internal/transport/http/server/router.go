package httpserver

import (
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/handlers"
	commonmiddleware "github.com/ZaiiiRan/messenger/backend/go-common/pkg/middleware/http/server"
)

// registerRoutes wires every HTTP route to its handler, applying the
// appropriate authentication / permission middleware tier:
//
//   - public:        no token required
//   - authenticated: a valid access token is required (user may be unconfirmed)
//   - protected:     a valid access token AND an active, confirmed account
func (s *Server) registerRoutes(h *handlers.Handlers) {
	secret := []byte(s.cfg.JWT.AccessTokenSecret)

	api := s.app.Group("/api")

	// Middleware tiers reused across the route groups.
	requireAuth := commonmiddleware.UserAuthMiddleware(secret, commonmiddleware.ProtectAll())
	requirePermission := commonmiddleware.UserPermissionMiddleware(commonmiddleware.ProtectAll())

	authenticated := api.Group("", requireAuth)
	protected := api.Group("", requireAuth, requirePermission)

	// ---- auth (public) ----
	auth := api.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/logout", h.Logout)
	auth.Post("/confirm/link", h.ConfirmByLink)
	auth.Post("/password/forgot", h.ForgotPassword)
	auth.Post("/password/reset/code", h.ResetPasswordByCode)
	auth.Post("/password/reset/link", h.ResetPasswordByLink)

	// ---- auth (authenticated, account may still be unconfirmed) ----
	authenticated.Post("/auth/confirm/code", h.ConfirmByCode)
	authenticated.Post("/auth/confirmation-code", h.GetNewConfirmationCode)

	// ---- auth (protected) ----
	protected.Post("/auth/password/change", h.ChangePassword)
	protected.Post("/auth/email/change", h.ChangeEmail)
	protected.Post("/auth/email/confirm/code", h.ConfirmNewEmailByCode)
	protected.Post("/auth/email/confirm/link", h.ConfirmNewEmailByLink)
	protected.Get("/auth/sessions", h.GetActiveSessions)
	protected.Post("/auth/sessions/invalidate", h.InvalidateSessions)

	// ---- users (protected) ----
	protected.Get("/users/me", h.GetMe)
	protected.Patch("/users/me", h.UpdateMe)
	protected.Patch("/users/me/privacy", h.UpdateMyPrivacySettings)
	protected.Delete("/users/me", h.DeleteMe)

	// ---- social (protected) ----
	protected.Post("/social/users/by-ids", h.GetSocialUsersByIds)
	protected.Post("/social/users/by-usernames", h.GetSocialUsersByUsernames)
	protected.Get("/social/users/username/:username", h.GetSocialUserByUsername)
	protected.Get("/social/users/:id", h.GetSocialUserById)
	protected.Get("/social/friends", h.GetFriends)
	protected.Post("/social/friends", h.AddUsersToFriends)
	protected.Delete("/social/friends", h.RemoveUsersFromFriends)
	protected.Get("/social/friends/incoming", h.GetIncomingFriendRequests)
	protected.Get("/social/friends/outgoing", h.GetOutgoingFriendRequests)
	protected.Get("/social/blocked", h.GetBlockedUsers)
	protected.Post("/social/blocked", h.BlockUsers)
	protected.Delete("/social/blocked", h.UnblockUsers)
	protected.Get("/social/search", h.SearchUsers)
	protected.Patch("/social/me/privacy", h.UpdateSocialPrivacySettings)

	// ---- files (protected) ----
	protected.Post("/files", h.CreateFiles)
	protected.Post("/files/confirm", h.ConfirmUploads)
	protected.Post("/files/query", h.GetFiles)

	// ---- media (protected) ----
	protected.Post("/media/process", h.ProcessFiles)
}
