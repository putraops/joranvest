package controllers

import (
	"net/http"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EducationCategoryController interface {
	GetPagination(context *gin.Context)
	Lookup(context *gin.Context)
	Save(context *gin.Context)
	GetById(context *gin.Context)
	DeleteById(context *gin.Context)
}

type educationCategoryController struct {
	educationCategoryService service.EducationCategoryService
	jwtService               service.JWTService
	db                       *gorm.DB
}

func NewEducationCategoryController(db *gorm.DB, jwtService service.JWTService) EducationCategoryController {
	return &educationCategoryController{
		db:                       db,
		jwtService:               jwtService,
		educationCategoryService: service.NewEducationCategoryService(db, jwtService),
	}
}

func (c educationCategoryController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	var result = c.educationCategoryService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

func (c educationCategoryController) Lookup(context *gin.Context) {
	var request helper.ReactSelectRequest

	errDTO := context.Bind(&request)
	if errDTO != nil {
		res := helper.StandartResult(false, errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	var result = c.educationCategoryService.Lookup(request)
	response := helper.StandartResult(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (r educationCategoryController) Save(c *gin.Context) {
	var result helper.Result
	var dto dto.EducationCategoryDto
	dto.Context = c

	errDto := c.Bind(&dto)
	if errDto != nil {
		res := helper.StandartResult(false, errDto.Error(), helper.EmptyObj{})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	result = r.educationCategoryService.Save(dto)
	c.JSON(http.StatusOK, helper.StandartResult(result.Status, result.Message, result.Data))
	return
}

func (c educationCategoryController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.educationCategoryService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c educationCategoryController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result = c.educationCategoryService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}
