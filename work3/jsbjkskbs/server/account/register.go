package account

import (
	"context"
	"log"
	"server/datastruct"
	"server/mysql"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Register(ctx context.Context, c *app.RequestContext) {
	var reigisterStruct struct {
		Username string `form:"username" json:"username" vd:"(len($)>0&&len($)<64); msg:'Illegal Username'"`
		Password string `form:"password" json:"password" vd:"(len($)>5&&len($)<16); msg:'Illegal Password'"`
	}

	if err := c.BindAndValidate(&reigisterStruct); err != nil {
		c.JSON(consts.StatusOK, datastruct.ShortResponse{
			Status:  consts.StatusOK,
			Message: "",
			Error:   err.Error(),
		})
		return
	}

	accountList, err := mysql.MySQLAccountSearch(reigisterStruct.Username)

	if err != nil {
		c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
			Status:  consts.StatusBadRequest,
			Message: "",
			Error:   err.Error(),
		})
		return
	}

	if len(accountList) != 0 {
		c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
			Status:  consts.StatusBadRequest,
			Message: "",
			Error:   "username has already existed",
		})
		return
	}

	err = mysql.MySQLAccountCreate(reigisterStruct.Username, reigisterStruct.Password)

	if err != nil {
		c.JSON(consts.StatusBadRequest, datastruct.ShortResponse{
			Status:  consts.StatusBadRequest,
			Message: "user creating failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, datastruct.ShortResponse{
		Status:  consts.StatusOK,
		Message: "ok",
		Error:   "",
	})

	log.Printf("[INFO] User [%s] has registered successfully.", reigisterStruct.Username)
}
