package repository

// ListFilters represents common filters for listing entities
type ListFilters struct {
	IsActive *bool
	Limit    int
	Offset   int
}

// ScreenTemplateFilter filters for listing templates
type ScreenTemplateFilter struct {
	Pattern string
	Offset  int
	Limit   int
}

// ScreenInstanceFilter filters for listing instances
type ScreenInstanceFilter struct {
	TemplateID *string
	Offset     int
	Limit      int
}
