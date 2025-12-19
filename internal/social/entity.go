package social

import "feedsystem_video_go/internal/account"

type Social struct {
	ID         uint `gorm:"primaryKey"`
	FollowerID uint `gorm:"index"`
	VloggerID  uint `gorm:"index"`
}

type FollowRequest struct {
	VloggerID uint `json:"vlogger_id"`
}

type UnfollowRequest struct {
	VloggerID uint `json:"vlogger_id"`
}

type GetAllFollowersRequest struct {
	VloggerID uint `json:"vlogger_id"`
}

type GetAllFollowersResponse struct {
	Followers []*account.Account `json:"followers"`
}

type GetAllVloggersRequest struct {
	FollowerID uint `json:"follower_id"`
}

type GetAllVloggersResponse struct {
	Vloggers []*account.Account `json:"vloggers"`
}
