package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gozero-learn/blog-api/internal/config"
	"gozero-learn/blog-api/internal/model"
)

type ServiceContext struct {
	Config        config.Config
	UsersModel    model.UsersModel
	PostsModel    model.PostsModel
	CommentsModel model.CommentsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Database.DataSource)
	return &ServiceContext{
		Config:        c,
		UsersModel:    model.NewUsersModel(conn),
		PostsModel:    model.NewPostsModel(conn),
		CommentsModel: model.NewCommentsModel(conn),
	}
}
