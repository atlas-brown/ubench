package social

import (
	"context"
	"github.com/atlas/slowpoke/pkg/invoke"
	"github.com/atlas/slowpoke/pkg/state"
	"net/http"
)

func ReadUserTimeline(ctx context.Context, userId string) []Post {
	// postIds, err := state.GetState[[]string](ctx, userId)
	postIds, err := state.GetState[[]string](ctx, userId)
	if err != nil {
		return []Post{}
	}
	req := ReadPostsRequest{PostIds: postIds}
	postsResp := invoke.Invoke[ReadPostsResponse](ctx, "poststorage", "ro_read_posts", req, http.Request{})
	//fmt.Printf("Stored: %+v\nReturned: %+v\n", req, postsResp)
	return postsResp.Posts
}

func WriteUserTimeline(ctx context.Context, userId string, newPostIds []string) {
	// postIds, err := state.GetState[[]string](ctx, userId)
	postIds, err := state.GetState[[]string](ctx, userId)
	//fmt.Printf("[WriteUserTimeline] old postIds: %+v\n", postIds)
	//fmt.Printf("[WriteUserTimeline] to store: %+v\n", newPostIds)
	if err != nil {
		postIds = []string{}
	}
	if len(postIds) >= 10 {
		postIds = postIds[1:]
	}
	postIds = append(postIds, newPostIds...)
	//fmt.Printf("[WriteUserTimeline] new postIds: %+v\n", postIds)
	// state.SetState(ctx, userId, postIds)
	state.SetState(ctx, userId, postIds)
}
