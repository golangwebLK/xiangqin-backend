package user

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"time"
	"xiangqin-backend/utils"
)

type UserService struct {
	DB  *gorm.DB
	Jwt *utils.JWT
}

func NewUserService(db *gorm.DB, jwt *utils.JWT) *UserService {
	return &UserService{
		DB:  db,
		Jwt: jwt,
	}
}

func (svc *UserService) ComparePassword(req LoginReq) (User, error) {
	var user User
	if err := svc.DB.
		Where("username=?", req.Password).
		Find(&user).Error; err != nil {
		return User{}, errors.New("没有此用户")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return User{}, errors.New("密码错误")
	}
	return user, nil
}

func (svc *UserService) SignByID(companyCodeAndID string, exp time.Time) (string, error) {
	tokenStr, err := svc.Jwt.Sign(jwt.RegisteredClaims{
		Subject:   companyCodeAndID,
		ExpiresAt: jwt.NewNumericDate(exp),
	})
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (svc *UserService) StrConcatenation(id uint, code string) string {
	idStr := strconv.Itoa(int(id))
	companyCodeAndID := code + "@" + idStr
	return companyCodeAndID
}

func (svc *UserService) GetContent(user User) ([]*Content, error) {
	var permissions []Permission
	if err := svc.DB.Where("role=?", user.Role).Find(&permissions).Error; err != nil {
		return nil, err
	}
	var contents []Content
	if err := svc.DB.Find(&contents).Error; err != nil {
		return nil, err
	}
	permissionContents := make([]Content, 0, len(permissions)*2)
	idContentMap := make(map[string]Content, len(contents))
	for _, content := range contents {
		idContentMap[content.Code] = content
	}
	for _, permission := range permissions {
		pContents := SearchContentByPermission(permission.ContentID, idContentMap)
		permissionContents = append(permissionContents, pContents...)
	}
	contentTree := buildTree(permissionContents)
	return contentTree, nil
}
func buildTree(contents []Content) []*Content {
	nodeMap := make(map[string]*Content)
	var roots []*Content

	for _, node := range contents {
		node.Children = []*Content{}
		nodeMap[node.Code] = &node
		if node.ParentCode == "" {
			roots = append(roots, &node)
		}
	}

	for _, node := range contents {
		if parent, ok := nodeMap[node.ParentCode]; ok {
			parent.Children = append(parent.Children, &node)
		}
	}

	return roots
}
func SearchContentByPermission(code string, idContentMap map[string]Content) []Content {
	if code == "" {
		return make([]Content, 0, 3)
	}
	content := idContentMap[code]
	contents := SearchContentByPermission(content.ParentCode, idContentMap)
	contents = append(contents, content)
	return contents
}
