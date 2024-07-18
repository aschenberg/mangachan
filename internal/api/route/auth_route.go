package route

import (
	"manga/api/controller"
	"manga/config"
	"manga/db"
	"manga/domain/models"
	"manga/repository"
	"manga/usecase"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func NewAuthRouter(env *config.Env, timeout time.Duration, db db.MongoDB, group *gin.RouterGroup, oidc *oauth2.Config, provider *oidc.Provider) {
	ur := repository.NewUserRepository(db, models.CollectionUser)
	lu := usecase.NewLoginUsecase(ur, timeout)
	lc := &controller.AuthController{
		LoginUsecase: lu,
		Env:          env,
		Oidc:         oidc,
		Provider:     provider,
	}
	group.POST("/authenticate", lc.Login)
}
