package router

import (
	controller "github.com/delta/arcadia-backend/server/controller/auction"
	"github.com/delta/arcadia-backend/server/middleware"
	"github.com/gin-gonic/gin"
)

// Router for the Auction entity
func auctionRouter(superRoute *gin.RouterGroup) {
	auctionRoutes := superRoute.Group("/auction")

	auctionRoutes.Use(middleware.Auth)
	{
		auctionRoutes.GET("/tobuy/:id", controller.AuctionToBuyGET)
		auctionRoutes.GET("/auctioned/past/:id", controller.AuctionAuctionedPastGET)
		auctionRoutes.GET("/auctioned/current/:id", controller.AuctionAuctionedCurrentGET)
		auctionRoutes.POST("/sell", controller.SellPOST)

	}

}
