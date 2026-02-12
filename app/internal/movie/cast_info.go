package movie

import (
	"context"
	"github.com/atlas/slowpoke/pkg/state"
)

func StoreCastInfo(ctx context.Context, castId string, name string, info string) string {
	castInfo := CastInfo{
		CastId: castId,
		Name:   name,
		Info:   info,
	}
	state.SetState(ctx, castId, castInfo)
	return castId
}

func ReadCastInfos(ctx context.Context, castIds []string) []CastInfo {
	// Bulk
	var castInfos []CastInfo
	if len(castIds) > 0 {
		castInfos = state.GetBulkStateDefault[CastInfo](ctx, castIds, CastInfo{})
	} else {
		castInfos = make([]CastInfo, len(castIds))
	}
	return castInfos
}
