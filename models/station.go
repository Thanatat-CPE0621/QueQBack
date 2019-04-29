package models

import (
  // "fmt"
)

// Station models
type Station struct {
  StationID                   uint32        `gorm:"primary_key;AUTO_INCREMENT" json:"station_id"`
  StationCode                 *string       `json:"station_code,omitempty"`
  StationName                 *string       `json:"station_name"`
  StationName2                *string       `json:"station_name2,omitempty"`
  StationName3                *string       `json:"station_name3,omitempty"`
  StationName4                *string       `json:"station_name4,omitempty"`
  StationName5                *string       `json:"station_name5,omitempty"`
  StationName6                *string       `json:"station_name6,omitempty"`
  StationName7                *string       `json:"station_name7,omitempty"`
  StationName8                *string       `json:"station_name8,omitempty"`
  QueuePrefix                 *string       `json:"queue_prefix,omitempty"`
  NoRoomQueuePrefix           *string       `json:"no_room_queue_prefix,omitempty"`
  AppointmentQueuePrefix      *string       `json:"appointment_queue_prefix,omitempty"`
  HospitalID                  uint32        `json:"hospital_id,omitempty"`
  QueueNumberType             *uint32       `json:"queue_number_type,omitempty"`
  QueueNumberIndex            *uint32       `json:"queue_number_index,omitempty" `
  QueueShowTime               *int          `json:"queue_show_time,omitempty"`
  ReasonText                  string        `json:"reason_text,omitempty"`
  StationMode                 int           `gorm:"default:0" json:"station_mode,omitempty"`
  Status                      *uint32       `gorm:"default:0" json:"status,omitempty"`
  OrderNo                     *uint32       `json:"order_no,omitempty"`
  CreatedDate                 *string       `json:"created_date,omitempty"`
  UpdatedDate                 *string       `json:"updated_date,omitempty"`
  StatGray                    *int          `json:"stat_gray,omitempty"`
  StatGreen                   *int          `json:"stat_green,omitempty"`
  StatYellow                  *int          `json:"stat_yellow,omitempty"`
  StatRed                     *int          `json:"stat_red,omitempty"`
  ApAllowCustQueue            *int          `json:"ap_allow_cust_queue,omitempty"`
  ApAllowCustAppoint          *int          `json:"ap_allow_cust_appoint,omitempty"`
  ApSlotQuota                 *int          `json:"ap_slot_quota,omitempty"`
  ApMinutesBeforeSubmit       *int          `json:"ap_minutes_before_submit,omitempty"`
  ApMinutesBeforeConfirm      *int          `json:"ap_minutes_before_confirm,omitempty"`
  ApRemarkMessage             *string       `json:"ap_remark_message,omitempty"`
}

// HighWaittime models for json
type HighWaittime struct {
  Number        *string     `gorm:"queue_number_text" json:"queue_number_text"`
  RoomID        *uint32     `gorm:"room_id" json:"room_id"`
  RoomNumber    *string     `gorm:"room_name" json:"room_name"`
  WaitingTime   *int32      `gorm:"queueing_time" json:"queueing_time"`
}

// CheckStationNameAvailable : check if stationame is available in hospital (return isAvailable)
func CheckStationNameAvailable (sName string, hID uint64) bool {
  var station Station
  db.Table("Station").Where("station_name=? AND hospital_id = ?", sName, hID).Find(&station)
  if (station.StationCode != nil) {
    return false
  }
  return true
}

// CheckStationCodeAvailable : check if stationcode is available (return isAvailable)
func CheckStationCodeAvailable (sCode string) bool {
  var station Station
  db.Table("Station").Where("station_code=?", sCode).Find(&station)
  if (station.StationName != nil) {
    return false
  }
  return true
}

// GetHighestWaitingTime : get top 1 highest waiting time queue
func GetHighestWaitingTime (date string, sID uint64, model *HighWaittime) {
  row := db.Raw(`
    WITH db_to_query AS (
        SELECT S.station_id as s_id, R.room_id, R.room_name, DATEDIFF(minute, CONVERT(Varchar(10), Q.create_date, 112) + ' ' + CONVERT(Varchar(8), Q.create_time), CONVERT(Varchar(10), Q.create_date, 112) + ' ' + CONVERT(Varchar(8), Q.call_time))  as queueing_time, Q.queue_id, Q.queue_number_text
        FROM dbo.Hospital as H, dbo.Station as S, dbo.Queue as Q, dbo.Room as R
        WHERE H.hospital_id=S.hospital_id and Q.station_id=S.station_id and Q.room_id=R.room_id and Q.create_date=? and S.station_id=? and Q.call_time is not null
    )
    SELECT TOP 1 room_id, room_name, queueing_time, queue_number_text
    FROM db_to_query
    ORDER BY queueing_time desc, queue_id desc
  `, date, sID).Row()
  row.Scan(&model.RoomID, &model.RoomNumber, &model.WaitingTime, &model.Number)
}

// CreateStation : get top 1 highest waiting time queue
func CreateStation (station *Station) error {
  row := db.Select("MAX(order_no)+1 as next_no").Table("Station").Row()
  row.Scan(&station.OrderNo)

  tx := db.Begin()
  if err := tx.Debug().Table("Station").Create(&station).Error; err != nil {
    tx.Rollback()
    return err
  }
  tx.Commit()
  return nil
}
