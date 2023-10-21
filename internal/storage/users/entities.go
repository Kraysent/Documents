package users

const (
	TableName      = "documents.t_user"
	ColumnID       = "id"
	ColumnUsername = "username"
)

type GetUserRequest struct {
	Username string
}

type GetUserResult struct {
	UserID int64 `db:"user_id"`
}

type CreateUserRequest struct {
	Username string
}

type CreateUserResult struct {
	UserID int64 `db:"user_id"`
}
