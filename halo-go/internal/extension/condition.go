package extension

import "time"

// ConditionStatus 表示条件的状态值。
type ConditionStatus string

const (
	// ConditionTrue 表示条件为真（状态满足）。
	ConditionTrue ConditionStatus = "True"
	// ConditionFalse 表示条件为假（状态不满足）。
	ConditionFalse ConditionStatus = "False"
	// ConditionUnknown 表示条件状态未知。
	ConditionUnknown ConditionStatus = "Unknown"
)

// Condition 表示扩展对象状态中的单个条件项。
// 用于表达资源的异步状态信息，如"已就绪"、"可用"、"异常"等语义化状态。
// 类似 Kubernetes 的 Condition 概念，控制器通过更新 Condition 来上报运行时状态。
type Condition struct {
	// Type 是条件的唯一标识符（如 "Ready"、"Available"），
	// 在同一对象的 Conditions 列表中必须唯一。
	Type string `json:"type"`

	// Status 是条件的当前状态值：True / False / Unknown。
	Status ConditionStatus `json:"status"`

	// Reason 是导致当前状态的机器可读原因码（单字词，大写开头），
	// 如 "Created"、"Unavailable"、"Progressing" 等。
	Reason string `json:"reason,omitempty"`

	// Message 是人类可读的状态描述信息，用于解释当前状态的原因和上下文。
	Message string `json:"message,omitempty"`

	// LastTransitionTime 是该条件最近一次发生状态变更的时间戳。
	LastTransitionTime time.Time `json:"lastTransitionTime"`
}

// IsTrue 判断条件状态是否为 True。
func (c *Condition) IsTrue() bool {
	return c.Status == ConditionTrue
}

// IsFalse 判断条件状态是否为 False。
func (c *Condition) IsFalse() bool {
	return c.Status == ConditionFalse
}

// IsUnknown 判断条件状态是否为 Unknown。
func (c *Condition) IsUnknown() bool {
	return c.Status == ConditionUnknown
}

// Equal 判断两个 Condition 是否在类型和状态上完全一致。
func (c *Condition) Equal(other *Condition) bool {
	if c == nil && other == nil {
		return true
	}
	if c == nil || other == nil {
		return false
	}
	return c.Type == other.Type && c.Status == other.Status
}

// NewCondition 创建并返回一个新的 Condition 实例，自动设置 LastTransitionTime 为当前时间。
func NewCondition(condType string, status ConditionStatus, reason, message string) *Condition {
	return &Condition{
		Type:               condType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: time.Now(),
	}
}

// ConditionList 是 Condition 的集合管理器，提供对多个条件的增删改查操作。
// 扩展对象应在其 Status 中嵌入此结构体以管理自身的状态条件列表。
type ConditionList struct {
	// items 是底层条件切片。
	items []*Condition
}

// NewConditionList 创建并返回一个空的 ConditionList 实例。
func NewConditionList() *ConditionList {
	return &ConditionList{
		items: make([]*Condition, 0),
	}
}

// Get 根据条件类型获取对应的 Condition。如果不存在则返回 nil。
func (l *ConditionList) Get(condType string) *Condition {
	for _, c := range l.items {
		if c.Type == condType {
			return c
		}
	}
	return nil
}

// Set 向条件列表中添加或更新指定类型的条件。
// 当条件状态发生变化时，自动更新 LastTransitionTime；
// 当条件状态未变时，保留原有的 LastTransitionTime 和 Message/Reason（除非显式覆盖）。
func (l *ConditionList) Set(condition *Condition) {
	for i, existing := range l.items {
		if existing.Type == condition.Type {
			if existing.Status != condition.Status {
				condition.LastTransitionTime = time.Now()
			} else if condition.LastTransitionTime.IsZero() {
				condition.LastTransitionTime = existing.LastTransitionTime
			}
			l.items[i] = condition
			return
		}
	}
	if condition.LastTransitionTime.IsZero() {
		condition.LastTransitionTime = time.Now()
	}
	l.items = append(l.items, condition)
}

// Remove 从条件列表中移除指定类型的条件。如果不存在则无操作。
func (l *ConditionList) Remove(condType string) {
	result := make([]*Condition, 0, len(l.items))
	for _, c := range l.items {
		if c.Type != condType {
			result = append(result, c)
		}
	}
	l.items = result
}

// Clear 清空所有条件。
func (l *ConditionList) Clear() {
	l.items = l.items[:0]
}

// Items 返回所有条件的只读切片副本。
func (l *ConditionList) Items() []*Condition {
	result := make([]*Condition, len(l.items))
	copy(result, l.items)
	return result
}

// Len 返回条件列表的长度。
func (l *ConditionList) Len() int {
	return len(l.items)
}

// IsEmpty 判断条件列表是否为空。
func (l *ConditionList) IsEmpty() bool {
	return len(l.items) == 0

}

// ForEach 遍历条件列表中的每个条件，调用给定的回调函数。
func (l *ConditionList) ForEach(fn func(*Condition)) {
	for _, c := range l.items {
		fn(c)
	}
}

// AnyTrue 判断是否存在至少一个状态为 True 的条件。
func (l *ConditionList) AnyTrue() bool {
	for _, c := range l.items {
		if c.IsTrue() {
			return true
		}
	}
	return false
}

// AllTrue 判断所有条件的状态是否均为 True。空列表返回 false。
func (l *ConditionList) AllTrue() bool {
	if len(l.items) == 0 {
		return false
	}
	for _, c := range l.items {
		if !c.IsTrue() {
			return false
		}
	}
	return true
}
