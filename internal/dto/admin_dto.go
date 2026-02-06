package dto

// DashboardResponse represents the dashboard summary for the back office
type DashboardResponse struct {
	UsersCount       int64 `json:"users_count"`
	CurriculumsCount int64 `json:"curriculums_count"`
}

// UsersStatsResponse represents user statistics
type UsersStatsResponse struct {
	Total int64 `json:"total"`
}

// CurriculumsStatsResponse represents curriculum statistics
type CurriculumsStatsResponse struct {
	Total int64 `json:"total"`
}

// AdminUsersListResponse represents paginated users list for admin
type AdminUsersListResponse struct {
	Data       []UserResponse   `json:"data"`
	Pagination CursorPagination `json:"pagination"`
}

// AdminCurriculumsListResponse represents paginated curriculums list for admin
type AdminCurriculumsListResponse struct {
	Data       []CurriculumResponse `json:"data"`
	Pagination CursorPagination     `json:"pagination"`
}
