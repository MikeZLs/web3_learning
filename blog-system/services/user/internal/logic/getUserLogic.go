package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"user/internal/model"

	"user/internal/svc"
	"user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserRequest) (*user.GetUserResponse, error) {
	userId, err := l.ctx.Value("id").(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	userModel := model.NewUsersModel(l.svcCtx.Mysql)

	user, err := userModel.FindOne(l.ctx, userId)

	if err != nil && (errors.Is(err, model.ErrNotFound) || errors.Is(err, sql.ErrNoRows)) {
		return nil, err
	}

	if err != nil {
		l.Logger.Error("查询用户失败: ", err)
		return nil, err
	}

	return &user.GetUserResponse{
		username: user.Username,
		email:    user.Email,
	}, nil
}
