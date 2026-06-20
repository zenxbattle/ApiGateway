package model

import (
	"time"
)

// UserProfile represents a user's profile information
type UserProfile struct {
	UserID                  string                  `json:"userId"`
	UserName                string                  `json:"userName"`
	FirstName               string                  `json:"firstName"`
	LastName                string                  `json:"lastName"`
	AvatarURL               string                  `json:"avatarURL"`
	Email                   string                  `json:"email"`
	Role                    string                  `json:"role"`
	Bio                     string                  `json:"bio"`
	Country                 string                  `json:"country"`
	IsBanned                bool                    `json:"isBanned"`
	AuthType                string                  `json:"authType"`
	IsVerified              bool                    `json:"isVerified"`
	PrimaryLanguageID       string                  `json:"primaryLanguageID"`
	MuteNotifications       bool                    `json:"muteNotifications"`
	Socials                 Socials                 `json:"socials"`
	CreatedAt               int64                   `json:"createdAt"`
	BanID                   string                  `gorm:"type:varchar(255)" json:"banID"`
	BanReason               string                  `gorm:"type:varchar(255)" json:"banReason"`
	BanExpiration           int64                   `json:"banExpiration"`
	TwoFactorEnabled        bool                    `gorm:"default:false;not null" json:"twoFactorEnabled"`
	ProblemSolvedStatsCount ProblemSolvedStatsCount `json:"problemSolvedStatsCount"`
}

type ProblemSolvedStatsCount struct {
	EasyCount      int `json:"easyCount"`
	MediumCount    int `json:"mediumCount"`
	HardCount      int `json:"hardCount"`
	MaxEasyCount   int `json:"maxEasyCount"`
	MaxMediumCount int `json:"maxMediumCount"`
	MaxHardCount   int `json:"maxHardCount"`
}

// to add new problem submission entry that is unique and first, this is used for leaderboard querying, Need- inorder to do faster query we cut down the struct and store this in another place.
type ProblemDone struct {
	UserID      string    `json:"userId" bson:"userId"`
	ProblemID   string    `json:"problemId" bson:"problemId"`
	Title       string    `json:"title" bson:"title"`
	Language    string    `json:"language" bson:"language"`
	Difficulty  string    `json:"difficulty" bson:"difficulty"`
	SubmittedAt time.Time `json:"submittedAt" bson:"submitted_at"`
}

type SubmissionHistoryResponse struct {
	Submissions []Submission `json:"submissions"`
}

type Submission struct {
	ID            string    `json:"id" bson:"id"`
	UserID        string    `json:"userId" bson:"userId"`
	ProblemID     string    `json:"problemId" bson:"problemId"`
	ChallengeID   string    `json:"challengeId,omitempty" bson:"challenge_id"`
	SubmittedAt   time.Time `json:"submittedAt" bson:"submitted_at"`
	Status        string    `json:"status" bson:"status"`
	Output        string    `json:"output,omitempty" bson:"output"`
	UserCode      string    `json:"userCode" bson:"userCode"`
	Language      string    `json:"language" bson:"language"`
	Score         int       `json:"score" bson:"score"`
	ExecutionTime float64   `json:"executionTime,omitempty" bson:"execution_time"`
	Difficulty    string    `json:"difficulty" bson:"difficulty"`
	IsFirst       bool      `json:"isFirst" bson:"is_first"`
	Title         string    `json:"title"`
}

// BanHistory represents a single ban record
type BanHistory struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	BannedAt  int64  `json:"bannedAt"`
	BanType   string `json:"banType"`
	BanReason string `json:"banReason"`
	BanExpiry int64  `json:"banExpiry"`
}

// Response structs for each handler
type RegisterUserResponse struct {
	UserID       string      `json:"userId"`
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresIn    int32       `json:"expiresIn"`
	UserProfile  UserProfile `json:"userProfile"`
	Message      string      `json:"message"`
}

type LoginUserResponse struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	ExpiresIn    int32       `json:"expiresIn"`
	UserID       string      `json:"userId"`
	UserProfile  UserProfile `json:"userProfile"`
	Message      string      `json:"message"`
}

type LoginAdminResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int32  `json:"expiresIn"`
	AdminID      string `json:"adminID"`
	Message      string `json:"message"`
}

type TokenRefreshResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int32  `json:"expiresIn"`
	UserID      string `json:"userId"`
	Message     string `json:"message"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type ResendEmailVerificationResponse struct {
	Message  string `json:"message"`
	ExpiryAt int64  `json:"expiryAt"`
}

type VerifyUserResponse struct {
	Message string `json:"message"`
}

type ToggleTwoFactorAuthResponse struct {
	Message string `json:"message"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type FinishForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}

type UpdateProfileResponse struct {
	Message     string      `json:"message"`
	UserProfile UserProfile `json:"userProfile"`
}

type UpdateProfileImageResponse struct {
	Message   string `json:"message"`
	AvatarURL string `json:"avatarURL"`
}

type GetUserProfileResponse struct {
	UserProfile UserProfile `json:"userProfile"`
	Message     string      `json:"message,omitempty"`
}

type CheckBanStatusResponse struct {
	IsBanned      bool   `json:"isBanned"`
	Reason        string `json:"reason"`
	BanExpiration int64  `json:"banExpiration"`
	Message       string `json:"message"`
}

type FollowUserResponse struct {
	Message string `json:"message"`
}

type UnfollowUserResponse struct {
	Message string `json:"message"`
}

type GetFollowingResponse struct {
	Users         []UserProfile `json:"users"`
	TotalCount    int32         `json:"totalCount"`
	NextPageToken string        `json:"nextPageToken"`
	Message       string        `json:"message"`
}

type GetFollowersResponse struct {
	Users         []UserProfile `json:"users"`
	TotalCount    int32         `json:"totalCount"`
	NextPageToken string        `json:"nextPageToken"`
	Message       string        `json:"message"`
}

type CreateUserAdminResponse struct {
	UserID  string `json:"userId"`
	Message string `json:"message"`
}

type UpdateUserAdminResponse struct {
	Message string `json:"message"`
}

type BanUserResponse struct {
	Message string `json:"message"`
}

type UnbanUserResponse struct {
	Message string `json:"message"`
}

type VerifyAdminUserResponse struct {
	Message string `json:"message"`
}

type UnverifyUserAdminResponse struct {
	Message string `json:"message"`
}

type SoftDeleteUserAdminResponse struct {
	Message string `json:"message"`
}

type GetAllUsersResponse struct {
	Users         []UserProfile `json:"users"`
	TotalCount    int32         `json:"totalCount"`
	NextPageToken string        `json:"nextPageToken"`
	PrevPageToken string        `json:"prevPageToken"`
	Message       string        `json:"message"`
}

type BanHistoryResponse struct {
	Bans    []BanHistory `json:"bans"`
	Message string       `json:"message"`
}

type SearchUsersResponse struct {
	Users         []UserProfile `json:"users"`
	TotalCount    int32         `json:"totalCount"`
	NextPageToken string        `json:"nextPageToken"`
	Message       string        `json:"message"`
}

// GenericResponse remains unchanged
type GenericResponse struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	ErrorType string `json:"type"`
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
}

type SetUpTwoFactorAuthResponse struct {
	Image   string `json:"image"`
	Secret  string `json:"secret"`
	Message string `json:"message"`
}

type GetTwoFactorAuthStatusResponse struct {
	IsEnabled bool   `json:"isEnabled"`
	Message   string `json:"message"`
}

type ExecutionResultJSON struct {
	TestCaseIndex int    `json:"testCaseIndex"`
	Nums          []int  `json:"nums"`
	Target        int    `json:"target"`
	Expected      []int  `json:"expected"`
	Received      []int  `json:"received"`
	Passed        bool   `json:"passed"`
	Summary       string `json:"summary,omitempty"`
}

type ProblemsDoneStatistics struct {
	MaxEasyCount    int32 `json:"maxEasyCount" bson:"maxEasyCount"`
	DoneEasyCount   int32 `json:"doneEasyCount" bson:"doneEasyCount"`
	MaxMediumCount  int32 `json:"maxMediumCount" bson:"maxMediumCount"`
	DoneMediumCount int32 `json:"doneMediumCount" bson:"doneMediumCount"`
	MaxHardCount    int32 `json:"maxHardCount" bson:"maxHardCount"`
	DoneHardCount   int32 `json:"doneHardCount" bson:"doneHardCount"`
}
