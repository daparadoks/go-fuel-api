package member

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

// Service -
type Service struct {
	DB *gorm.DB
}

// Member -
type Member struct {
	gorm.Model
	Username string
	Password string
	Mail     string
}

type MemberToken struct {
	gorm.Model
	Token       string
	DeviceToken string
	MemberId    uint
	ExpireDate  time.Time
}

// MemberService -
type MemberService interface {
	GetMember(username string) (Member, error)
}

// NewService -
func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

// GetMember -
func (s *Service) GetMember(username string) (Member, error) {
	var member Member
	if result := s.DB.First(&member).Where("Username = ?", username); result.Error != nil {
		return Member{}, result.Error
	}
	return member, nil
}

// GetMember -
func (s *Service) GetMemberById(id uint) (Member, error) {
	var member Member
	if result := s.DB.First(&member).Where("ID = ?", id); result.Error != nil {
		return Member{}, result.Error
	}
	return member, nil
}

// GetMember -
func (s *Service) GetMemberByMail(mail string) (Member, error) {
	var member Member
	if result := s.DB.First(&member).Where("Mail = ?", mail); result.Error != nil {
		return Member{}, result.Error
	}
	return member, nil
}

// GetMember -
func (s *Service) GetMemberByToken(token string) (Member, error) {
	var member Member
	memberToken, err := s.GetToken(token)
	if err != nil {
		return member, err
	}
	member, err = s.GetMemberById(memberToken.MemberId)
	return member, err
}

func (s *Service) Register(member Member) (Member, error) {
	if result := s.DB.Save(&member); result.Error != nil {
		return Member{}, result.Error
	}
	return member, nil
}

func (s *Service) GetToken(token string) (MemberToken, error) {
	var memberToken MemberToken
	result := s.DB.First(&memberToken).Where("Token=?", token)

	return memberToken, result.Error
}

func (s *Service) GetTokenByMemberId(deviceToken string, memberId uint) (MemberToken, error) {
	var token MemberToken
	result := s.DB.First(&token).Where("DeviceToken=? and MemberId=?", deviceToken, memberId)
	if result.Error == nil {
		return token, nil
	}

	println("m2")
	token, err := s.AddToken(deviceToken, memberId)
	if err != nil {
		return MemberToken{}, err
	}
	return token, nil

}

func (s *Service) AddToken(deviceToken string, memberId uint) (MemberToken, error) {
	now := time.Now()
	memberToken := RandStringRunes(30)
	if memberToken == "" {
		return MemberToken{}, nil
	}

	var token = MemberToken{
		Model:       gorm.Model{},
		Token:       memberToken,
		DeviceToken: deviceToken,
		MemberId:    memberId,
		ExpireDate:  now.AddDate(1, 0, 0),
	}
	if result := s.DB.Save(&token); result.Error != nil {
		return MemberToken{}, result.Error
	}

	return token, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
