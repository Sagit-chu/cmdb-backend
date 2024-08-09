package models

import (
	"cmdb-backend/config"
	"log"
)

type Asset struct {
	ID                 int
	IP                 string
	ApplicationSystem  string
	ApplicationManager string
	OverallManager     string
	IsVirtualMachine   bool
	ResourcePool       string
	DataCenter         string
	RackLocation       string
	SNNumber           string
	OutOfBandIP        string
	CreatedAt          string
	UpdatedAt          string
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
