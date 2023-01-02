package logic

import (
	"context"

	"my-rpc-gen/example/feeds/feeds"
	"my-rpc-gen/example/feeds/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFeedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedListLogic {
	return &GetFeedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// return all chats with bool for user { chat_id: int64, peer_type: int32, state: bool } req: { user_id: int64 }
func (l *GetFeedListLogic) GetFeedList(in *feeds.TLFeedGetFeedList) (*feeds.FeedListState, error) {
	// todo: add your logic here and delete this line

	return &feeds.FeedListState{}, nil
}
