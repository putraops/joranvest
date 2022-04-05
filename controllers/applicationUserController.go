package controllers

import (
	"fmt"
	"net/http"

	"joranvest/commons"
	"joranvest/dto"
	"joranvest/helper"
	"joranvest/models"
	"joranvest/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApplicationUserController interface {
	GetPagination(context *gin.Context)
	Lookup(context *gin.Context)
	UserLookup(context *gin.Context)
	UpdateProfile(context *gin.Context)
	ChangePhone(context *gin.Context)
	ChangePassword(context *gin.Context)
	Profile(context *gin.Context)
	GetAll(context *gin.Context)
	GetById(context *gin.Context)
	GetViewById(context *gin.Context)
	DeleteById(context *gin.Context)
	RecoverPassword(context *gin.Context)
	ResetPassword(context *gin.Context)
	EmailVerificationById(context *gin.Context)
}

type applicationUserController struct {
	DB                     *gorm.DB
	applicationUserService service.ApplicationUserService
	jwtService             service.JWTService
}

func NewApplicationUserController(db *gorm.DB, jwtService service.JWTService) ApplicationUserController {
	return &applicationUserController{
		DB:                     db,
		applicationUserService: service.NewApplicationUserService(db),
		jwtService:             jwtService,
	}
}

// @Tags         ApplicationUser
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body commons.Pagination2ndRequest true "body"
// @Success      200 {object} object
// @Failure 	 400,404 {object} object
// @Router       /application_user/getPagination [post]
func (c *applicationUserController) GetPagination(context *gin.Context) {
	var req commons.Pagination2ndRequest
	errDTO := context.Bind(&req)
	if errDTO != nil {
		res := helper.BuildResponse(false, errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	var result = c.applicationUserService.GetPagination(req)
	context.JSON(http.StatusOK, result)
}

func (c *applicationUserController) Lookup(context *gin.Context) {
	var request helper.ReactSelectRequest
	qry := context.Request.URL.Query()

	if _, found := qry["q"]; found {
		request.Q = fmt.Sprint(qry["q"][0])
	}
	request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	var result = c.applicationUserService.Lookup(request)
	response := helper.BuildResponse(true, "Ok", result.Data)
	context.JSON(http.StatusOK, response)
}

func (c *applicationUserController) UserLookup(context *gin.Context) {
	var request helper.ReactSelectRequest
	//qry := context.Request.URL.Query()

	// if _, found := qry["q"]; found {
	// 	request.Q = fmt.Sprint(qry["q"][0])
	// }
	//request.Field = helper.StringifyToArray(fmt.Sprint(qry["field"]))

	// var result = c.applicationUserService.Lookup(request)

	err := context.Bind(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, helper.StandartResult(false, err.Error(), nil))
		return
	}

	var result = c.applicationUserService.UserLookup(request)
	context.JSON(http.StatusOK, result)
}

func (c *applicationUserController) GetAll(context *gin.Context) {
	var users = c.applicationUserService.GetAll()
	res := helper.BuildResponse(true, "Ok", users)
	context.JSON(http.StatusOK, res)
}

func (c *applicationUserController) UpdateProfile(context *gin.Context) {
	var dto dto.ApplicationUserDescriptionDto
	errDTO := context.ShouldBind(&dto)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)

	dto.UpdatedBy = fmt.Sprintf("%v", claims["user_id"])
	result := c.applicationUserService.UpdateProfile(dto)
	if result.Status {
		response := helper.BuildResponse(result.Status, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildResponse(result.Status, result.Message, helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *applicationUserController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	user := c.applicationUserService.UserProfile(id)

	res := helper.BuildResponse(true, "Ok!", user)
	context.JSON(http.StatusOK, res)
}

func (c *applicationUserController) DeleteById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get Id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	var result helper.Response
	result = c.applicationUserService.DeleteById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	} else {
		response := helper.BuildResponse(true, "Ok", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *applicationUserController) GetById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.applicationUserService.GetById(id)
	if !result.Status {
		response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	} else {
		response := helper.BuildResponse(true, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	}
}

func (c *applicationUserController) GetViewById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	result := c.applicationUserService.GetViewById(id)
	if result.Status {
		response := helper.BuildResponse(result.Status, "Ok", result.Data)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildResponse(result.Status, result.Message, helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
}

func (c *applicationUserController) ChangePassword(context *gin.Context) {
	var recordDto dto.ChangePasswordDto
	err := context.ShouldBind(&recordDto)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to request login", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	result := c.applicationUserService.ChangePassword(recordDto)
	if result.Status {
		if v, ok := (result.Data).(models.ApplicationUser); ok {
			generatedToken := c.jwtService.GenerateToken(v.Id, v.EntityId)
			v.Token = generatedToken

			response := helper.BuildResponse(true, "Ok!", v)
			context.JSON(http.StatusOK, response)
			return
		}
	} else {
		if result.Errors == "NotFound" {
			response := helper.BuildErrorResponse("Error", result.Message, helper.EmptyObj{})
			context.JSON(http.StatusBadRequest, response)
		} else {
			response := helper.BuildResponse(false, result.Message, helper.EmptyObj{})
			context.JSON(http.StatusOK, response)
		}
	}
}

func (c *applicationUserController) ChangePhone(context *gin.Context) {
	var recordDto dto.ChangePhoneDto
	err := context.ShouldBind(&recordDto)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to bind request", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	res := c.applicationUserService.ChangePhone(recordDto)

	if res.Status {
		response := helper.BuildResponse(true, "Ok!", res.Data)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildResponse(false, res.Message, helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	}
	return
}

func (c *applicationUserController) RecoverPassword(context *gin.Context) {
	var recordDto dto.RecoverPasswordDto
	err := context.ShouldBind(&recordDto)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to request dto", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	result := c.applicationUserService.RecoverPassword(recordDto)
	response := helper.BuildResponse(true, "Ok!", result)
	context.JSON(http.StatusOK, response)
}

func (c *applicationUserController) ResetPassword(context *gin.Context) {
	var email string
	qry := context.Request.URL.Query()

	if _, found := qry["email"]; found {
		email = fmt.Sprint(qry["email"][0])
	}

	response := c.applicationUserService.ResetPasswordByEmail(email)
	context.JSON(http.StatusOK, response)
}

func (c *applicationUserController) EmailVerificationById(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		response := helper.BuildErrorResponse("Failed to get id", "Error", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}

	result := c.applicationUserService.EmailVerificationById(id)
	response := helper.BuildResponse(result.Status, result.Message, helper.EmptyObj{})
	context.JSON(http.StatusOK, response)
}
