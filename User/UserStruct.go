package User

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"notNull" json:"nama"`
	NIK      string `gorm:"uniqueIndex;notNull;size:20" json:"nik"`
	Password string `gorm:"notNull" json:"pass"`
	Telepon  string `gorm:"uniqueIndex;notNull;size:20" json:"telp"`
	KotaAsal string `gorm:"null" json:"kota_asal"`
}

// Google Login & Register Google
type Google struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	EmailVerif bool   `json:"email_verified"`
	Locale     string `json:"locale"`
}

// Register regular user
type Register struct {
	Name     string `gorm:"notNull" json:"nama"`
	NIK      string `gorm:"uniqueIndex;notNull;size:20" json:"nik"`
	Password string `gorm:"notNull" json:"pass"`
	Telepon  string `gorm:"uniqueIndex;notNull;size:20" json:"telp"`
	KotaAsal string `gorm:"null" json:"kota_asal"`
}

// Login regular user
type Login struct {
	Telepon  string `gorm:"uniqueIndex;notNull;size:20" json:"telp"`
	Password string `gorm:"notNull" json:"pass"`
}
