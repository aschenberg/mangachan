package route

import (
	"manga/api/controller"
	"manga/config"
	"manga/db"
	"manga/domain/models"
	"manga/repository"
	"manga/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewMangaRouter(env *config.Env, timeout time.Duration, db db.MongoDB, group *gin.RouterGroup, redis *redis.Client) {
	mr := repository.NewMangaRepository(db, models.CollectionManga)
	mu := usecase.NewMangaUsecase(mr, timeout, redis)
	mc := &controller.MangaController{
		MangaUsecase: mu,
	}
	group.POST("/manga", mc.Create)
	group.GET("/manga/:id", mc.FindByID)
}
