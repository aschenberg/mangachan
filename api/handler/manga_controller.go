package handler

import (
	"manga/api"
	"manga/internal/domain"
	"manga/internal/domain/dtos"
	"manga/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MangaController struct {
	MangaUsecase domain.MangaUsecase
}

func (mc *MangaController) Create(c *gin.Context) {
	var err error
	var create dtos.CreateManga

	err = c.ShouldBind(&create)
	if err != nil {
		api.RenderErrorResponse(c, "invalid request",
			pkg.WrapErrorf(err, pkg.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	err = mc.MangaUsecase.Create(c, create)
	if err != nil {
		api.RenderErrorResponse(c, "An error occurred: ",
			pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "An error occurred: "+err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dtos.SuccessResponse{Message: "created successfuly"})
}

func (mc *MangaController) FindByID(c *gin.Context) {
	var err error
	id := c.Param("id")

	result, err := mc.MangaUsecase.FindById(c, id)
	if err != nil {
		api.RenderErrorResponse(c, "An error occurred: ",
			pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "An error occurred: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, dtos.SuccessResponse{Message: "created successfuly", Data: result})
}
