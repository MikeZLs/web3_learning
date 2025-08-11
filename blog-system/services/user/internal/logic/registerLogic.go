package logic

import (
	"context"
	"user/internal/model"

	"user/internal/svc"
	"user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	userModel := model.NewUsersModel(l.svcCtx.Mysql)
	username, err := userModel.FindOneByUsername(l.ctx, in.Username)

	if err != nil {
		l.Logger.Error("查询用户失败：", err)
		return nil, err
	}

	if username != nil {
		l.Logger.Error("用户已存在：", err)
		return nil, err
	}

	_, err = userModel.Insert(l.ctx, &model.Users{
		Username: in.Username,
		Password: in.Password,
		Email:    in.Email,
	})
	if err != nil {
		return nil, err
	}
	return &user.RegisterResponse{
		Message: "注册成功",
	}, nil
}
