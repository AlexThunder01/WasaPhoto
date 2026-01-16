package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login enpoint
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User Endpoint
	rt.router.GET("/users", rt.wrap(rt.searchUsers))
	rt.router.PUT("/users/:user_id", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users/:user_id", rt.wrap(rt.getUsername))
	rt.router.GET("/users/:user_id/profile", rt.wrap(rt.getUserProfile))

	// Ban endpoint
	rt.router.PUT("/users/:user_id/ban/:ban_user_id", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:user_id/ban/:ban_user_id", rt.wrap(rt.unbanUser))

	// Followers endpoint
	rt.router.PUT("/users/:user_id/followed/:followed_user_id", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:user_id/followed/:followed_user_id", rt.wrap(rt.unfollowUser))

	// Stream endpoint
	rt.router.GET("/users/:user_id/home", rt.wrap(rt.getMyStream))

	// Photo Endpoint
	rt.router.POST("/users/:user_id/photos", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:user_id/photos/:photo_id", rt.wrap(rt.deletePhoto))
	rt.router.GET("/users/:user_id/photos/:photo_id", rt.wrap(rt.getPhoto))

	// Comments endpoint
	rt.router.POST("/users/:user_id/photos/:photo_id/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:user_id/photos/:photo_id/comments/:comment_id", rt.wrap(rt.uncommentPhoto))

	// Likes endpoint
	rt.router.PUT("/users/:user_id/photos/:photo_id/likes/:like_id", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:user_id/photos/:photo_id/likes/:like_id", rt.wrap(rt.unlikePhoto))

	// Icons endpoint
	rt.router.GET("/icons/:icon_id/", rt.wrap(rt.getServerIcon))
	rt.router.PUT("/users/:user_id/icon", rt.wrap(rt.setUserIcon))
	rt.router.GET("/users/:user_id/icon", rt.wrap(rt.getUserIconID))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
