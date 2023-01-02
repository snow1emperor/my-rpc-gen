package logic

import (
	"context"

	"my-rpc-gen/example/feeds_old/feeds"
	"my-rpc-gen/example/feeds_old/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFeedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFeedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFeedListLogic {
	return &UpdateFeedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// send array with { chat_id: int64, peer_type: int32, state: bool }
func (l *UpdateFeedListLogic) UpdateFeedList(in *feeds.TLUpdateFeedList) (*feeds.UpdateFeedListStatus, error) {
	// todo: add your logic here and delete this line

	return &feeds.UpdateFeedListStatus{}, nil
}
