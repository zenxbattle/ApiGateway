package model

import "encoding/json"

// Socials represents social media links, matching the proto message Socials
type Socials struct {
	Github   string `json:"github"`
	Twitter  string `json:"twitter"`
	Linkedin string `json:"linkedin"`
}

// RegisterUserRequest for POST /api/v1/auth/register (JSON body)
type RegisterUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	// Country          string   `json:"country"`
	// Role             string   `json:"role"`
	// PrimaryLanguageID string  `json:"primaryLanguageID"`
	Email string `json:"email"`
	// AuthType         string   `json:"authType"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	// MuteNotifications bool    `json:"muteNotifications"`
	// Socials          Socials  `json:"socials"`
	// TwoFactorAuth    bool     `json:"twoFactorAuth"`
}

// LoginUserRequest for POST /api/v1/auth/login (JSON body)
type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

// TokenRefreshRequest for POST /api/v1/auth/token/refresh (JSON body)
type TokenRefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// LogoutRequest for POST /api/v1/auth/logout (JSON body)
type LogoutRequest struct {
	UserID string `json:"userId"`
}

// ResendOTPRequest for GET /api/v1/auth/otp/resend (query parameters)
type ResendOTPRequest struct {
	UserID string `form:"userId"`
}

// VerifyUserRequest for GET /api/v1/auth/verify (query parameters)
type VerifyUserRequest struct {
	Email string `form:"email"`
	Token string `form:"token"`
}

// ToggleTwoFactorAuthRequest for POST /api/v1/auth/2fa (JSON body)
type ToggleTwoFactorAuthRequest struct {
	// UserID        string `json:"userId"`
	Password      string `json:"password"`
	TwoFactorAuth bool   `json:"twoFactorAuth"`
}

// ForgotPasswordRequest for GET /api/v1/auth/password/forgot (query parameters)
type ForgotPasswordRequest struct {
	Email string `form:"email"`
}

// FinishForgotPasswordRequest for POST /api/v1/auth/password/reset (JSON body)
type FinishForgotPasswordRequest struct {
	Email           string `json:"email"`
	Token           string `json:"token"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

// ChangePasswordRequest for POST /api/v1/auth/password/change (JSON body)
type ChangePasswordRequest struct {
	// UserID          string `json:"userId"`
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

// LoginAdminRequest for POST /api/v1/admin/login (JSON body)
type LoginAdminRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateProfileRequest for PUT /api/v1/users/profile (JSON body)
type UpdateProfileRequest struct {
	// UserID            string  `json:"userId"`
	UserName          string  `json:"userName"`
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	Country           string  `json:"country"`
	PrimaryLanguageID string  `json:"primaryLanguageID"`
	MuteNotifications bool    `json:"muteNotifications"`
	Socials           Socials `json:"socials"`
	Bio               string  `json:"bio"`
}

// UpdateProfileImageRequest for PUT /api/v1/users/profile/image (JSON body)
// type UpdateProfileImageRequest struct {
// 	AvatarImageData string `json:"avatarImageData"`
// }

// GetUserProfileRequest for GET /api/v1/users/profile (query parameters)
type GetUserProfileRequest struct {
	UserID string `form:"userId"`
}

// CheckBanStatusRequest for GET /api/v1/users/ban/status (query parameters)
type CheckBanStatusRequest struct {
	UserID string `form:"userId"`
}

// BanHistoryRequest for GET /api/v1/users/ban/history (query parameters)
type BanHistoryRequest struct {
	UserID string `form:"userId"`
}

// SearchUsersRequest for GET /api/v1/users/search (query parameters)
type SearchUsersRequest struct {
	Query     string `form:"query"`
	PageToken string `form:"pageToken"`
	Limit     int32  `form:"limit"`
}

// FollowUserRequest for POST /api/v1/users/follow (JSON body)
type FollowUserRequest struct {
	FollowerID string `json:"followerID"`
	FolloweeID string `json:"followeeID"`
}

// UnfollowUserRequest for DELETE /api/v1/users/follow (query parameters)
type UnfollowUserRequest struct {
	FollowerID string `form:"followerID"`
	FolloweeID string `form:"followeeID"`
}

// GetFollowingRequest for GET /api/v1/users/following (query parameters)
type GetFollowingRequest struct {
	UserID    string `form:"userId"`
	PageToken string `form:"pageToken"`
	Limit     int32  `form:"limit"`
}

// GetFollowersRequest for GET /api/v1/users/followers (query parameters)
type GetFollowersRequest struct {
	UserID    string `form:"userId"`
	PageToken string `form:"pageToken"`
	Limit     int32  `form:"limit"`
}

// CreateUserAdminRequest for POST /api/v1/admin/users (JSON body)
type CreateUserAdminRequest struct {
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	Country           string  `json:"country"`
	Role              string  `json:"role"`
	PrimaryLanguageID string  `json:"primaryLanguageID"`
	Email             string  `json:"email"`
	AuthType          string  `json:"authType"`
	Password          string  `json:"password"`
	ConfirmPassword   string  `json:"confirmPassword"`
	MuteNotifications bool    `json:"muteNotifications"`
	Socials           Socials `json:"socials"`
}

// UpdateUserAdminRequest for PUT /api/v1/admin/users (JSON body)
type UpdateUserAdminRequest struct {
	UserID            string  `json:"userId"`
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
	Country           string  `json:"country"`
	Role              string  `json:"role"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	PrimaryLanguageID string  `json:"primaryLanguageID"`
	MuteNotifications bool    `json:"muteNotifications"`
	Socials           Socials `json:"socials"`
}

// BanUserRequest for POST /api/v1/admin/users/ban (JSON body)
type BanUserRequest struct {
	UserID    string `json:"userId"`
	Reason    string `json:"reason"`
	BanType   string `json:"banType"`
	BanReason string `json:"banReason"`
	BannedAt  int64  `json:"bannedAt"`
	BanExpiry int64  `json:"banExpiry"`
}

// UnbanUserRequest for DELETE /api/v1/admin/users/ban (query parameters)
type UnbanUserRequest struct {
	UserID string `form:"userId"`
}

// VerifyAdminUserRequest for POST /api/v1/admin/users/verify (JSON body)
type VerifyAdminUserRequest struct {
	UserID string `json:"userId"`
}

// UnverifyUserAdminRequest for POST /api/v1/admin/users/unverify (JSON body)
type UnverifyUserAdminRequest struct {
	UserID string `json:"userId"`
}

// SoftDeleteUserAdminRequest for DELETE /api/v1/admin/users (query parameters)
type SoftDeleteUserAdminRequest struct {
	UserID string `form:"userId"`
}

// GetAllUsersRequest for GET /api/v1/admin/users (query parameters)
type GetAllUsersRequest struct {
	PageToken      string `form:"pageToken"`
	Limit          int32  `form:"limit"`
	RoleFilter     string `form:"roleFilter"`
	StatusFilter   string `form:"statusFilter"`
	NameFilter     string `form:"nameFilter"`
	EmailFilter    string `form:"emailFilter"`
	FromDateFilter int64  `form:"fromDateFilter"`
	ToDateFilter   int64  `form:"toDateFilter"`
}

// SetUpTwoFactorAuthRequest for POST /api/v1/auth/2fa/setup (JSON body)
type SetUpTwoFactorAuthRequest struct {
	UserID   string `json:"userId"` //will be taken from JWT
	Password string `json:"password"`
}

type DisableTwoFactorAuthRequest struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

// type ExecuteUserCodeProblemIDRequest struct{
//     problems
// }

type UniversalExecutionResult struct {
	TotalTestCases  int            `json:"totalTestCases"`
	PassedTestCases int            `json:"passedTestCases"`
	FailedTestCases int            `json:"failedTestCases"`
	FailedTestCase  TestCaseResult `json:"failedTestCase,omitempty"`
	OverallPass     bool           `json:"overallPass"`
	SyntaxError     string         `json:"syntaxError"`
}
type TestCaseResult struct {
	TestCaseIndex int         `json:"testCaseIndex"`
	Input         interface{} `json:"input"`
	Expected      interface{} `json:"expected"`
	Received      interface{} `json:"received"`
	Error         string      `json:"error,omitempty"`
}

type TestCase struct {
	Input    json.RawMessage `json:"input"`
	Expected json.RawMessage `json:"expected"`
}

type GetSubmissionHistoryOptionalProblemId struct {
	UserID    string `json:"userId"`
	ProblemID string `json:"problemId,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}
