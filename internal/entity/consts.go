package entity

const (
	LenRefreshToken          = 32
	ExpiresMinuteAccessToken = 15
	ExpiresDayRefreshToken   = 30
	HeaderAccessToken        = "AccessToken"
	HeaderRefreshToken       = "RefreshToken"

	LenRegistrationCode = 6
)

var AllSymbol = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
