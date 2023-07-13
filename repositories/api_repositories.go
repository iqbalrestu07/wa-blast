package repositories

// User ...
var User UserRepository

// Message ...
var Message MessageRepository

// Init ...
func Init() {
	Config()
	User = initUserRepository()
	Message = initMessageRepository()
}
