package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/paiduay/queq-hospital-api/classes"
)

// Hospital : models for db & json
type Hospital struct {
	HospitalID             uint32   `gorm:"primary_key;AUTO_INCREMENT" json:"hospital_id"`
	HospitalName           *string  `json:"hospital_name" binding:"required"`
	HospitalName2          *string  `json:"hospital_name2,omitempty"`
	HospitalName3          *string  `json:"hospital_name3,omitempty"`
	HospitalName4          *string  `json:"hospital_name4,omitempty"`
	HospitalName5          *string  `json:"hospital_name5,omitempty"`
	HospitalName6          *string  `json:"hospital_name6,omitempty"`
	HospitalName7          *string  `json:"hospital_name7,omitempty"`
	HospitalName8          *string  `json:"hospital_name8,omitempty"`
	HospitalUID            *string  `json:"hospital_uid,omitempty"`
	HospitalLogoPath       *string  `json:"hospital_logo_path,omitempty"`
	HospitalPrintLogoPath  *string  `json:"hospital_print_logo_path,omitempty"`
	HospitalMobileLogoPath *string  `json:"hospital_mobile_logo_path,omitempty"`
	Status                 uint32   `json:"status,omitempty" binding:"required"`
	OrderNo                *uint32  `json:"order_no,omitempty"`
	CreatedDate            string   `gorm:"default:CURRENT_TIMESTAMP" json:"created_date,omitempty"`
	UpdatedDate            string   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_date,omitempty"`
	ApProvince             *string  `json:"ap_province,omitempty"`
	ApLatitude             *float64 `json:"ap_latitude,omitempty"`
	ApLongitude            *float64 `json:"ap_longitude,omitempty"`
	ApDistanceFlag         *uint32  `json:"ap_distance_flag,omitempty"`
	ApDistanceLimit        *uint32  `json:"ap_distance_limit,omitempty"`
	ApOpenCloseRules       *string  `json:"ap_open_close_rules,omitempty"`
	ApBoxTypeCode          *string  `json:"ap_box_type_code,omitempty"`

	StationAmount *uint32 `gorm:"-" json:"station_amount,omitempty"`
	QueueAmount   *uint32 `gorm:"-" json:"queues_amount,omitempty"`
	Hours         *[]hour `gorm:"-" json:"hours,omitempty"`
}

type hour struct {
	Hour   string `json:"hour,omitempty"`
	Queues uint32 `json:"queues,omitempty"`
}

// GetHospitalList : get hospitals list
func GetHospitalList(hospitals *[]Hospital) error {
	rows, err := db.
		Raw(`
			SELECT H.hospital_id, H.hospital_name, H.created_date, H.updated_date, H.status, H.hospital_logo_path, H.hospital_print_logo_path, H.order_no, COUNT(S.station_id) as station_amount
			FROM Hospital as H
				LEFT JOIN Station as S ON H.hospital_id=S.hospital_id
			GROUP BY H.hospital_id, H.hospital_name, H.created_date, H.updated_date, H.status, H.hospital_logo_path, H.hospital_print_logo_path, H.order_no
			ORDER BY H.order_nohospital asc
		`).
		Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var tmp Hospital
		rows.Scan(&tmp.HospitalID, &tmp.HospitalName, &tmp.CreatedDate, &tmp.UpdatedDate, &tmp.Status, &tmp.HospitalLogoPath, &tmp.HospitalPrintLogoPath, &tmp.OrderNo, &tmp.StationAmount)
		*hospitals = append(*hospitals, tmp)
	}
	return nil
}

// GetAllRoomInHospital : get all rooms in hospital
func GetAllRoomInHospital(hID string, rooms *[]Room) {
	db.
		Table("Room JOIN Station ON Room.station_id = Station.station_id").
		Where("Station.hospital_id=?", hID).
		Select("*").
		Find(&rooms)
}

// GetStationsListInHospital : get all
func GetStationsListInHospital(hID string, stations *[]Station) {
	db.
		Table("Station as S LEFT JOIN Room as R ON S.station_id=R.station_id").
		Where("S.hospital_id=?", hID).
		Select("S.station_id, S.station_name, S.created_date, S.updated_date, S.status, COUNT(R.room_id) as room_amount, S.order_no, ROW_NUMBER() OVER(ORDER BY S.order_no asc) AS row_index").
		Group("S.station_id, S.station_name, S.created_date, S.updated_date, S.status, S.order_no").
		Find(&stations)
}

// GetQueueAmountInHospital : get queue amount in hospital by date
func GetQueueAmountInHospital(hospital *Hospital, fDate, tDate string) {
	row := db.Raw(`
		SELECT H.hospital_name, count(queue_id) as queues_amount
		FROM Queue as Q, Hospital as H, Station as S
		WHERE Q.station_id = S.station_id
			and S.hospital_id = H.hospital_id
			and H.hospital_id = ?
			and create_date between ? and ?
			and Q.status NOT IN (10, 11)
		GROUP BY H.hospital_name
	`, hospital.HospitalID, fDate, tDate).Row()
	row.Scan(&hospital.HospitalName, &hospital.QueueAmount)
}

// GetQueueDuringTheDay : get queue during the day
func GetQueueDuringTheDay(hospital *Hospital, fDate, tDate string) error {
	rows, err := db.Raw(`
		SELECT SUBSTRING(create_time, 1, 2) as Hour, COUNT(queue_id) as Queues
		FROM dbo.Hospital as H, dbo.Station as S, dbo.Queue as Q
		WHERE Q.station_id = S.station_id
			and S.hospital_id = H.hospital_id
			and H.hospital_id = ?
			and create_date between ? and ?
			and Q.status NOT IN (10, 11)
		GROUP BY SUBSTRING(create_time, 1, 2)
		ORDER BY SUBSTRING(create_time, 1, 2)
	`, hospital.HospitalID, fDate, tDate).Rows()
	if err != nil {
		return err
	}
	var tmpHourArr []hour
	for rows.Next() {
		var tmpHour hour
		rows.Scan(&tmpHour.Hour, &tmpHour.Queues)
		tmpHourArr = append(tmpHourArr, tmpHour)
	}
	hospital.Hours = &tmpHourArr
	return nil
}

// GetWeeklyQueueInHos : get weekly queue amount
func GetWeeklyQueueInHos (days *[]map[string]interface{}, year, weeknum, hID string) error {
	rows, err := db.Raw(`
		SELECT create_date, count(queue_id) as Queues
    FROM Queue as Q, Hospital as H, Station as S
    WHERE Q.station_id=S.station_id and S.hospital_id=H.hospital_id and H.hospital_id=?
        and SUBSTRING(create_date, 1, 4)=? and DATEPART(wk, create_date) = ?
    GROUP BY create_date
	`, hID, year, weeknum).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var date string
		var amount int
		tmp := make(map[string]interface{})
		rows.Scan(&date, &amount)
		tmp[date] = amount
		*days = append(*days, tmp)
	}
	return nil
}

// GetMonthlyQueueInHos : get mountly queue amount
func GetMonthlyQueueInHos (m *[]map[string]interface{}, lm *map[string]interface{}, hID, year string) error {
	now := time.Now()
	rows, err := db.Raw(`
		SELECT SUBSTRING(create_date, 6, 2) as month, count(queue_id) as Queues
    FROM Queue as Q, dbo.Hospital as H, dbo.Station as S
    WHERE Q.station_id=S.station_id and S.hospital_id=H.hospital_id and H.hospital_id=? and SUBSTRING(create_date, 1, 4)=?
    GROUP BY SUBSTRING(create_date, 6, 2)
	`, hID, year).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var month string
		var amount int
		tmp := make(map[string]interface{})
		rows.Scan(&month, &amount)
		tmp[month] = amount
		*m = append(*m, tmp)
	}
	row := db.Raw(`
		SELECT count(queue_id) as Queues
    FROM Queue as Q, dbo.Hospital as H, dbo.Station as S
    WHERE Q.station_id=S.station_id and S.hospital_id=H.hospital_id and H.hospital_id=? and SUBSTRING(create_date, 1, 7)=?
    GROUP BY SUBSTRING(create_date, 1, 7)
	`, hID, now.Format("2006-01")).Row()
	var queues int
	row.Scan(&queues)
	(*lm)["date"] = now.Format("2006-01")
	(*lm)["queues"] = queues
	return nil
}

// GetYearlyQueueInHos : get yearly queue amount
func GetYearlyQueueInHos (years *[]map[string]interface{}, startyear, endyear, hID string) error {
	rows, err := db.Raw(`
		SELECT SUBSTRING(create_date, 1, 4) as year, count(queue_id) as Queues
    FROM Queue as Q, Hospital as H, Station as S
    WHERE Q.station_id=S.station_id and S.hospital_id=H.hospital_id and H.hospital_id=? and (SUBSTRING(create_date, 1, 4) BETWEEN ?  AND ?)
    GROUP BY SUBSTRING(create_date, 1, 4)
	`, hID, startyear, endyear).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var date string
		var amount int
		tmp := make(map[string]interface{})
		rows.Scan(&date, &amount)
		tmp[date] = amount
		*years = append(*years, tmp)
	}
	return nil
}

// GetHospitalQueueInDateRange : get queue in hosptal filter by date range
func GetHospitalQueueInDateRange(hID, fromDate, toDate string, data *classes.Week) error {
	rows, err := db.Raw(`
		DECLARE @toDate date, @fromDate date, @hospitalId int
        SET @hospitalId = ?
        SET @fromDate = ?
        SET @toDate = ?

        IF @fromDate != @toDate
			SELECT create_date, count(queue_id) as Queues
			FROM Hospital as H, Station as S, Queue as Q
			WHERE
            Q.station_id=S.station_id
            and S.hospital_id=H.hospital_id
            and H.hospital_id=@hospitalId
            and create_date between @fromDate and @toDate
            and Q.status != 10
            and Q.status != 11
			GROUP BY create_date
		ELSE
			SELECT create_date, count(queue_id) as Queues
			FROM Hospital as H, Station as S, Queue as Q
			WHERE Q.station_id=S.station_id
            and S.hospital_id=H.hospital_id
            and H.hospital_id=@hospitalId
            and DATEPART(week, create_date) = DATEPART(week, @fromDate)
            and Q.status != 10
            and Q.status != 11
			GROUP BY create_date
	`, hID, fromDate, toDate).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var date string
		var amount int32
		rows.Scan(&date, &amount)
		t, err:= time.Parse("2006-01-02", date)
		if err != nil {
			return err
		}

		data.SetDayInWeek(int(t.Weekday()), amount)
	}
	return nil
}

// CheckIfHospitalNameExist : check if name exist in table hospital
func CheckIfHospitalNameExist(name string) bool {
	var tmp Hospital
	db.
		Table("Hospital").
		Where("hospital_name = ?", name).
		Scan(&tmp)

	if tmp.HospitalName == nil {
		return true
	}
	return false
}

// CreateNewHospital :
func CreateNewHospital(hospital *Hospital) error {
	tx := db.Begin()
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	row := db.Raw("SELECT MAX(order_no)+1 as next_no FROM dbo.Hospital").Row()
	row.Scan(&hospital.OrderNo)
	uuID := fmt.Sprintf("%v", id)
	hospital.HospitalUID = &uuID
	if err := tx.Debug().Table("Hospital").Create(&hospital).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
