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

type AssetResponse struct {
	ID                 int    `json:"id"`
	IP                 string `json:"ip"`
	ApplicationSystem  string `json:"application_system"`
	ApplicationManager string `json:"application_manager"`
	OverallManager     string `json:"overall_manager"`
	IsVirtualMachine   bool   `json:"is_virtual_machine"`
	ResourcePool       string `json:"resource_pool"`
	DataCenter         string `json:"data_center"`
	RackLocation       string `json:"rack_location"`
	SNNumber           string `json:"sn_number"`
	OutOfBandIP        string `json:"out_of_band_ip"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

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

	// 将数据转换为前端友好的格式
	var assetResponses []AssetResponse
	for _, asset := range assets {
		assetResponses = append(assetResponses, AssetResponse{
			ID:                 asset.ID,
			IP:                 toString(asset.IP),
			ApplicationSystem:  toString(asset.ApplicationSystem),
			ApplicationManager: toString(asset.ApplicationManager),
			OverallManager:     toString(asset.OverallManager),
			IsVirtualMachine:   asset.IsVirtualMachine,
			ResourcePool:       toString(asset.ResourcePool),
			DataCenter:         toString(asset.DataCenter),
			RackLocation:       toString(asset.RackLocation),
			SNNumber:           toString(asset.SNNumber),
			OutOfBandIP:        toString(asset.OutOfBandIP),
			CreatedAt:          toString(asset.CreatedAt),
			UpdatedAt:          toString(asset.UpdatedAt),
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, assetResponses)
}

// toString 将 sql.NullString 转换为普通字符串
func toString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func UpdateAsset(w http.ResponseWriter, r *http.Request) {
	var asset models.Asset

	// 打印请求体以调试数据格式问题
	body, _ := io.ReadAll(r.Body)
	fmt.Println("Received data:", string(body))

	if err := json.Unmarshal(body, &asset); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

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
