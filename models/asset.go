package models

import (
	"cmdb-backend/config"
	"database/sql"
	"log"
)

type Asset struct {
	ID                 int
	IP                 sql.NullString
	ApplicationSystem  sql.NullString
	ApplicationManager sql.NullString
	OverallManager     sql.NullString
	IsVirtualMachine   bool
	ResourcePool       sql.NullString
	DataCenter         sql.NullString
	RackLocation       sql.NullString
	SNNumber           sql.NullString
	OutOfBandIP        sql.NullString
	CreatedAt          sql.NullString
	UpdatedAt          sql.NullString
}

func (asset *Asset) Create() error {
	query := `INSERT INTO cmdb_assets (ip, application_system, application_manager, overall_manager, is_virtual_machine, resource_pool, data_center, rack_location, sn_number, out_of_band_ip)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := config.DB.Exec(query, asset.IP, asset.ApplicationSystem, asset.ApplicationManager, asset.OverallManager, asset.IsVirtualMachine, asset.ResourcePool, asset.DataCenter, asset.RackLocation, asset.SNNumber, asset.OutOfBandIP)
	if err != nil {
		log.Println("Error creating asset:", err)
		return err
	}
	return nil
}

func GetAllAssets() ([]Asset, error) {
	query := "SELECT * FROM cmdb_assets"
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		err := rows.Scan(&asset.ID, &asset.IP, &asset.ApplicationSystem, &asset.ApplicationManager, &asset.OverallManager, &asset.IsVirtualMachine, &asset.ResourcePool, &asset.DataCenter, &asset.RackLocation, &asset.SNNumber, &asset.OutOfBandIP, &asset.CreatedAt, &asset.UpdatedAt)
		if err != nil {
			log.Println("Error scanning asset:", err)
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (asset *Asset) Update() error {
	query := `UPDATE cmdb_assets SET ip=?, application_system=?, application_manager=?, overall_manager=?, is_virtual_machine=?, resource_pool=?, data_center=?, rack_location=?, sn_number=?, out_of_band_ip=? WHERE id=?`

	_, err := config.DB.Exec(query, asset.IP, asset.ApplicationSystem, asset.ApplicationManager, asset.OverallManager, asset.IsVirtualMachine, asset.ResourcePool, asset.DataCenter, asset.RackLocation, asset.SNNumber, asset.OutOfBandIP, asset.ID)
	if err != nil {
		log.Println("Error updating asset:", err)
		return err
	}
	return nil
}

func DeleteAsset(id int) error {
	query := "DELETE FROM cmdb_assets WHERE id=?"
	_, err := config.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting asset:", err)
		return err
	}
	return nil
}

func GetAssetsByQuery(ip, sn string) ([]Asset, error) {
	query := "SELECT * FROM cmdb_assets WHERE 1=1"
	args := []interface{}{}

	if ip != "" {
		query += " AND ip = ?"
		args = append(args, ip)
	}

	if sn != "" {
		query += " AND sn_number = ?"
		args = append(args, sn)
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var asset Asset
		err := rows.Scan(&asset.ID, &asset.IP, &asset.ApplicationSystem, &asset.ApplicationManager, &asset.OverallManager, &asset.IsVirtualMachine, &asset.ResourcePool, &asset.DataCenter, &asset.RackLocation, &asset.SNNumber, &asset.OutOfBandIP, &asset.CreatedAt, &asset.UpdatedAt)
		if err != nil {
			log.Println("Error scanning asset:", err)
			return nil, err
		}
		assets = append(assets, asset)
	}

	// 确保返回的始终是一个数组，即使没有匹配的记录
	if assets == nil {
		return []Asset{}, nil
	}

	return assets, nil
}
