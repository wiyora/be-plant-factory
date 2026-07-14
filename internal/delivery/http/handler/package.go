package handler

import "github.com/samber/do/v2"

var Package = do.Package(
	do.Lazy(NewHealthHandler),
	do.Lazy(NewAuthHandler),
	do.Lazy(NewUserMeHandler),
	do.Lazy(NewUserHandler),
	do.Lazy(NewRoleHandler),
	do.Lazy(NewTenantHandler),
	do.Lazy(NewUserTenantHandler),
	do.Lazy(NewPermissionHandler),
	do.Lazy(NewStorageHandler),
)
