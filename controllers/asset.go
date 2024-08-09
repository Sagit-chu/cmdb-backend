package controllers

import (
	"cmdb-backend/models"
	"cmdb-backend/utils"
	"encoding/json"
	"net/http"
)

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var asset models.Asset
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := asset.Create(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create asset")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, asset)
}

func GetAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := models.GetAllAssets()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve assets")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, assets)
}
