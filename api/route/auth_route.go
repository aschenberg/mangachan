package route

import (
	"manga/api/handler"
	"manga/config"

	pgsql "manga/internal/infra/pgsql/repository"
	"manga/internal/usecase"
	"manga/pkg"
	"manga/pkg/logging"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func Auth(group *gin.RouterGroup, cfg *config.Config, pg *pkg.Postgres, oidcCfg *oauth2.Config, provider *oidc.Provider, log logging.Logger) {
	ar := pgsql.NewUserRepository(pg, log)
	au := usecase.NewLoginUsecase(ar, cfg, 10*time.Second, log)
	ah := handler.NewAuthHandler(au, cfg, oidcCfg, provider)
	group.POST("/login", ah.Login)
}
