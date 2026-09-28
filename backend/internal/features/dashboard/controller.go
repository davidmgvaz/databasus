package dashboard

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	users_enums "databasus-backend/internal/features/users/enums"
	users_middleware "databasus-backend/internal/features/users/middleware"
)

type DashboardController struct {
	dashboardService *DashboardService
	logger           *slog.Logger
}

func (c *DashboardController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard", c.GetWorkspaceDashboard)
	router.GET("/dashboard/storages", c.GetWorkspaceStorages)

	adminOnly := router.Group("/dashboard")
	adminOnly.Use(users_middleware.RequireRole(users_enums.UserRoleAdmin))

	adminOnly.GET("/installation", c.GetInstallationDashboard)
	adminOnly.GET("/installation/storages", c.GetInstallationStorages)
}

// GetWorkspaceDashboard
// @Summary Get workspace dashboard
// @Description Lists every database of a workspace with its health, its last healthcheck attempts and the totals of the backups it currently stores, plus the workspace totals. Sizes count completed backups only; for physical databases they also include WAL segments. The backup count includes every stored backup regardless of status, and for physical databases counts FULL and INCREMENTAL backups but not WAL segments. Any workspace member can read it.
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Param workspace_id query string true "Workspace ID"
// @Success 200 {object} WorkspaceDashboard
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dashboard [get]
func (c *DashboardController) GetWorkspaceDashboard(ctx *gin.Context) {
	user, ok := users_middleware.GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workspaceIDStr := ctx.Query("workspace_id")
	if workspaceIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id query parameter is required"})
		return
	}

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace_id"})
		return
	}

	dashboard, err := c.dashboardService.GetWorkspaceDashboard(ctx.Request.Context(), user, workspaceID)
	if errors.Is(err, ErrDashboardUnavailable) {
		c.logger.ErrorContext(ctx.Request.Context(), "failed to build workspace dashboard", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrDashboardUnavailable.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dashboard)
}

// GetInstallationDashboard
// @Summary Get installation dashboard (ADMIN only)
// @Description Returns the database count and backup totals across every workspace, including databases created by restores that belong to no workspace. Sizes and counts follow the same rules as the workspace dashboard.
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} DashboardTotals
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dashboard/installation [get]
func (c *DashboardController) GetInstallationDashboard(ctx *gin.Context) {
	totals, err := c.dashboardService.GetInstallationDashboard()
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "failed to build installation dashboard", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrDashboardUnavailable.Error()})
		return
	}

	ctx.JSON(http.StatusOK, totals)
}

// GetWorkspaceStorages
// @Summary Get workspace storages with their space
// @Description Lists every storage of a workspace with the number of its databases that use it, the size of their stored backups, and the storage's total, used and free space when the storage can report it. Local, NAS, SFTP, rclone remotes that support it and Google Drive accounts with a quota report space; S3, Azure Blob and FTP do not. A storage whose probe fails is returned with status ERROR instead of failing the request. freeSpaceBytes sums the free space of the reporting storages, counting the local disk once. Any workspace member can read it.
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Param workspace_id query string true "Workspace ID"
// @Success 200 {object} WorkspaceStorages
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dashboard/storages [get]
func (c *DashboardController) GetWorkspaceStorages(ctx *gin.Context) {
	user, ok := users_middleware.GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workspaceIDStr := ctx.Query("workspace_id")
	if workspaceIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "workspace_id query parameter is required"})
		return
	}

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace_id"})
		return
	}

	workspaceStorages, err := c.dashboardService.GetWorkspaceStorages(ctx.Request.Context(), user, workspaceID)
	if errors.Is(err, ErrDashboardUnavailable) {
		c.logger.ErrorContext(ctx.Request.Context(), "failed to build workspace storages", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrDashboardUnavailable.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, workspaceStorages)
}

// GetInstallationStorages
// @Summary Get installation free space (ADMIN only)
// @Description Returns the free space summed across every storage of every workspace that can report it, counting the local disk once.
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} InstallationStorages
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dashboard/installation/storages [get]
func (c *DashboardController) GetInstallationStorages(ctx *gin.Context) {
	installationStorages, err := c.dashboardService.GetInstallationStorages(ctx.Request.Context())
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "failed to build installation storages", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ErrDashboardUnavailable.Error()})
		return
	}

	ctx.JSON(http.StatusOK, installationStorages)
}
