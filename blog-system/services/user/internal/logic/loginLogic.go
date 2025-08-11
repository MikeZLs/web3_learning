package logic

import (
	"context"
	"user/internal/model"
	"user/userservice"

	"user/internal/svc"
	"user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	userModel := model.NewUsersModel(l.svcCtx.Mysql)
	user, err := userModel.FindByUsernameAndPwd(l.ctx, in.Username, in.Password)
	if err != nil {
		l.Logger.Error("查询用户失败: ", err)
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return &userservice.LoginResponse{
		Id:       user.Id,
		Username: user.Username,
	}, nil

}
