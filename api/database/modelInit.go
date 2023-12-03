package database

func Migrate() {
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Token{})
	DB.AutoMigrate(&Credential{})
}
