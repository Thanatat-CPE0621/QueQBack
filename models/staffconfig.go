package models

import (
  "fmt"

  "github.com/jinzhu/gorm"
)

// StaffConfig ...
type StaffConfig struct {
  ID                        uint32          `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
  StaffID    								uint32     			`json:"staff_id"`
  StaffType    							uint32     			`json:"staff_type"`
  StationID    							uint32     			`json:"station_id"`
  RoomID  									uint32     			`json:"room_id"`
  Status    								uint32     			`json:"status"`
  Online    								uint32     			`json:"online"`

  IsDelete                  bool            `gorm:"-" json:"is_del"`
  IsNew                     bool            `gorm:"-" json:"is_new"`
}


// CreateStaffConfigArray : create from staffconfig's array
func CreateStaffConfigArray (staffID uint32, configs []StaffConfig,  tx *gorm.DB) error {
  query := `
    INSERT INTO StaffConfig
      (staff_id, staff_type, station_id, room_id, status, online)
    VALUES
  `
  for _, item := range configs {
    query += fmt.Sprintf("(%v, %v, %v, %v, %v, %v),",
      staffID, item.StaffType, item.StationID, item.RoomID,
      item.Status, item.Online,
    )
  }
  if err := tx.Exec(query[:len(query) - 1]).Error; err != nil {
    return err
  }
  return nil
}


// UpdateStaffConfigArray : update from staffconfig's array
func UpdateStaffConfigArray (staff Staff, configs []StaffConfig,  tx *gorm.DB) error {
  insert := `
    INSERT INTO StaffConfig
      (staff_id, staff_type, station_id, room_id, status, online)
    VALUES
  `
  delete := `DELETE FROM StaffConfig WHERE id in (?)`
  var deleted []uint32
  inserted := 0
  for _, element := range configs {
    if element.IsDelete {
      deleted = append(deleted, element.ID)
    } else if element.IsNew {
      insert += fmt.Sprintf("(%v, %v, %v, %v, %v, %v),",
      staff.StaffID, staff.StaffType, element.StationID, element.RoomID,
      element.Status, element.Online)
      inserted++
    } else {
      tx.Table("StaffConfig").Where("id = ?", element.ID).Save(&element)
    }
  }
  if inserted > 0 {
    if err := tx.Exec(insert).Error; err != nil {
      return err
    }
  }
  if len(deleted) > 0 {
    if err := tx.Exec(delete, deleted).Error; err != nil {
      return err
    }
  }
  return nil
}
