package route

import (
	"manga/config"
	"manga/db"
	"manga/internal/api/controller"
	"manga/internal/domain/models"
	"manga/internal/repository"
	"manga/internal/usecase"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewMangaRouter(env *config.Config, timeout time.Duration, db db.MongoDB, group *gin.RouterGroup, redis *redis.Client) {
	mr := repository.NewMangaRepository(db, models.CollectionManga)
	mu := usecase.NewMangaUsecase(mr, timeout, redis)
	mc := &controller.MangaController{
		MangaUsecase: mu,
	}
	group.POST("/manga", mc.Create)
	group.GET("/manga/:id", mc.FindByID)
}
