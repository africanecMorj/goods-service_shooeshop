package domain

type Permission string

const (
	ProductRead   Permission = "product:read"
	ProductCreate Permission = "product:create"
	ProductDelete Permission = "product:delete"
	ProductUpdate Permission = "product:update"
	UserRead 	  Permission = "user:read"
	UserUpdate 	  Permission = "user:update"
	UserDelete    Permission = "user:delete" 
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
	"super-admin": {
		ProductRead,
		ProductCreate,
		ProductDelete,
		ProductUpdate,
		UserRead,
		UserUpdate,
		UserDelete,
	},
}