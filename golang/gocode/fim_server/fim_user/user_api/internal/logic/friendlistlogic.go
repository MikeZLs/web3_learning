package logic

import (
	"context"
	"fim_server/common/list_query"
	"fim_server/common/models"
	"fim_server/fim_user/user_api/internal/svc"
	"fim_server/fim_user/user_api/internal/types"
	"fim_server/fim_user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListRequest) (resp *types.FriendListResponse, err error) {
	//var count int64
	//l.svcCtx.DB.Model(user_models.FriendModel{}).Where("send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID).Count(&count)
	//var friends []user_models.FriendModel
	//
	//if req.Limit <= 0 {
	//	req.Limit = 10
	//}
	//if req.Page <= 0 {
	//	req.Page = 1
	//}
	//
	//offset := (req.Page - 1) * req.Limit
	//
	//l.svcCtx.DB.Preload("SendUserModel").Preload("RevUserModel").Limit(req.Limit).Offset(offset).Find(&friends, "send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID)

	friends, count, _ := list_query.ListQuery(l.svcCtx.DB, user_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Where:   l.svcCtx.DB.Where("send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID),
		Preload: []string{"SendUserModel", "RevUserModel"},
	})

	// 查哪些用户在线
	//onlineMap := l.svcCtx.Redis.HGetAll("online").Val()
	//var onlineUserMap = map[uint]bool{}
	//for key, _ := range onlineMap {
	//	val, err1 := strconv.Atoi(key)
	//	if err1 != nil {
	//		logx.Error(err1)
	//		continue
	//	}
	//	onlineUserMap[uint(val)] = true
	//}

	var list []types.FriendInfoResponse
	for _, friend := range friends {

		info := types.FriendInfoResponse{}
		if friend.SendUserID == req.UserID {
			// 我是发起方
			info = types.FriendInfoResponse{
				UserID:   friend.RevUserID,
				Nickname: friend.RevUserModel.Nickname,
				Abstract: friend.RevUserModel.Abstract,
				Avatar:   friend.RevUserModel.Avatar,
				Notice:   friend.SenUserNotice,
				//IsOnline: onlineUserMap[friend.RevUserID],
			}
		}
		if friend.RevUserID == req.UserID {
			// 我是接收方
			info = types.FriendInfoResponse{
				UserID:   friend.SendUserID,
				Nickname: friend.SendUserModel.Nickname,
				Abstract: friend.SendUserModel.Abstract,
				Avatar:   friend.SendUserModel.Avatar,
				Notice:   friend.RevUserNotice,
				//IsOnline: onlineUserMap[friend.SendUserID],
			}
		}
		list = append(list, info)
	}

	return &types.FriendListResponse{List: list, Count: int(count)}, nil
}
