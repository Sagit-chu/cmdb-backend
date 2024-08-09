package controllers

import (
	"cmdb-backend/models"
	"cmdb-backend/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx/v3"
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
	queryParams := r.URL.Query()
	ip := queryParams.Get("ip")
	sn := queryParams.Get("sn")

	assets, err := models.GetAssetsByQuery(ip, sn)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve assets")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, assets)
}

func UpdateAsset(w http.ResponseWriter, r *http.Request) {
	var asset models.Asset
	if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}
	asset.ID = id

	if err := asset.Update(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update asset")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, asset)
}

func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	if err := models.DeleteAsset(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete asset")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func ImportAssets(w http.ResponseWriter, r *http.Request) {
	// 处理文件上传
	file, _, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error retrieving the file")
		return
	}
	defer file.Close()

	// 解析 Excel 文件
	var assets []models.Asset
	if err := parseExcelFile(file, &assets); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error parsing the file")
		return
	}

	// 插入数据到数据库
	for _, asset := range assets {
		if err := asset.Create(); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error inserting data into database")
			return
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Data imported successfully"})
}

func parseExcelFile(file multipart.File, assets *[]models.Asset) error {
	// 读取文件内容到内存中
	buf, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// 使用 OpenBinary 方法解析文件内容
	xlFile, err := xlsx.OpenBinary(buf)
	if err != nil {
		fmt.Println("Error opening Excel file:", err)
		return err
	}

	// 检查是否存在工作表
	if len(xlFile.Sheets) == 0 {
		return errors.New("no sheets found in Excel file")
	}

	// 读取数据
	for _, sheet := range xlFile.Sheets {
		err := sheet.ForEachRow(func(row *xlsx.Row) error {
			var asset models.Asset

			// 将每个单元格的值转换为 sql.NullString
			asset.IP = toNullString(row.GetCell(0).String())
			asset.ApplicationSystem = toNullString(row.GetCell(1).String())
			asset.ApplicationManager = toNullString(row.GetCell(2).String())
			asset.OverallManager = toNullString(row.GetCell(3).String())
			asset.IsVirtualMachine = row.GetCell(4).Bool() // Bool 类型不需要转换
			asset.ResourcePool = toNullString(row.GetCell(5).String())
			asset.DataCenter = toNullString(row.GetCell(6).String())
			asset.RackLocation = toNullString(row.GetCell(7).String())
			asset.SNNumber = toNullString(row.GetCell(8).String())
			asset.OutOfBandIP = toNullString(row.GetCell(9).String())

			*assets = append(*assets, asset)
			return nil
		})

		if err != nil {
			fmt.Println("Error reading rows:", err)
			return err
		}
	}

	return nil
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
