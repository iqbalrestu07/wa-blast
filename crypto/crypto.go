package crypto

const (
	TokenTypeBearer = "Bearer"
	TokenTypeSecret = "Secret"
)

func Init() {
	initJwtKey()
}
