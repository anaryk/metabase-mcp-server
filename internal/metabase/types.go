package metabase

import "time"

// Card represents a Metabase saved question/card.
type Card struct {
	ID                    int              `json:"id,omitempty"`
	Name                  string           `json:"name,omitempty"`
	Description           *string          `json:"description,omitempty"`
	Display               string           `json:"display,omitempty"`
	DatasetQuery          map[string]any   `json:"dataset_query,omitempty"`
	VisualizationSettings map[string]any   `json:"visualization_settings,omitempty"`
	CollectionID          *int             `json:"collection_id,omitempty"`
	Archived              *bool            `json:"archived,omitempty"`
	EnableEmbedding       *bool            `json:"enable_embedding,omitempty"`
	EmbeddingParams       map[string]any   `json:"embedding_params,omitempty"`
	DatabaseID            *int             `json:"database_id,omitempty"`
	TableID               *int             `json:"table_id,omitempty"`
	QueryType             *string          `json:"query_type,omitempty"`
	CreatorID             *int             `json:"creator_id,omitempty"`
	CreatedAt             *time.Time       `json:"created_at,omitempty"`
	UpdatedAt             *time.Time       `json:"updated_at,omitempty"`
	ResultMetadata        []map[string]any `json:"result_metadata,omitempty"`
}

// Dashboard represents a Metabase dashboard.
type Dashboard struct {
	ID                    int              `json:"id,omitempty"`
	Name                  string           `json:"name,omitempty"`
	Description           *string          `json:"description,omitempty"`
	CollectionID          *int             `json:"collection_id,omitempty"`
	Parameters            []map[string]any `json:"parameters,omitempty"`
	Archived              *bool            `json:"archived,omitempty"`
	EnableEmbedding       *bool            `json:"enable_embedding,omitempty"`
	EmbeddingParams       map[string]any   `json:"embedding_params,omitempty"`
	DashCards             []DashCard       `json:"dashcards,omitempty"`
	Tabs                  []map[string]any `json:"tabs,omitempty"`
	CreatorID             *int             `json:"creator_id,omitempty"`
	CreatedAt             *time.Time       `json:"created_at,omitempty"`
	UpdatedAt             *time.Time       `json:"updated_at,omitempty"`
	CacheTTL              *int             `json:"cache_ttl,omitempty"`
	VisualizationSettings map[string]any   `json:"visualization_settings,omitempty"`
}

// DashCard represents a card placed on a dashboard.
type DashCard struct {
	ID                    int              `json:"id,omitempty"`
	DashboardID           int              `json:"dashboard_id,omitempty"`
	CardID                *int             `json:"card_id,omitempty"`
	Row                   int              `json:"row"`
	Col                   int              `json:"col"`
	SizeX                 int              `json:"size_x,omitempty"`
	SizeY                 int              `json:"size_y,omitempty"`
	Series                []map[string]any `json:"series,omitempty"`
	ParameterMappings     []map[string]any `json:"parameter_mappings,omitempty"`
	VisualizationSettings map[string]any   `json:"visualization_settings,omitempty"`
}

// Collection represents a Metabase collection.
type Collection struct {
	ID              any     `json:"id,omitempty"` // can be int or "root"
	Name            string  `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	ParentID        *int    `json:"parent_id,omitempty"`
	Color           *string `json:"color,omitempty"`
	Archived        *bool   `json:"archived,omitempty"`
	Namespace       *string `json:"namespace,omitempty"`
	PersonalOwnerID *int    `json:"personal_owner_id,omitempty"`
}

// CollectionItem represents an item within a collection.
type CollectionItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Model       string  `json:"model"`
}

// Database represents a Metabase database connection.
type Database struct {
	ID                int            `json:"id,omitempty"`
	Name              string         `json:"name,omitempty"`
	Engine            string         `json:"engine,omitempty"`
	Details           map[string]any `json:"details,omitempty"`
	Tables            []Table        `json:"tables,omitempty"`
	Features          []string       `json:"features,omitempty"`
	NativePermissions string         `json:"native_permissions,omitempty"`
	CreatedAt         *time.Time     `json:"created_at,omitempty"`
	UpdatedAt         *time.Time     `json:"updated_at,omitempty"`
}

// Table represents a database table in Metabase.
type Table struct {
	ID          int     `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	Schema      *string `json:"schema,omitempty"`
	DBID        int     `json:"db_id,omitempty"`
	EntityType  *string `json:"entity_type,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
	Rows        *int    `json:"rows,omitempty"`
}

// Field represents a field/column in a table.
type Field struct {
	ID              int     `json:"id,omitempty"`
	Name            string  `json:"name,omitempty"`
	DisplayName     *string `json:"display_name,omitempty"`
	Description     *string `json:"description,omitempty"`
	DatabaseType    string  `json:"database_type,omitempty"`
	BaseType        string  `json:"base_type,omitempty"`
	SemanticType    *string `json:"semantic_type,omitempty"`
	TableID         int     `json:"table_id,omitempty"`
	FKTargetFieldID *int    `json:"fk_target_field_id,omitempty"`
	Visibility      string  `json:"visibility_type,omitempty"`
}

// FieldValues represents the distinct values of a field.
type FieldValues struct {
	FieldID int     `json:"field_id"`
	Values  [][]any `json:"values"`
}

// ForeignKey represents a foreign key relationship.
type ForeignKey struct {
	Origin       *FKField `json:"origin,omitempty"`
	Destination  *FKField `json:"destination,omitempty"`
	Relationship string   `json:"relationship,omitempty"`
}

// FKField represents one side of a foreign key.
type FKField struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	TableID int    `json:"table_id,omitempty"`
	Table   *Table `json:"table,omitempty"`
}

// DatasetQueryRequest represents a query execution request.
type DatasetQueryRequest struct {
	Database int            `json:"database"`
	Type     string         `json:"type"`
	Native   *NativeQuery   `json:"native,omitempty"`
	Query    map[string]any `json:"query,omitempty"`
}

// NativeQuery represents a native SQL query.
type NativeQuery struct {
	Query        string         `json:"query"`
	TemplateTags map[string]any `json:"template-tags,omitempty"`
}

// DatasetQueryResponse represents query results.
type DatasetQueryResponse struct {
	Data       DatasetData `json:"data"`
	Status     string      `json:"status"`
	Context    string      `json:"context,omitempty"`
	RowCount   int         `json:"row_count"`
	DatabaseID int         `json:"database_id,omitempty"`
	Error      *string     `json:"error,omitempty"`
}

// DatasetData holds the actual query result data.
type DatasetData struct {
	Rows       [][]any      `json:"rows"`
	Cols       []DatasetCol `json:"cols"`
	NativeForm *NativeForm  `json:"native_form,omitempty"`
}

// DatasetCol represents a column in query results.
type DatasetCol struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	BaseType    string `json:"base_type"`
	FieldRef    []any  `json:"field_ref,omitempty"`
}

// NativeForm holds the generated SQL for a query.
type NativeForm struct {
	Query  string `json:"query,omitempty"`
	Params []any  `json:"params,omitempty"`
}

// User represents a Metabase user.
type User struct {
	ID          int        `json:"id,omitempty"`
	Email       string     `json:"email,omitempty"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	CommonName  *string    `json:"common_name,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
	IsSuperuser *bool      `json:"is_superuser,omitempty"`
	GroupIDs    []int      `json:"group_ids,omitempty"`
	LastLogin   *time.Time `json:"last_login,omitempty"`
	CreatedAt   *time.Time `json:"date_joined,omitempty"`
}

// PermissionGroup represents a Metabase permission group.
type PermissionGroup struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Members []User `json:"members,omitempty"`
}

// SearchResult represents a search result item.
type SearchResult struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	Model        string  `json:"model"`
	CollectionID *int    `json:"collection_id,omitempty"`
	TableID      *int    `json:"table_id,omitempty"`
	DatabaseID   *int    `json:"database_id,omitempty"`
}

// SearchResponse wraps search results with pagination info.
type SearchResponse struct {
	Data   []SearchResult `json:"data"`
	Total  int            `json:"total"`
	Limit  *int           `json:"limit,omitempty"`
	Offset *int           `json:"offset,omitempty"`
}

// Alert represents a Metabase alert.
type Alert struct {
	ID             int              `json:"id,omitempty"`
	CardID         int              `json:"card_id,omitempty"`
	Card           *Card            `json:"card,omitempty"`
	Channels       []map[string]any `json:"channels,omitempty"`
	AlertCondition string           `json:"alert_condition,omitempty"`
	AlertAboveGoal *bool            `json:"alert_above_goal,omitempty"`
	AlertFirstOnly bool             `json:"alert_first_only,omitempty"`
	Creator        *User            `json:"creator,omitempty"`
	CreatedAt      *time.Time       `json:"created_at,omitempty"`
}

// Action represents a Metabase action.
type Action struct {
	ID          int              `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	ModelID     int              `json:"model_id,omitempty"`
	Type        string           `json:"type,omitempty"`
	Parameters  []map[string]any `json:"parameters,omitempty"`
	CreatedAt   *time.Time       `json:"created_at,omitempty"`
	UpdatedAt   *time.Time       `json:"updated_at,omitempty"`
}

// Timeline represents a Metabase timeline.
type Timeline struct {
	ID           int             `json:"id,omitempty"`
	Name         string          `json:"name,omitempty"`
	Description  *string         `json:"description,omitempty"`
	CollectionID *int            `json:"collection_id,omitempty"`
	Icon         *string         `json:"icon,omitempty"`
	Default      *bool           `json:"default,omitempty"`
	Events       []TimelineEvent `json:"events,omitempty"`
	Archived     *bool           `json:"archived,omitempty"`
	CreatedAt    *time.Time      `json:"created_at,omitempty"`
}

// TimelineEvent represents an event on a timeline.
type TimelineEvent struct {
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Timestamp   string     `json:"timestamp,omitempty"`
	TimeZone    string     `json:"time_zone,omitempty"`
	Icon        *string    `json:"icon,omitempty"`
	Archived    *bool      `json:"archived,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

// ActivityItem represents an activity log entry.
type ActivityItem struct {
	ID        int            `json:"id,omitempty"`
	Topic     string         `json:"topic,omitempty"`
	UserID    *int           `json:"user_id,omitempty"`
	Model     *string        `json:"model,omitempty"`
	ModelID   *int           `json:"model_id,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
	Timestamp *time.Time     `json:"timestamp,omitempty"`
}

// RecentItem represents a recently viewed item.
type RecentItem struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Model   string `json:"model"`
	ModelID int    `json:"model_id"`
}

// Setting represents a Metabase setting.
type Setting struct {
	Key         string `json:"key"`
	Value       any    `json:"value"`
	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`
}
