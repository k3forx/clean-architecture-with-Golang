package http

import (
	"net/http"
	"strconv"

	"github.com/k3forx/clean-architecture-with-Golang/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// ResponseError represents the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler represents the httphandler for article
type ArticleHandler struct {
	ArticleUseCase domain.ArticleUseCase
}

// NewArticleHandler will initialize the articles/resources endpoint
func NewArticleHandler(e *echo.Echo, us domain.ArticleUseCase) {
	handler := &ArticleHandler{
		ArticleUseCase: us,
	}
	e.GET("/articles", handler.Fetch)
	e.GET("/articles/:id", handler.GetById)
}

func (a *ArticleHandler) Fetch(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listAr, nextCursor, err := a.ArticleUseCase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

func (a *ArticleHandler) GetById(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}
	id := int64(idP)
	ctx := c.Request().Context()

	art, err := a.ArticleUseCase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, art)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
