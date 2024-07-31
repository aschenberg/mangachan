package route

import (
	"manga/config"
	"manga/db"
	"manga/internal/api/controller"
	"manga/internal/domain/models"
	"manga/internal/repository"
	"manga/internal/usecase"

	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func NewAuthRouter(conf *config.Config, timeout time.Duration, db db.MongoDB, group *gin.RouterGroup, oidc *oauth2.Config, provider *oidc.Provider) {
	ur := repository.NewUserRepository(db, models.CollectionUser)
	lu := usecase.NewLoginUsecase(ur, timeout)
	lc := &controller.AuthController{
		LoginUsecase: lu,
		Conf:         conf,
		Oidc:         oidc,
		Provider:     provider,
	}
	group.POST("/authenticate", lc.Login)
}
