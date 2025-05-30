package controller

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type DBScopes struct {
	PageSize int
	Page     int
	Offset   int
	UserID   string
	Global   bool
}

func NewDBScopes(request *http.Request) *DBScopes {
	return &DBScopes{
		PageSize: GetPageSize(request),
		Page:     GetPage(request),
		Offset:   GetOffset(request),
		UserID:   GetOwnerID(request),
		Global:   IsGlobal(request),
	}
}

func (dbs *DBScopes) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(dbs.Offset).Limit(dbs.PageSize)
	}
}

func (dbs *DBScopes) Owned() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if dbs.Global {
			return db
		} else {
			return db.Where("user_id = ?", dbs.UserID)
		}
	}
}

// IsGlobal Check if current request is global
func IsGlobal(request *http.Request) bool {
	return (request.Context().Value(GLOBAL_SCOPE) != nil)
}

// GetOwnerID returns the current request user ID
func GetOwnerID(request *http.Request) string {
	return request.Context().Value(CURRENT_USER_ID).(string)
}

func GetPageSize(request *http.Request) int {
	query := request.URL.Query()
	pageSize, _ := strconv.Atoi(query.Get("page_size")) // Error is ignored because wrong or missing parameters are handled as 0
	switch {
	case pageSize > 500:
		pageSize = 500
	case pageSize <= 0:
		pageSize = 10
	}
	return pageSize
}

func GetPage(request *http.Request) int {
	query := request.URL.Query()
	page, _ := strconv.Atoi(query.Get("page")) // Error is ignored because wrong or missing parameters are handled as 0
	if page <= 0 {
		page = 1
	}
	return page
}

func GetOffset(request *http.Request) int {
	return (GetPage(request) - 1) * GetPageSize(request)
}
