package handler

import (
	"manga/api/helper"
	"manga/internal/domain/dtos"
	"manga/internal/domain/enum"
	"manga/internal/usecase"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type MangaHandler struct {
	mu usecase.IMangaUsecase
}

func NewMangaHandler(mu usecase.IMangaUsecase) *MangaHandler {
	return &MangaHandler{
		mu: mu}
}

func (h *MangaHandler) Create(c *gin.Context) {
	var err error
	var create dtos.CreateManga

	err = c.ShouldBind(&create)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	res, err := h.mu.Create(c, create)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}
	if res.MangaID == 0 {
		c.JSON(http.StatusConflict, helper.GenerateBaseResponse(res, true, helper.Success))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, helper.Success))

}

func (h *MangaHandler) Source(c *gin.Context) {

	var sourceList []enum.SourceListResponse
	for _, name := range enum.SourceName {
		sourceList = append(sourceList, enum.SourceListResponse{Name: name.Name, Value: name.Source, Url: name.Url})
	}
	sort.Slice(sourceList, func(i, j int) bool {
		return sourceList[i].Value < sourceList[j].Value
	})
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(sourceList, true, helper.Success))

}

func (h *MangaHandler) GetById(c *gin.Context) {
	idParam := c.Param("id")

	res, err := h.mu.FindById(c, idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.GenerateBaseResponseWithAnyError(nil, false, helper.InternalError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, helper.Success))

}

// func (h *MangaHandler) FindByID(c *gin.Context) {
// 	var err error
// 	id := c.Param("id")

// 	result, err := h.MangaUsecase.FindById(c, id)
// 	if err != nil {
// 		api.RenderErrorResponse(c, "An error occurred: ",
// 			pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "An error occurred: "+err.Error()))
// 		return
// 	}
// 	c.JSON(http.StatusOK, dtos.SuccessResponse{Message: "created successfuly", Data: result})
// }
