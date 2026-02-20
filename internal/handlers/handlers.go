package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/debt-optimization-engine/internal/models"
	"github.com/user/debt-optimization-engine/internal/repositories"
	"github.com/user/debt-optimization-engine/internal/services"
)

type Handler struct {
	repo              repositories.Repository
	settlementService *services.SettlementService
}

func NewHandler(repo repositories.Repository, ss *services.SettlementService) *Handler {
	return &Handler{repo: repo, settlementService: ss}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.CreateGroup(c.Request.Context(), &group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, group)
}

func (h *Handler) AddMember(c *gin.Context) {
	groupID := c.Param("id")
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.AddMemberToGroup(c.Request.Context(), groupID, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user added to group"})
}

func (h *Handler) CreateExpense(c *gin.Context) {
	groupID := c.Param("id")
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	gid, _ := models.ParseUUID(groupID) 
	expense.GroupID = gid

	// Logic for equal splits if only user IDs are provided without amounts
	if expense.SplitType == models.SplitEqual && len(expense.Splits) > 0 {
		userIDs := make([]string, len(expense.Splits))
		for i, s := range expense.Splits {
			userIDs[i] = s.UserID.String()
		}
		amounts := services.CalculateEqualSplits(expense.Amount, userIDs)
		for i := range expense.Splits {
			expense.Splits[i].Amount = amounts[i]
		}
	}

	if err := services.ValidateSplits(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateExpense(c.Request.Context(), &expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, expense)
}

func (h *Handler) GetSettlement(c *gin.Context) {
	groupID := c.Param("id")
	
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to *time.Time
	if fromStr != "" {
		if t, err := time.Parse("2006-01-02", fromStr); err == nil { from = &t }
	}
	if toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil { to = &t }
	}

	resp, err := h.settlementService.GetSettlement(c.Request.Context(), groupID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetBalances(c *gin.Context) {
	groupID := c.Param("id")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to *time.Time
	if fromStr != "" {
		if t, err := time.Parse("2006-01-02", fromStr); err == nil { from = &t }
	}
	if toStr != "" {
		if t, err := time.Parse("2006-01-02", toStr); err == nil { to = &t }
	}

	balances, err := h.settlementService.CalculateBalances(c.Request.Context(), groupID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

func (h *Handler) CompareStrategies(c *gin.Context) {
	groupID := c.Param("id")
	cmp, err := h.settlementService.CompareStrategies(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cmp)
}
