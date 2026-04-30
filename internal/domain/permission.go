package domain

type Permission string

const (
	ProductRead   Permission = "product:read"
	ProductCreate Permission = "product:create"
	ProductDelete Permission = "product:delete"
	ProductUpdate Permission = "product:update"
)

var RolePermissions = map[string][]Permission{
	"user": {
		ProductRead,
	},
	"admin": {
		ProductRead,
		ProductCreate,
		ProductDelete,
		ProductUpdate,
	},
}