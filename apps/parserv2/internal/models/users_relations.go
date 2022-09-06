package model

type UserWitchToken struct {
	Users
	Token Tokens `gorm:"foreignKey:tokenId"`
}
