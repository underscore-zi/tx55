package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelNone, "GET", "/news/list", getNewsList)
	restapi.Register(restapi.AuthLevelNone, "GET", "/policy", getPolicy)
}

// getNewsList godoc
// @Summary      Retrieve News List
// @Description  Retrieves all active news lists, one may have the title "policy" which is the terms of service
// @Tags         News
// @Produce      json
// @Success      200  {object}  restapi.ResponseJSON{data=[]restapi.NewsJSON{}}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /news/list [get]
func getNewsList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var news []models.News
	if err := db.Where("topic != 'policy'").Order("created_at desc").Find(&news).Error; err != nil {
		l.WithError(err).Error("Error getting news list")
		restapi.Error(c, 500, "Error getting news list")
	} else {
		out := make([]restapi.NewsJSON, len(news))
		for i, n := range news {
			out[i] = restapi.NewsJSON{
				ID:        n.ID,
				CreatedAt: n.CreatedAt,
				Title:     n.Topic,
				Content:   n.Body,
			}
		}
		restapi.Success(c, out)
	}
}

// getPolicy godoc
// @Summary      Retrieve Policy
// @Description  Retrieves the terms of service page shown ingame
// @Tags         News
// @Produce      json
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.NewsJSON{}}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /policy [get]
func getPolicy(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var policy models.News
	if err := db.Model(&policy).Find(&policy, "topic = 'policy'").Error; err != nil && err != gorm.ErrRecordNotFound {
		l.WithError(err).Error("Error getting policy")
		restapi.Error(c, 500, "Error getting news list")
	} else {
		restapi.Success(c, restapi.NewsJSON{
			ID:        policy.ID,
			CreatedAt: policy.CreatedAt,
			Title:     policy.Topic,
			Content:   policy.Body,
		})
	}
}
