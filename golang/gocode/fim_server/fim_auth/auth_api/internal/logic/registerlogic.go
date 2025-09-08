package logic

import (
	"context"
	"errors"
	"fim_server/fim_auth/auth_api/internal/svc"
	"fim_server/fim_auth/auth_api/internal/types"
	"fim_server/fim_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	if req.Pwd != req.RePwd {
		err = errors.New("两次密码不一致")
		return
	}

	response, err := l.svcCtx.UserRpc.UserCreate(context.Background(), &user_rpc.UserCreateRequest{
		NickName:       req.Nickname,
		Password:       req.RePwd,
		Role:           2,
		RegisterSource: "user",
		//OpenId:         uuid.New().String(),
	})
	if err != nil {
		return
	}

	return &types.RegisterResponse{UserID: uint(response.UserId)}, nil
}
