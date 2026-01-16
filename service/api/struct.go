package api

// Comment structure

type Comment struct {
	Text      string `json:"text"`       // The comment text
	PhotoID   int64  `json:"photo_id"`   // The photo ID
	AuthorID  int64  `json:"author_id"`  // The author ID
	CommentID int64  `json:"comment_id"` // The comment ID
	Datetime  string `json:"datetime"`   // The Date and time of a photo
}

// Like structure

type Like struct {
	UserID int64 `json:"user_id"` // The Author ID
}

// Photo structure

type Photo struct {
	Datetime string    `json:"datetime"`  // The Date and time of a photo
	Likes    []Like    `json:"likes"`     // The Likes number of a photo
	Comments []Comment `json:"comments"`  // The Comments number of a photo
	PhotoID  int64     `json:"photo_id"`  // The Photo ID
	AuthorID int64     `json:"author_id"` // The Author ID
	Username string    `json:"username"`  // The Author username
}

// Error structure

type Error struct {
	Code    string `json:"code"`    // The Code error
	Message string `json:"message"` // The Message error
}

// User structure

type User struct {
	UserID   int64  `json:"userID"`   // The User ID
	Username string `json:"username"` // The username
}

// Profile structure

type Profile struct {
	Owner       string  `json:"owner"`     // The Owner ID
	Photos      []Photo `json:"photos"`    // The Profile's photos
	FollowersID []int64 `json:"followers"` // The Followers of a profile
	FollowedsID []int64 `json:"followed"`  // The Followed of a profile
	IconID      string  `json:"icon_id"`   // The IconID of a profile
}
