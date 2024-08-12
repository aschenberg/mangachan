package route

import (
	"manga/api/handler"
	"manga/config"
	"manga/internal/domain"
	"manga/internal/infra/ard"
	"manga/internal/usecase"
	"time"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func Auth(group *gin.RouterGroup, cfg *config.Config, db arangodb.Database, oidcCfg *oauth2.Config, provider *oidc.Provider) {
	ar := ard.NewUserRepository(db, domain.CollectionUser)
	au := usecase.NewLoginUsecase(ar, 10*time.Second)
	ah := handler.NewAuthHandler(au, cfg, oidcCfg, provider)
	group.POST("/login", ah.Login)
}
