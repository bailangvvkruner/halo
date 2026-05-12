package extension

import "strings"

// ListOptions 封装了扩展资源列表查询的所有可选参数。
// 用于控制分页、排序、标签选择器和字段选择器等查询行为。
type ListOptions struct {
	// Page 表示请求的页码，默认为 1。
	Page int

	// Size 表示每页的记录条数，默认为 20。0 表示不分页返回全部结果。
	Size int

	// Sort 表示排序字段和方向（如 "creationTimestamp,desc"）。
	Sort string

	// LabelSelector 是基于标签的选择器表达式，
	// 支持等值匹配（key=value）、不等匹配（key!=value）、
	// 存在判断（key）、集合匹配（key in (v1,v2)）等语法。
	LabelSelector string

	// FieldSelector 是基于字段的选择器表达式，
	// 用于按元数据字段（如 name）进行精确筛选。
	FieldSelector string
}

// DefaultListOptions 返回默认的 ListOptions 配置（第一页，每页 20 条）。
func DefaultListOptions() *ListOptions {
	return &ListOptions{
		Page: 1,
		Size: 20,
	}
}

// Offset 根据 Page 和 Size 计算数据库查询偏移量。
func (o *ListOptions) Offset() int {
	if o.Page <= 0 {
		o.Page = 1
	}
	if o.Size <= 0 {
		return 0
	}
	return (o.Page - 1) * o.Size
}

// Limit 返回每页的最大记录数。当 Size 为 0 时表示不限制（返回全部）。
func (o *ListOptions) Limit() int {
	if o.Size <= 0 {
		return 0
	}
	return o.Size
}

// Validate 校验 ListOptions 的参数合法性，非法参数将被修正为默认值。
func (o *ListOptions) Validate() *ListOptions {
	if o == nil {
		return DefaultListOptions()
	}
	if o.Page < 1 {
		o.Page = 1
	}
	if o.Size < 0 {
		o.Size = 20
	}
	return o
}

// LabelSelectorParser 提供对 LabelSelector 表达式的解析能力。
// 返回一个可迭代的标签匹配规则列表。
type LabelSelectorParser struct {
	raw string
}

// NewLabelSelectorParser 创建一个标签选择器解析器。
func NewLabelSelectorParser(selector string) *LabelSelectorParser {
	return &LabelSelectorParser{raw: selector}
}

// Parse 将原始选择器字符串解析为多个独立的标签要求（requirement）。
// 每个要求之间以逗号分隔，表示 AND 关系。
//
// 支持的语法：
//   - key=value        等值匹配
//   - key==value       等值匹配（同上）
//   - key!=value       不等匹配
//   - key              存在性检查
//   - key in (v1,v2)   集合包含
//   - key notin (v1,v2) 集合排除
func (p *LabelSelectorParser) Parse() []Requirement {
	if p.raw == "" {
		return nil
	}
	parts := splitTrim(p.raw, ',')
	requirements := make([]Requirement, 0, len(parts))
	for _, part := range parts {
		req := parseRequirement(strings.TrimSpace(part))
		if req != nil {
			requirements = append(requirements, *req)
		}
	}
	return requirements
}

// Requirement 表示一个标签选择器的匹配条件。
type Requirement struct {
	Key    string
	Op     Operator
	Values []string
}

// Operator 定义标签选择器支持的运算符类型。
type Operator string

const (
	// OpEquals 表示等值匹配运算符（= 或 ==）
	OpEquals Operator = "="
	// OpNotEquals 表示不等匹配运算符（!=）
	OpNotEquals Operator = "!="
	// OpExists 表示存在性检查运算符
	OpExists Operator = "Exists"
	// OpIn 表示集合包含运算符（in）
	OpIn Operator = "In"
	// OpNotIn 表示集合排除运算符（notin）
	OpNotIn Operator = "NotIn"
)

// Matches 判断给定的标签字典是否满足此匹配条件。
func (r *Requirement) Matches(labels map[string]string) bool {
	switch r.Op {
	case OpExists:
		_, ok := labels[r.Key]
		return ok
	case OpEquals:
		val, ok := labels[r.Key]
		return ok && val == r.Values[0]
	case OpNotEquals:
		val, ok := labels[r.Key]
		return !ok || val != r.Values[0]
	case OpIn:
		val, ok := labels[r.Key]
		if !ok {
			return false
		}
		for _, v := range r.Values {
			if val == v {
				return true
			}
		}
		return false
	case OpNotIn:
		val, ok := labels[r.Key]
		if !ok {
			return true
		}
		for _, v := range r.Values {
			if val == v {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// parseRequirement 将单个表达式解析为 Requirement 结构体。
func parseRequirement(expr string) *Requirement {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil
	}

	for _, op := range []struct {
		op    Operator
		prefix string
	}{
		{OpNotIn, "notin"},
		{OpIn, "in"},
		{OpNotEquals, "!="},
		{OpEquals, "=="},
		{OpEquals, "="},
	} {
		idx := strings.Index(expr, string(op.prefix))
		if idx > 0 {
			key := strings.TrimSpace(expr[:idx])
			valueStr := strings.TrimSpace(expr[idx+len(op.prefix):])
			var values []string
			if op.op == OpIn || op.op == OpNotIn {
				values = parseValueSet(valueStr)
				if values == nil {
					continue
				}
			} else {
				values = []string{valueStr}
			}
			return &Requirement{
				Key:    key,
				Op:     op.op,
				Values: values,
			}
		}
	}

	return &Requirement{
		Key:    expr,
		Op:     OpExists,
		Values: nil,
	}
}

// parseValueSet 解析括号包裹的值集合（如 "(v1,v2)"），返回值的切片。
func parseValueSet(s string) []string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "(") || !strings.HasSuffix(s, ")") {
		return nil
	}
	inner := s[1 : len(s)-1]
	if inner == "" {
		return nil
	}
	return splitTrim(inner, ',')
}

// splitTrim 按指定分隔符切割字符串并去除每个元素的空白字符。
func splitTrim(s string, sep rune) []string {
	parts := strings.Split(s, string(sep))
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
