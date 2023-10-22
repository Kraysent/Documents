package users

const (
	TableName      = "documents.t_user"
	ColumnID       = "id"
	ColumnUsername = "username"
	ColumnGoogleID = "google_id"
)

type GetUserRequest struct {
	Fields map[string]any
}

type GetUserResult struct {
	UserID int64 `db:"user_id"`
}

type CreateUserRequest struct {
	Username string
	GoogleID string
}

type CreateUserResult struct {
	UserID int64 `db:"user_id"`
}
