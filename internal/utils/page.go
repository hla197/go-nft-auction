package utils

import "gorm.io/gorm"

// Pagination 分页请求参数
type Pagination struct {
	Page     int `json:"page" form:"page" example:"1"`            // 页码
	PageSize int `json:"page_size" form:"page_size" example:"10"` // 每页数量
}

// PageResult 分页响应结果
type PageResult struct {
	Data       interface{} `json:"data"`        // 查询的数据列表
	Total      int64       `json:"total"`       // 总记录数
	Page       int         `json:"page"`        // 当前页
	PageSize   int         `json:"page_size"`   // 每页数量
	TotalPages int         `json:"total_pages"` // 总页数
	HasMore    bool        `json:"has_more"`    // 是否有下一页
}

// 初始化默认值
func (p *Pagination) setDefault() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	} else if p.PageSize > 100 { // 防止恶意请求导致内存溢出
		p.PageSize = 100
	}
}

// Paginate 返回一个 GORM Scopes，用于处理分页
func Paginate(p *Pagination) func(db *gorm.DB) *gorm.DB {
	// 设置默认值
	p.setDefault()

	// 计算偏移量
	offset := (p.Page - 1) * p.PageSize

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(p.PageSize)
	}
}

func GetPaginatedData(db *gorm.DB, dest interface{}, pagination *Pagination) (*PageResult, error) {
	// 如果为空，设置默认
	if pagination == nil {
		pagination = &Pagination{
			Page:     1,
			PageSize: 10,
		}
	}

	var total int64

	// 1. 先获取总记录数
	// Count 会忽略 Limit 和 Offset，所以可以正确获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 2. 应用分页 Scopes 并查询数据
	// 注意：这里需要传入一个新构建的 DB 实例，或者克隆原 DB，避免污染原条件
	// 使用 SubQuery 或者在 Count 后重建链式调用是安全的
	err := db.Scopes(Paginate(pagination)).Find(dest).Error
	if err != nil {
		return nil, err
	}

	// 3. 计算总页数
	totalPages := int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))

	// 4. 构造响应
	PageResult := &PageResult{
		Data:       dest,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
		HasMore:    int64(pagination.Page) < (total/int64(pagination.PageSize) + 1),
	}

	return PageResult, nil
}
