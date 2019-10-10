package router

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/traPtitech/booQ/model"
)

// PostComments POST /items/:id/comments
func PostComments(c echo.Context) error {
	ID := c.Param("id")
	itemID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user := c.Get("user").(model.User)
	user, err = model.GetUserByName(user.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	commentRequest := model.RequestPostCommentBody{}
	if err := c.Bind(&commentRequest); err != nil {
		return err
	}
	comment := model.Comment{
		ItemID: uint(itemID),
		UserID: user.ID,
		User:   user,
		Text:   commentRequest.Text,
	}
	res, err := model.CreateComment(comment)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, res)
}

// GetComments GET /comments
func GetComments(c echo.Context) error {
	userName := c.QueryParam("user")
	if userName != "" {
		user, err := model.GetUserByName(userName)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		res, err := model.GetCommentsByUserID(user.ID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, res)
	}
	res := []model.Comment{}
	return c.JSON(http.StatusCreated, res)
}