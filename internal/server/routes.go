package server

import "github.com/gin-gonic/gin"

// registerRoutes registers API routes and returns the router
func (s *GinServer) registerRoutes(router *gin.Engine) *gin.Engine {
	callGroup := router.Group("/calls")

	callGroup.POST("/", s.AddCall)

	callGroup.GET("/", s.GetCalls)
	callGroup.GET("/:id", s.GetCallInfo)

	callGroup.PATCH("/:id/status", s.UpdateCallStatus)

	callGroup.DELETE("/:id", s.DeleteCall)

	return router
}
