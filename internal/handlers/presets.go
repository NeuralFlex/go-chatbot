package handlers

import (
	"github.com/gin-gonic/gin"

	"go-chatbot/internal/models"
	"go-chatbot/internal/services"
)

// PresetsHandler godoc
// @Summary      List the available preset analysis prompts
// @Tags         presets
// @Produce      json
// @Success      200  {array}  models.Preset
// @Router       /presets [get]
func PresetsHandler(c *gin.Context) {
	c.JSON(200, []models.Preset{
		{ID: "cost_savings", Label: "Cost savings", Text: services.PresetCostSavings},
		{ID: "meeting_notes", Label: "Meeting notes", Text: services.PresetMeetingNotes},
		{ID: "client_proposals", Label: "Client proposals", Text: services.PresetClientProposals},
	})
}
