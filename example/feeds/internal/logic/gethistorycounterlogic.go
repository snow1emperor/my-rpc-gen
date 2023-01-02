package logic

import (
	"context"

	"my-rpc-gen/example/feeds/feeds"
	"my-rpc-gen/example/feeds/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHistoryCounterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHistoryCounterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHistoryCounterLogic {
	return &GetHistoryCounterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// for user req: { user_id: int64 }
func (l *GetHistoryCounterLogic) GetHistoryCounter(in *feeds.TLGetHistoryCounter) (*feeds.HistoryCounterState, error) {
	// todo: add your logic here and delete this line

	return &feeds.HistoryCounterState{}, nil
}
