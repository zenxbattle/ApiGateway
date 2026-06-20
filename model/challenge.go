package model

import (
	"time"
)

const ChallengeOpen = "CHALLENGEOPEN"
const ChallengeClose = "CHALLENGECLOSE"
const ChallengeStarted = "CHALLENGESTARTED"
const ChallengeForfieted = "CHALLENGEFORFIETED"
const ChallengeEnded = "CHALLENGEENDED"

// Challenge represents a coding challenge
type Challenge struct {
	ID                  string                         `bson:"_id"`
	Title               string                         `bson:"title"`
	CreatorID           string                         `bson:"creator_id"`
	Difficulty          string                         `bson:"difficulty"`
	IsPrivate           bool                           `bson:"is_private"`
	Status              string                         `bson:"status"`
	Password            string                         `bson:"password"` // Only for private challenges
	ProblemIDs          []string                       `bson:"problemIds"`
	TimeLimit           int32                          `bson:"time_limit"`
	CreatedAt           int64                          `bson:"created_at"`
	IsActive            bool                           `bson:"is_active"`
	ParticipantIDs      []string                       `bson:"participant_ids"`
	UserProblemMetadata map[string]ProblemMetadataList `bson:"user_problem_metadata"`
	StartTime           int64                          `bson:"start_time"`
	EndTime             int64                          `bson:"end_time"`
}

// ChallengeProblemMetadata represents metadata for a problem in a challenge
type ChallengeProblemMetadata struct {
	ProblemID   string `bson:"problemId"`
	Score       int32  `bson:"score"`
	TimeTaken   int64  `bson:"time_taken"`
	CompletedAt int64  `bson:"completed_at"`
}

// ProblemMetadataList holds a list of challenge problem metadata
type ProblemMetadataList struct {
	ChallengeProblemMetadata []ChallengeProblemMetadata `bson:"challenge_problem_metadata"`
}

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	UserID            string `bson:"userId"`
	ProblemsCompleted int32  `bson:"problems_completed"`
	TotalScore        int32  `bson:"total_score"`
	Rank              int32  `bson:"rank"`
}

// UserStats represents user statistics across challenges
type UserStats struct {
	UserID              string                   `bson:"userId"`
	ProblemsCompleted   int32                    `bson:"problems_completed"`
	TotalTimeTaken      int64                    `bson:"total_time_taken"`
	ChallengesCompleted int32                    `bson:"challenges_completed"`
	Score               float64                  `bson:"score"`
	ChallengeStats      map[string]ChallengeStat `bson:"challenge_stats"`
}

// ChallengeStat represents user stats for a specific challenge
type ChallengeStat struct {
	Rank              int32 `bson:"rank"`
	ProblemsCompleted int32 `bson:"problems_completed"`
	TotalScore        int32 `bson:"total_score"`
}

// CreateChallengeRequest represents the request to create a challenge
type CreateChallengeRequest struct {
	Title      string    `bson:"title"`
	CreatorID  string    `bson:"creator_id"`
	Difficulty string    `bson:"difficulty"`
	IsPrivate  bool      `bson:"is_private"`
	ProblemIDs []string  `bson:"problemIds"`
	TimeLimit  int32     `bson:"time_limit"`
	StartAt    time.Time `bson:"start_at"`
}

// CreateChallengeResponse represents the response for creating a challenge
type CreateChallengeResponse struct {
	ID       string `bson:"id"`
	Password string `bson:"password"` // Only for private challenges
	JoinURL  string `bson:"join_url"`
}

// GetChallengeDetailsRequest represents the request to get challenge details
type GetChallengeDetailsRequest struct {
	ID     string `bson:"id"`
	UserID string `bson:"userId"`
}

// GetChallengeDetailsResponse represents the response for getting challenge details
type GetChallengeDetailsResponse struct {
	Challenge    Challenge           `bson:"challenge"`
	Leaderboard  []LeaderboardEntry  `bson:"leaderboard"`
	UserMetadata ProblemMetadataList `bson:"user_metadata"`
}

// GetPublicChallengesRequest represents the request to get public challenges
type GetPublicChallengesRequest struct {
	Difficulty     string `bson:"difficulty"`
	IsActive       bool   `bson:"is_active"`
	Page           int32  `bson:"page"`
	PageSize       int32  `bson:"page_size"`
	UserID         string `bson:"userId"`
	IncludePrivate bool   `bson:"include_private"`
}

// GetPublicChallengesResponse represents the response for getting public challenges
type GetPublicChallengesResponse struct {
	Challenges []Challenge `bson:"challenges"`
}

// JoinChallengeRequest represents the request to join a challenge
type JoinChallengeRequest struct {
	ChallengeID string  `bson:"challenge_id"`
	UserID      string  `bson:"userId"`
	Password    *string `bson:"password"` // Optional, required for private challenges
}

// JoinChallengeResponse represents the response for joining a challenge
type JoinChallengeResponse struct {
	ChallengeID string `bson:"challenge_id"`
	Success     bool   `bson:"success"`
	Message     string `bson:"message"`
}

// StartChallengeRequest represents the request to start a challenge
type StartChallengeRequest struct {
	ChallengeID string `bson:"challenge_id"`
	UserID      string `bson:"userId"`
}

// StartChallengeResponse represents the response for starting a challenge
type StartChallengeResponse struct {
	Success   bool  `bson:"success"`
	StartTime int64 `bson:"start_time"`
}

// EndChallengeRequest represents the request to end a challenge
type EndChallengeRequest struct {
	ChallengeID string `bson:"challenge_id"`
	UserID      string `bson:"userId"`
}

// EndChallengeResponse represents the response for ending a challenge
type EndChallengeResponse struct {
	Success     bool               `bson:"success"`
	Leaderboard []LeaderboardEntry `bson:"leaderboard"`
}

// GetSubmissionStatusRequest represents the request to get submission status
type GetSubmissionStatusRequest struct {
	SubmissionID string `bson:"submission_id"`
}

// GetSubmissionStatusResponse represents the response for getting submission status
type GetSubmissionStatusResponse struct {
	Submission Submission `bson:"submission"`
}

// GetChallengeSubmissionsRequest represents the request to get challenge submissions
type GetChallengeSubmissionsRequest struct {
	ChallengeID string `bson:"challenge_id"`
}

// GetChallengeSubmissionsResponse represents the response for getting challenge submissions
type GetChallengeSubmissionsResponse struct {
	Submissions []Submission `bson:"submissions"`
}

// GetUserStatsRequest represents the request to get user stats
type GetUserStatsRequest struct {
	UserID string `bson:"userId"`
}

// GetUserStatsResponse represents the response for getting user stats
type GetUserStatsResponse struct {
	Stats UserStats `bson:"stats"`
}

// GetChallengeUserStatsRequest represents the request to get challenge-specific user stats
type GetChallengeUserStatsRequest struct {
	ChallengeID string `bson:"challenge_id"`
	UserID      string `bson:"userId"`
}

// GetChallengeUserStatsResponse represents the response for getting challenge-specific user stats
type GetChallengeUserStatsResponse struct {
	UserID                   string                     `bson:"userId"`
	ProblemsCompleted        int32                      `bson:"problems_completed"`
	TotalScore               int32                      `bson:"total_score"`
	Rank                     int32                      `bson:"rank"`
	ChallengeProblemMetadata []ChallengeProblemMetadata `bson:"challenge_problem_metadata"`
}
