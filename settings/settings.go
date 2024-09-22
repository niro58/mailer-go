package setting

import "time"

const ACCESS_TOKEN_EXPIRY = 30 * time.Minute
const REFRESH_TOKEN_EXPIRY = 30 * 24 * time.Hour
var JWT_HEADER = map[string]string{
	"alg":"HS256",
	"typ":"JWT",
}