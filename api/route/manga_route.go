package route

import (
	"manga/api/handler"
	"manga/config"
	"manga/pkg/logging"
	"manga/pkg/postgres"

	pgsql "manga/internal/infra/pgsql/repository"

	"manga/internal/usecase"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

func Manga(group *gin.RouterGroup, cfg *config.Config, pg *postgres.Postgres, redis *redis.Client, minio *minio.Client, log logging.Logger) {
	uow := pgsql.NewUOWRepository(pg, log)
	mu := usecase.NewMangaUsecase(uow, 10*time.Second, redis)
	mh := handler.NewMangaHandler(mu)
	group.POST("/", mh.Create)
	// group.GET("/manga/:id", mc.FindByID)
}
