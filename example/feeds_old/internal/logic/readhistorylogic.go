package logic

import (
	"context"

	"my-rpc-gen/example/feeds_old/feeds"
	"my-rpc-gen/example/feeds_old/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReadHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadHistoryLogic {
	return &ReadHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// for user req: { user_id: int64 }
func (l *ReadHistoryLogic) ReadHistory(in *feeds.TLFeedReadHistory) (*feeds.HistoryList, error) {
	// todo: add your logic here and delete this line

	return &feeds.HistoryList{}, nil
}
