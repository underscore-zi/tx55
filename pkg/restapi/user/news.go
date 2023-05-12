package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelNone, "GET", "/news/list", getNewsList, nil, restapi.NewsJSON{})
}

// getNewsList godoc
// @Summary      Retrieve News List
// @Description  Retrieves all active news lists, one may have the title "policy" which is the terms of service
// @Tags         News
// @Produce      json
// @Success      200  {object}  restapi.ResponseJSON{data=[]restapi.NewsJSON{}}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /lobby/list [get]
func getNewsList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var news []models.News
	if err := db.Find(&news).Error; err != nil {
		l.WithError(err).Error("Error getting news list")
		restapi.Error(c, 500, "Error getting news list")
	} else {
		out := make([]restapi.NewsJSON, len(news))
		for i, n := range news {
			out[i] = restapi.NewsJSON{
				CreatedAt: n.CreatedAt,
				Title:     n.Topic,
				Content:   n.Body,
			}
		}
		restapi.Success(c, out)
	}
}