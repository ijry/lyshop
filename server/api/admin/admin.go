package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ijry/lyshop/config"
	"github.com/ijry/lyshop/core/middleware"
	"github.com/ijry/lyshop/core/plugin"
	"github.com/ijry/lyshop/core/response"
	adminsvc "github.com/ijry/lyshop/service/admin"
)

// RegisterRoutes registers admin/role management routes.
// All require system:admin permission.
func RegisterRoutes(g *gin.RouterGroup) {
	sys := g.Group("").Use(middleware.RequirePermission("system:admin"))

	// Admins
	sys.GET("/admins", listAdmins)
	sys.POST("/admins", createAdmin)
	sys.PUT("/admins/:id", updateAdmin)
	sys.DELETE("/admins/:id", deleteAdmin)

	// Roles
	sys.GET("/roles", listRoles)
	sys.POST("/roles", createRole)
	sys.PUT("/roles/:id", updateRole)
	sys.DELETE("/roles/:id", deleteRole)

	// All available permissions (from plugins)
	sys.GET("/permissions", listPermissions)
}

func listAdmins(c *gin.Context) {
	list, err := adminsvc.ListAdmins(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func createAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		RoleID   uint64 `json:"role_id"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	a, err := adminsvc.CreateAdmin(c.Request.Context(), req.Username, req.Password, req.RoleID)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, a)
}

func updateAdmin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates map[string]any
	c.ShouldBindJSON(&updates)
	if err := adminsvc.UpdateAdmin(c.Request.Context(), id, updates); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func deleteAdmin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := adminsvc.DeleteAdmin(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func listRoles(c *gin.Context) {
	list, err := adminsvc.ListRoles(c.Request.Context())
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, list)
}

func createRole(c *gin.Context) {
	var req struct {
		Name  string   `json:"name" binding:"required"`
		Perms []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	r, err := adminsvc.CreateRole(c.Request.Context(), req.Name, req.Perms)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, r)
}

func updateRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Name  string   `json:"name"`
		Perms []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	if err := adminsvc.UpdateRole(c.Request.Context(), id, req.Name, req.Perms); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func deleteRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := adminsvc.DeleteRole(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, nil)
}

func listPermissions(c *gin.Context) {
	perms := plugin.AllPermissions(config.Global.Plugins.Enabled)
	// Add system-level permissions
	perms = append([]string{"system:admin", "system:config"}, perms...)
	response.OK(c, perms)
}
