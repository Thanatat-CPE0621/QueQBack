package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/paiduay/queq-hospital-api/classes"
	"gitlab.com/paiduay/queq-hospital-api/worker"
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

// ReorderHospitalList :
type ReorderHospitalList struct {
	ID      uint32 `json:"id" binding:"required"`
	OrderNO uint32 `json:"order_no" binding:"required"`
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

// GetInformation : get hospital information
func (h *Hospital) GetInformation(hID string) {
	db.
		Table("Hospital").
		Where("hospital_id=?", hID).
		Select("*").
		Scan(&h)
}

// GetStationinHoswithInfo : get stations in hospital with information
func (model *Station) GetStationinHoswithInfo (fd, td, hID string) (s []Station, err error) {
	s = make([]Station, 0)
	rows, err := db.Raw(`
		DECLARE @toDate date, @fromDate date, @hospitalId int
		SET @fromDate = ?
		SET @toDate = ?
		SET @hospitalId = ?;
		WITH db_to_query_clone AS (
			 SELECT S.station_id as s_id, Q.create_date, Q.create_time, Q.call_time, DATEDIFF(minute, CONVERT(Varchar(10), Q.create_date, 112) + ' ' + CONVERT(Varchar(8), Q.create_time),CONVERT(Varchar(10), Q.create_date, 112) + ' ' + CONVERT(Varchar(8), Q.call_time)) as queueing_time, Q.queue_id, Q.queue_number_text
			 FROM dbo.Hospital as H
					 JOIN Station as S ON H.hospital_id = S.hospital_id
					 JOIN Queue as Q ON S.station_id = Q.station_id
			 WHERE Q.create_date between @fromDate and @toDate
					 and H.hospital_id=@hospitalId
					 and Q.status=0
		 AND Q.status!=91
					 and Q.status!=10
		 and Q.status!=11
	 ), db_to_query AS (
			 SELECT
				 S.station_id as s_id,
				 Q.create_date,
				 Q.create_time,
				 Q.call_time,
				 Q.wait_min as 'queueing_time',
				 Q.queue_id,
				 Q.queue_number_text
			 FROM dbo.Hospital as H
				 JOIN Station as S ON H.hospital_id = S.hospital_id
				 JOIN Queue as Q ON S.station_id = Q.station_id
			 WHERE Q.create_date between @fromDate and @toDate
				 and H.hospital_id=@hospitalId
				 and Q.status=0
				 and Q.status!=91
				 and Q.status!=10
				 and Q.status!=11
	 ), stations AS (
		 SELECT S.station_id as station_id, station_name, count(room_id) as rooms, S.stat_gray, S.stat_green, S.stat_yellow, S.stat_red
		 FROM Hospital as H, Station as S LEFT JOIN Room as R ON R.station_id=S.station_id
		 WHERE S.hospital_id=H.hospital_id and H.hospital_id=@hospitalId
		 GROUP BY S.station_id, station_name, S.stat_gray, S.stat_green, S.stat_yellow, S.stat_red
	 )
	 SELECT
			 stations.station_id,
			 station_name,
			 rooms,
			 (
					 SELECT avg(Q.wait_min)
					 FROM Queue as Q
					 WHERE Q.create_date between @fromDate and @toDate
						 and Q.station_id= stations.station_id
						 AND Q.status!=91
			 ) as avg_time , sub1.queues 'queues', (
					 SELECT
					 COUNT(Q.queue_id) as queues_all
					 FROM Queue as Q
					 WHERE
						 Q.create_date between  @fromDate and  @toDate
						 and Q.station_id = stations.station_id
						 and Q.status!=10
						 and Q.status!=11
			 )
			 as 'queues_all',
			 sub2.max_q_time 'max_queueing_time',
			 sub2.max_q_number 'max_queueing_number',
			 sub2.min_q_time 'min_queueing_time',
			 sub2.min_q_number 'min_queueing_number',
			 stat_gray, stat_green, stat_yellow, stat_red
	 FROM stations LEFT JOIN
			 (
					 SELECT
					 t1.s_id as station_id,
					 AVG(t1.queueing_time) as 'avg',
					 COUNT(t1.queue_id) as 'queues'
									 FROM db_to_query_clone as t1
									 GROUP BY t1.s_id
			 ) sub1 ON (stations.station_id=sub1.station_id)

			 LEFT JOIN
			 (
					 SELECT
					 t2.s_id as station_id,
					 AVG(t2.queueing_time) as 'avg',
					 max_table.queueing_time as 'max_q_time',
					 max_table.queue_number_text as 'max_q_number',
					 min_table.queueing_time as 'min_q_time',
					 min_table.queue_number_text as 'min_q_number'
			 FROM db_to_query as t2
					 LEFT JOIN
					 (
									 SELECT tmp.*, ROW_NUMBER() OVER(PARTITION BY tmp.s_id ORDER BY ISNULL(tmp.queueing_time, 0) DESC, tmp.queueing_time DESC) as rank
					 FROM db_to_query as tmp
							 ) as max_table ON t2.s_id = max_table.s_id,
					 (
									 SELECT tmp.*, ROW_NUMBER() OVER(PARTITION BY tmp.s_id ORDER BY ISNULL(tmp.queueing_time, 999999) ASC, tmp.queueing_time ASC) as rank
					 FROM db_to_query as tmp
							 ) as min_table
			 WHERE t2.s_id=max_table.s_id and max_table.rank = 1 and t2.s_id=min_table.s_id and min_table.rank = 1
			 GROUP BY
			 t2.s_id,
			 max_table.queueing_time,
			 max_table.queue_number_text,
			 min_table.queueing_time,
			 min_table.queue_number_text
			 ) sub2 ON (stations.station_id=sub2.station_id)
	`, fd, td, hID).Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var tmp Station
		rows.Scan(
			&tmp.StationID, &tmp.StationName, &tmp.RoomAmount, &tmp.AvgTime, &tmp.QueueAmount,
			&tmp.QueueAll, &tmp.MaxQTime, &tmp.MaxQNumber, &tmp.MinQTime, &tmp.MinQNumber,
			&tmp.StatGray, &tmp.StatGreen, &tmp.StatYellow, &tmp.StatRed,
		)
		s = append(s, tmp)
	}
	return
}

// EditHospitalInformation : edit hospital information
func (h *Hospital) EditHospitalInformation(hID string) error {
	tx := db.Begin()
	now := time.Now()
	if err := tx.Exec(`
		UPDATE Hospital
			SET hospital_name=?, status=?, hospital_logo_path=?,
			hospital_print_logo_path=?, hospital_mobile_logo_path=?, updated_date=?
		WHERE hospital_id=?
	`, h.HospitalName, h.Status, h.HospitalLogoPath, h.HospitalPrintLogoPath,
	h.HospitalMobileLogoPath, now.Format("2006-01-02 15:04:05"), hID).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
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
func GetWeeklyQueueInHos(days *[]map[string]interface{}, year, weeknum, hID string) error {
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
func GetMonthlyQueueInHos(m *[]map[string]interface{}, lm *map[string]interface{}, hID, year string) error {
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
func GetYearlyQueueInHos(years *[]map[string]interface{}, startyear, endyear, hID string) error {
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
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return err
		}

		data.SetDayInWeek(int(t.Weekday()), amount)
	}
	return nil
}

// GetAvgAllDay : get all average queuing time in day
func GetAvgAllDay(hID, date string, data *map[string]interface{}) error {
	rows, err := db.Raw(`
		SELECT SUBSTRING(create_time, 1, 2) as Hour, AVG(DATEDIFF(mi, call_time, commit_time)) AS queueing_time
        FROM Queue as Q, dbo.Hospital as H, dbo.Station as S
        WHERE Q.station_id=S.station_id and S.hospital_id=H.hospital_id and H.hospital_id=? and create_date=?
		GROUP BY SUBSTRING(create_time, 1, 2)
	`, hID, date).Rows()
	if err != nil {
		return err
	}
	tmp := make(map[string]interface{})
	for rows.Next() {
		var hour string
		var queues int32
		rows.Scan(&hour, &queues)
		tmp[hour] = queues
	}
	*data = tmp

	return nil
}

// GetAvgAllDayWeek : get all average queuing time in week
func GetAvgAllDayWeek(hID, date string, data *map[string]interface{}) error {
	rows, err := db.Raw(`
		declare @hid int = ?
    declare @date datetime = ?
    declare @wn int = (datediff(week, dateadd(week, datediff(week, 0, dateadd(month, datediff(month, 0, @date), 0)), 0), @date) + 1)
    SELECT SUBSTRING(create_time, 1, 2) as Hour, AVG(DATEDIFF(mi, call_time, commit_time)) AS queueing_time
    FROM Queue as Q, dbo.Hospital as H, dbo.Station as S
    WHERE Q.station_id=S.station_id and S.hospital_id=H.hospital_id and H.hospital_id=@hid and SUBSTRING(create_date, 1, 7)=SUBSTRING(CONVERT(varchar, @date, 21), 1, 7)
		AND (datediff(week, dateadd(week, datediff(week, 0, dateadd(month, datediff(month, 0, create_date), 0)), 0), create_date) + 1)=@wn
		GROUP BY SUBSTRING(create_time, 1, 2)
	`, hID, date).Rows()
	if err != nil {
		return err
	}
	tmp := make(map[string]interface{})
	for rows.Next() {
		var hour string
		var queues int32
		rows.Scan(&hour, &queues)
		tmp[hour] = queues
	}
	*data = tmp

	return nil
}

// GetAvgAllDayMonth : get all average queuing time in month
func GetAvgAllDayMonth(hID, date string, data *map[string]interface{}) error {
	rows, err := db.Raw(`
		declare @hid int = ?
    declare @date datetime = ?
    declare @wn int = (datediff(week, dateadd(week, datediff(week, 0, dateadd(month, datediff(month, 0, @date), 0)), 0), @date) + 1)
    SELECT SUBSTRING(create_time, 1, 2) as Hour, AVG(DATEDIFF(mi, call_time, commit_time)) AS queueing_time
    FROM Queue as Q, dbo.Hospital as H, dbo.Station as S
    WHERE Q.station_id=S.station_id and S.hospital_id=H.hospital_id and H.hospital_id=@hid and SUBSTRING(create_date, 1, 7)=SUBSTRING(CONVERT(varchar, @date, 21), 1, 7)
		GROUP BY SUBSTRING(create_time, 1, 2)
	`, hID, date).Rows()
	if err != nil {
		return err
	}
	tmp := make(map[string]interface{})
	for rows.Next() {
		var hour string
		var queues int32
		rows.Scan(&hour, &queues)
		tmp[hour] = queues
	}
	*data = tmp

	return nil
}

// GetAverageDayQueueingTime : get average day queuing time
func GetAverageDayQueueingTime(hID, fdate, tdate string, data *map[string]interface{}) error {
	row := db.Raw(`
		SELECT AVG(Q.wait_min)
    FROM Queue as Q, Hospital as H, Station as S
    WHERE
	    Q.station_id=S.station_id
	    AND S.hospital_id=H.hospital_id
	    AND H.hospital_id=?
	    AND create_date between ? AND ?
     	AND Q.status!=91
	`, hID, fdate, tdate).Row()
	tmp := make(map[string]interface{})
	var tmpAvg int32
	row.Scan(&tmpAvg)
	tmp["time"] = tmpAvg
	*data = tmp
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

// ReorderHospitals ...
func ReorderHospitals(hospital []ReorderHospitalList) error {
	tx := db.Begin()
	c := len(hospital)
	jobs := make(chan func(), c)
	err := make(chan error, c)
	pool := worker.Pool{
		Amount: 3,
		Jobs:   jobs,
		Done:   err,
	}
	for _, element := range hospital {
		jobs <- func() {
			if fail := tx.Exec("UPDATE Hospital SET order_no = ? WHERE hospital_id = ?", element.OrderNO, element.ID).Error; fail != nil {
				err <- fail
			}
			err <- nil
		}
	}
	worker.Worker(pool)
	for idx := 0; idx < c; idx++ {
		e := <-err
		if e != nil {
			tx.Rollback()
			return e
		}
	}
	tx.Commit()
	return nil
}
