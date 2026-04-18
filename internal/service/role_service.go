package service

import (
	"errors"

	"halo/internal/model"

	"gorm.io/gorm"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

func (s *RoleService) EnsureDefaults() error {
	roles := []model.Role{
		{Name: "super-role", DisplayName: "超级管理员", Description: "拥有全部权限"},
		{Name: "operator-role", DisplayName: "操作员", Description: "拥有内容管理权限"},
		{Name: "guest-role", DisplayName: "访客", Description: "基础权限"},
	}

	for _, role := range roles {
		var count int64
		s.db.Model(&model.Role{}).Where("name = ?", role.Name).Count(&count)
		if count == 0 {
			s.db.Create(&role)
		}
	}

	rules := []model.PolicyRule{
		{Name: "posts.*", Resource: "posts", Action: "*"},
		{Name: "pages.*", Resource: "pages", Action: "*"},
		{Name: "categories.*", Resource: "categories", Action: "*"},
		{Name: "tags.*", Resource: "tags", Action: "*"},
		{Name: "menus.*", Resource: "menus", Action: "*"},
		{Name: "comments.*", Resource: "comments", Action: "*"},
		{Name: "attachments.*", Resource: "attachments", Action: "*"},
		{Name: "backups.*", Resource: "backups", Action: "*"},
		{Name: "settings.*", Resource: "settings", Action: "*"},
		{Name: "users.*", Resource: "users", Action: "*"},
		{Name: "plugins.*", Resource: "plugins", Action: "*"},
		{Name: "themes.*", Resource: "themes", Action: "*"},
	}

	for _, rule := range rules {
		var count int64
		s.db.Model(&model.PolicyRule{}).Where("name = ?", rule.Name).Count(&count)
		if count == 0 {
			s.db.Create(&rule)
		}
	}

	return nil
}

func (s *RoleService) AssignRole(userID uint, roleName string) error {
	var role model.Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return errors.New("role not found")
	}

	binding := model.UserRole{
		UserID: userID,
		RoleID: role.ID,
	}

	return s.db.Create(&binding).Error
}

func (s *RoleService) GetUserRoles(userID uint) ([]model.Role, error) {
	var roles []model.Role
	var bindings []model.UserRole

	if err := s.db.Where("user_id = ?", userID).Find(&bindings).Error; err != nil {
		return nil, err
	}

	for _, b := range bindings {
		var role model.Role
		if err := s.db.First(&role, b.RoleID).Error; err == nil {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

func (s *RoleService) HasPermission(userID uint, resource, action string) bool {
	roles, _ := s.GetUserRoles(userID)

	for _, role := range roles {
		if role.Name == "super-role" {
			return true
		}

		var roleRules []model.RoleRule
		s.db.Where("role_id = ?", role.ID).Find(&roleRules)
		for _, rr := range roleRules {
			var rule model.PolicyRule
			if s.db.First(&rule, rr.RuleID).Error == nil {
				if (rule.Resource == resource || rule.Resource == "*") && (rule.Action == action || rule.Action == "*") {
					return true
				}
			}
		}
	}

	return false
}
