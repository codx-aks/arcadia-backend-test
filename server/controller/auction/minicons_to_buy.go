package controller

import (
	"net/http"

	"github.com/delta/arcadia-backend/database"
	helper "github.com/delta/arcadia-backend/server/helper/general"
	"github.com/delta/arcadia-backend/server/model"
	"github.com/delta/arcadia-backend/utils"
	"github.com/gin-gonic/gin"
)

func AuctionToBuyGET(c *gin.Context) {
	id := c.Param("id")
	var auctionItems []AuctionItem

	db := database.GetDB()

	log := utils.GetControllerLogger("/api/auction/tobuy/:id [GET]")

	if err := db.Model(model.Auction{}).
		Omit("updated_at", "deleted_at").
		Preload("OwnedMinicon").
		Preload("OwnedMinicon.Minicon").
		Where("seller_id != ? AND status = ?", id, "ongoing").
		Find(&auctionItems).Error; err != nil {
		log.Errorln(err)
		helper.SendError(c, http.StatusInternalServerError, "Unable to get auction items, Please try again later")
		return
	}

	res := GetAuctionBuyerResponse{
		AuctionItems: auctionItems,
	}

	helper.SendResponse(c, http.StatusOK, res)
}
