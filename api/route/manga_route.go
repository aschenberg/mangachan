package route

import (
	"manga/api/handler"
	"manga/config"
	"manga/pkg"
	"manga/pkg/logging"

	"manga/internal/infra/meili"
	pgsql "manga/internal/infra/pgsql/repository"
	"manga/internal/infra/rabbitmq"

	"manga/internal/usecase"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

func Manga(group *gin.RouterGroup, cfg *config.Config, pg *pkg.Postgres, redis *redis.Client, minio *minio.Client, rmq *rabbitmq.Task, ml meilisearch.ServiceManager, log logging.Logger) {
	meili := meili.NewManga(ml)
	uow := pgsql.NewUOWRepository(pg, cfg, minio, rmq, ml, log)
	mr := pgsql.NewMangaRepository(pg, log)
	mu := usecase.NewMangaUsecase(uow, mr, 10*time.Second, redis, cfg, meili,log)
	mh := handler.NewMangaHandler(mu)
	group.POST("/", mh.Create)
	group.GET("/source", mh.Source)
	group.GET("/:id", mh.GetById)
	group.GET("/all", mh.Find)
	// group.GET("/manga/:id", mc.FindByID)
}
