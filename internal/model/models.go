package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64" json:"username"`
	Password  string    `json:"-"`
	Role      string    `gorm:"size:32" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255" json:"title"`
	Slug      string    `gorm:"uniqueIndex;size:255" json:"slug"`
	Content   string    `gorm:"type:text" json:"content"`
	Excerpt   string    `gorm:"type:text" json:"excerpt"`
	Template  string    `gorm:"size:255" json:"template"`
	Category  string    `gorm:"size:128" json:"category"`
	Tags      string    `gorm:"type:text" json:"tags"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Page struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255" json:"title"`
	Slug      string    `gorm:"uniqueIndex;size:255" json:"slug"`
	Content   string    `gorm:"type:text" json:"content"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:128" json:"name"`
	Slug        string    `gorm:"uniqueIndex;size:128" json:"slug"`
	DisplayName string    `gorm:"size:255" json:"displayName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:128" json:"name"`
	Slug      string    `gorm:"uniqueIndex;size:128" json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Menu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:128" json:"name"`
	Items     string    `gorm:"type:text" json:"items"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `json:"postId"`
	Author    string    `gorm:"size:128" json:"author"`
	Email     string    `gorm:"size:255" json:"email"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Reply struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CommentID uint      `gorm:"index" json:"commentId"`
	Author    string    `gorm:"size:128" json:"author"`
	Email     string    `gorm:"size:255" json:"email"`
	Content   string    `gorm:"type:text" json:"content"`
	Status   string    `gorm:"size:32;default:pending" json:"status"`
	Allow    bool      `gorm:"-" json:"allow"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserRegistration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64" json:"username"`
	Email     string    `gorm:"uniqueIndex;size:255" json:"email"`
	Password  string    `json:"-"`
	Token     string    `gorm:"size:128" json:"token"`
	Verified  bool      `json:"verified"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PasswordReset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"index;size:255" json:"email"`
	Token     string    `gorm:"size:128" json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PolicyRule struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"uniqueIndex;size:128" json:"name"`
	Resource string `gorm:"size:64" json:"resource"`
	Action   string `gorm:"size:64" json:"action"`
}

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:64" json:"name"`
	DisplayName string    `gorm:"size:255" json:"displayName"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RoleRule struct {
	ID   uint `gorm:"primaryKey" json:"id"`
	RoleID uint `gorm:"index" json:"roleId"`
	RuleID uint `gorm:"index" json:"ruleId"`
}

type UserRole struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"index" json:"userId"`
	RoleID uint `gorm:"index" json:"roleId"`
}

type Theme struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:128" json:"name"`
	DisplayName string    `gorm:"size:255" json:"displayName"`
	Activated   bool      `json:"activated"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Plugin struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:128" json:"name"`
	DisplayName string    `gorm:"size:255" json:"displayName"`
	Path        string    `gorm:"size:512" json:"path"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Attachment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Filename  string    `gorm:"size:255" json:"filename"`
	Path      string    `gorm:"size:512" json:"path"`
	URL       string    `gorm:"size:512" json:"url"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Backup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Filename  string    `gorm:"size:255" json:"filename"`
	Status    string    `gorm:"size:32" json:"status"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Setting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:255" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
