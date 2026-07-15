package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/handler"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/middleware"
	"github.com/samber/do/v2"
)

type Router interface {
	Register(app *fiber.App)
}

type router struct {
	mid        middleware.Middleware     `do:""`
	health     handler.HealthHandler     `do:""`
	auth       handler.AuthHandler       `do:""`
	userMe     handler.UserMeHandler     `do:""`
	user       handler.UserHandler       `do:""`
	role       handler.RoleHandler       `do:""`
	tenant     handler.TenantHandler     `do:""`
	userTenant handler.UserTenantHandler `do:""`
	permission handler.PermissionHandler `do:""`
	storage    handler.StorageHandler    `do:""`
}

func New(i do.Injector) (Router, error) {
	return do.InvokeStruct[*router](i)
}

func (r *router) Register(app *fiber.App) {
	app.Get("/health", r.health.Check)

	r.authRoute(app.Group("/auth"))
	r.userMeRoute(app.Group("/user/me"))
	r.userRoute(app.Group("/user"))
	r.roleRoute(app.Group("/role"))
	r.tenantRoute(app.Group("/tenant"))
	r.permissionRoute(app.Group("/permission"))
	r.storageRoute(app.Group("/storage"))
}

func (r *router) authRoute(route fiber.Router) {
	// TODO: add rate limiter for refresh token
	route.Post("/refresh", r.auth.RefreshToken)
	route.Get("/google", r.auth.GoogleLogin)
	route.Get("/google/callback", r.auth.GoogleCallback)

	route.Get("/me", r.mid.Auth(), r.auth.Me)
	route.Post("/logout", r.mid.Auth(), r.auth.Logout)
}

func (r *router) userRoute(route fiber.Router) {
	route.Get("/", r.mid.Auth(), r.user.List)
	route.Get("/dropdown", r.mid.Auth(), r.user.Dropdown)
	route.Get("/dropdown/selected", r.mid.Auth(), r.user.SelectedDropdown)
	route.Get("/:id", r.mid.Auth(), r.user.Detail)
	route.Get("/:id/tenants", r.mid.Auth(), r.userTenant.ListUserTenants)
	route.Post("/", r.mid.Auth(), r.user.Create)
	route.Put("/:id", r.mid.Auth(), r.user.Update)
	route.Post("/:id/status/:status", r.mid.Auth(), r.user.UpdateStatus)
}

func (r *router) userMeRoute(route fiber.Router) {
	route.Post("/getting-started", r.mid.Auth(), r.userMe.GettingStarted)
}

func (r *router) roleRoute(route fiber.Router) {
	route.Get("", r.mid.Auth(), r.role.List)
	route.Get("/dropdown", r.mid.Auth(), r.role.Dropdown)
	route.Get("/dropdown/selected", r.mid.Auth(), r.role.SelectedDropdown)
	route.Get("/:id", r.mid.Auth(), r.role.Detail)
	route.Post("", r.mid.Auth(), r.role.Create)
	route.Put("/:id", r.mid.Auth(), r.role.Update)
	route.Delete("/:id", r.mid.Auth(), r.role.Delete)
}

func (r *router) tenantRoute(route fiber.Router) {
	route.Get("/", r.mid.Auth(), r.tenant.List)
	route.Get("/dropdown", r.mid.Auth(), r.tenant.Dropdown)
	route.Get("/dropdown/selected", r.mid.Auth(), r.tenant.SelectedDropdown)
	route.Get("/:id", r.mid.Auth(), r.tenant.Detail)
	route.Get("/:id/users", r.mid.Auth(), r.userTenant.ListTenantUsers)
	route.Post("/", r.mid.Auth(), r.tenant.Create)
	route.Put("/:id", r.mid.Auth(), r.tenant.Update)
	route.Put("/:id/users/:user_id", r.mid.Auth(), r.userTenant.AssignUserRole)
	route.Delete("/:id/users/:user_id", r.mid.Auth(), r.userTenant.RemoveTenantUser)
	route.Post("/:id/status/:status", r.mid.Auth(), r.tenant.UpdateStatus)
}

func (r *router) permissionRoute(route fiber.Router) {
	route.Get("", r.mid.Auth(), r.permission.All)
}

func (r *router) storageRoute(route fiber.Router) {
	route.Post("/presigned-upload", r.mid.Auth(), r.storage.PresignedUpload)
}
