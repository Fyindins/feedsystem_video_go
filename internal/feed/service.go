package feed

import (
	"context"
	"feedsystem_video_go/internal/video"
	"time"
)

type FeedService struct {
	repo     *FeedRepository
	likeRepo *video.LikeRepository
}

func NewFeedService(repo *FeedRepository, likeRepo *video.LikeRepository) *FeedService {
	return &FeedService{repo: repo, likeRepo: likeRepo}
}

func (f *FeedService) ListLatest(ctx context.Context, limit int, latestBefore time.Time) (ListLatestResponse, error) {
	videos, err := f.repo.ListLatest(ctx, limit, latestBefore)
	if err != nil {
		return ListLatestResponse{}, err
	}
	var nextTime time.Time
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].CreateTime
	} else {
		nextTime = time.Time{}
	}
	hasMore := len(videos) == limit
	feedVideos := make([]FeedVideoItem, 0, len(videos))
	for _, video := range videos {
		isLiked, err := f.likeRepo.IsLiked(ctx, video.ID, video.AuthorID)
		if err != nil {
			return ListLatestResponse{}, err
		}
		feedVideos = append(feedVideos, FeedVideoItem{
			ID:          video.ID,
			Author:      FeedAuthor{ID: video.AuthorID, Username: video.Username},
			Title:       video.Title,
			Description: video.Description,
			PlayURL:     video.PlayURL,
			CoverURL:    video.CoverURL,
			CreateTime:  video.CreateTime.Unix(),
			LikesCount:  video.LikesCount,
			IsLiked:     isLiked,
		})
	}
	resp := ListLatestResponse{
		VideoList: feedVideos,
		NextTime:  nextTime.Unix(),
		HasMore:   hasMore,
	}
	return resp, nil
}

func (f *FeedService) ListLikesCount(ctx context.Context, limit int, likesCountBefore int64) (ListLikesCountResponse, error) {
	videos, err := f.repo.ListLikesCount(ctx, limit, likesCountBefore)
	if err != nil {
		return ListLikesCountResponse{}, err
	}
	var nextLikesCountBefore int64
	if len(videos) > 0 {
		nextLikesCountBefore = videos[len(videos)-1].LikesCount
	} else {
		nextLikesCountBefore = 0
	}
	hasMore := len(videos) == limit
	feedVideos := make([]FeedVideoItem, 0, len(videos))
	for _, video := range videos {
		isLiked, err := f.likeRepo.IsLiked(ctx, video.ID, video.AuthorID)
		if err != nil {
			return ListLikesCountResponse{}, err
		}
		feedVideos = append(feedVideos, FeedVideoItem{
			ID:          video.ID,
			Author:      FeedAuthor{ID: video.AuthorID, Username: video.Username},
			Title:       video.Title,
			Description: video.Description,
			PlayURL:     video.PlayURL,
			CoverURL:    video.CoverURL,
			LikesCount:  video.LikesCount,
			IsLiked:     isLiked,
		})
	}
	resp := ListLikesCountResponse{
		VideoList:            feedVideos,
		NextLikesCountBefore: nextLikesCountBefore,
		HasMore:              hasMore,
	}
	return resp, nil
}
