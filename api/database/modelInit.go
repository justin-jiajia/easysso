package database

func Migrate() {
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&ServerCode2Token{})
	DB.AutoMigrate(&Credential{})
	DB.AutoMigrate(&UserToken{})
	DB.AutoMigrate(&UserLog{})
}
