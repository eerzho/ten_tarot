package constant

type UserState string

const (
	UserDonateState  UserState = "donate"
	UserSupportState UserState = "support"
	UserDefaultState UserState = ""
)
