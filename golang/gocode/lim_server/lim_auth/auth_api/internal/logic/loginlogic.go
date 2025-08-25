package logic

import (
	"context"
	"errors"
	"fim_server/lim_auth/auth_models"
	"fim_server/utils/jwts"
	"fim_server/utils/pwd"
	"fmt"

	"fim_server/lim_auth/auth_api/internal/svc"
	"fim_server/lim_auth/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	fmt.Println("AUTH_LOGIN")
	var user auth_models.UserModel
	err = l.svcCtx.DB.Take(&user, "id = ?", req.UserName).Error

	if err != nil {
		err = errors.New("用户名或密码错误")
		return
	}

	if !pwd.CheckPwd(user.Pwd, req.Password) {
		err = errors.New("用户名或密码错误")
		return
	}

	token, err := jwts.GenToken(jwts.JwtPayLoad{
		Nickname: user.Nickname,
		Role:     user.Role,
		UserID:   user.ID,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Error(err)
		err = errors.New("服务内部错误")
		return
	}

	return &types.LoginResponse{
		Token: token,
	}, nil

}
