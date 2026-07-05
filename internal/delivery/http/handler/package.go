package handler

import "github.com/samber/do/v2"

var Package = do.Package(
	do.Lazy(NewHealthHandler),
	do.Lazy(NewAuthHandler),
	do.Lazy(NewUserMeHandler),
	do.Lazy(NewUserHandler),
)
