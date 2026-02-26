package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
)

type MembershipHandler struct {
	membershipService service.MembershipService
	logger            logger.Logger
}

func NewMembershipHandler(membershipService service.MembershipService, logger logger.Logger) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService, logger: logger}
}

// CreateMembership godoc
// @Summary Create a membership
// @Tags memberships
// @Accept json
// @Produce json
// @Param request body dto.CreateMembershipRequest true "Membership data"
// @Success 201 {object} dto.MembershipResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /memberships [post]
func (h *MembershipHandler) CreateMembership(c *gin.Context) {
	var req dto.CreateMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	m, err := h.membershipService.CreateMembership(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, m)
}

// ListMembershipsByUnit godoc
// @Summary List memberships by unit
// @Tags memberships
// @Accept json
// @Produce json
// @Param unit_id query string true "Unit ID"
// @Success 200 {array} dto.MembershipResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /memberships [get]
func (h *MembershipHandler) ListMembershipsByUnit(c *gin.Context) {
	unitID := c.Query("unit_id")
	if unitID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "unit_id query parameter is required", Code: "INVALID_REQUEST"})
		return
	}
	var filters sharedrepo.ListFilters
	if search := c.Query("search"); search != "" {
		filters.Search = search
		if fields := c.Query("search_fields"); fields != "" {
			filters.SearchFields = strings.Split(fields, ",")
		}
	}
	memberships, err := h.membershipService.ListMembershipsByUnit(c.Request.Context(), unitID, filters)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, memberships)
}

// ListMembershipsByRole godoc
// @Summary List memberships by role
// @Tags memberships
// @Accept json
// @Produce json
// @Param unit_id query string true "Unit ID"
// @Param role query string true "Role name"
// @Success 200 {array} dto.MembershipResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /memberships/by-role [get]
func (h *MembershipHandler) ListMembershipsByRole(c *gin.Context) {
	unitID := c.Query("unit_id")
	role := c.Query("role")
	if unitID == "" || role == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "unit_id and role query parameters are required", Code: "INVALID_REQUEST"})
		return
	}
	var filters sharedrepo.ListFilters
	if search := c.Query("search"); search != "" {
		filters.Search = search
		if fields := c.Query("search_fields"); fields != "" {
			filters.SearchFields = strings.Split(fields, ",")
		}
	}
	memberships, err := h.membershipService.ListMembershipsByRole(c.Request.Context(), unitID, role, filters)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, memberships)
}

// GetMembership godoc
// @Summary Get a membership by ID
// @Tags memberships
// @Accept json
// @Produce json
// @Param id path string true "Membership ID (UUID)"
// @Success 200 {object} dto.MembershipResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /memberships/{id} [get]
func (h *MembershipHandler) GetMembership(c *gin.Context) {
	id := c.Param("id")
	m, err := h.membershipService.GetMembership(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, m)
}

// UpdateMembership godoc
// @Summary Update a membership
// @Tags memberships
// @Accept json
// @Produce json
// @Param id path string true "Membership ID (UUID)"
// @Param request body dto.UpdateMembershipRequest true "Membership update data"
// @Success 200 {object} dto.MembershipResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /memberships/{id} [put]
func (h *MembershipHandler) UpdateMembership(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}
	m, err := h.membershipService.UpdateMembership(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, m)
}

// DeleteMembership godoc
// @Summary Delete a membership
// @Tags memberships
// @Accept json
// @Produce json
// @Param id path string true "Membership ID (UUID)"
// @Success 204 "No content"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /memberships/{id} [delete]
func (h *MembershipHandler) DeleteMembership(c *gin.Context) {
	id := c.Param("id")
	if err := h.membershipService.DeleteMembership(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ExpireMembership godoc
// @Summary Expire a membership
// @Tags memberships
// @Accept json
// @Produce json
// @Param id path string true "Membership ID (UUID)"
// @Success 200 {object} dto.MembershipResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /memberships/{id}/expire [post]
func (h *MembershipHandler) ExpireMembership(c *gin.Context) {
	id := c.Param("id")
	m, err := h.membershipService.ExpireMembership(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, m)
}

// ListMembershipsByUser godoc
// @Summary List memberships for a user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path string true "User ID (UUID)"
// @Success 200 {array} dto.MembershipResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/{user_id}/memberships [get]
func (h *MembershipHandler) ListMembershipsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	var filters sharedrepo.ListFilters
	if search := c.Query("search"); search != "" {
		filters.Search = search
		if fields := c.Query("search_fields"); fields != "" {
			filters.SearchFields = strings.Split(fields, ",")
		}
	}
	memberships, err := h.membershipService.ListMembershipsByUser(c.Request.Context(), userID, filters)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, memberships)
}
