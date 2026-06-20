package model

// Common error messages
const (
	// Request validation errors
	ErrInvalidRequestFormat       = "Invalid request format"
	ErrInvalidEmailFormat         = "Invalid email format, must be a valid email address"
	ErrPasswordTooShort           = "Password must be at least 8 characters and contain only letters, numbers, and special characters"
	ErrInvalidNameFormat          = "Name must be 2-50 characters long and contain only letters and spaces"
	ErrInvalidPhoneFormat         = "Phone number must be 10 digits"
	ErrInvalidPincodeFormat       = "Pincode must be 6 digits"
	ErrEmptyStreetName            = "Street name cannot be empty"
	ErrEmptyLocality              = "Locality cannot be empty"
	ErrEmptyState                 = "State cannot be empty"
	ErrUserIDRequired             = "User ID is required"
	ErrAddressIDRequired          = "Address ID is required"
	ErrAuthorizationTokenRequired = "Authorization token required"
	ErrFailedGenerateToken        = "Failed to generate token"

	// Authentication errors
	ErrUserIDNotFound     = "User ID not found in context"
	ErrUnauthorizedModify = "Cannot modify another user's address"
	ErrUnauthorizedDelete = "Cannot delete another user's address"

	// Operation failures
	ErrLoginFailed             = "Login failed"
	ErrSignupFailed            = "Signup failed"
	ErrEmailVerificationFailed = "Email verification failed"
	ErrFailedRetrieveProfile   = "Failed to retrieve profile"
	ErrFailedUpdateProfile     = "Failed to update profile"
	ErrFailedRetrieveUser      = "Failed to retrieve user information"
	ErrFailedAddAddress        = "Failed to add address"
	ErrFailedRetrieveAddresses = "Failed to retrieve addresses"
	ErrFailedUpdateAddress     = "Failed to update address"
	ErrFailedDeleteAddress     = "Failed to delete address"
	ErrFailedBanUser           = "Failed to ban user"
	ErrFailedUnbanUser         = "Failed to unban user"
	ErrFailedCheckBan          = "Failed to check ban status"
	ErrFailedRetrieveUsers     = "Failed to retrieve users"
)

// Response messages
const (
	MsgAddressUpdated = "Address updated successfully"
	MsgAddressDeleted = "Address deleted successfully"
	MsgUserBanned     = "User banned successfully"
	MsgUserUnbanned   = "User unbanned successfully"
)