package models

import (
	"gitlab.com/paiduay/queq-hospital-api/worker"
)

// Station models
type Station struct {
	StationID              uint32  `gorm:"primary_key;AUTO_INCREMENT" json:"station_id"`
	StationCode            *string `json:"station_code,omitempty"`
	StationName            *string `json:"station_name"`
	StationName2           *string `json:"station_name2,omitempty"`
	StationName3           *string `json:"station_name3,omitempty"`
	StationName4           *string `json:"station_name4,omitempty"`
	StationName5           *string `json:"station_name5,omitempty"`
	StationName6           *string `json:"station_name6,omitempty"`
	StationName7           *string `json:"station_name7,omitempty"`
	StationName8           *string `json:"station_name8,omitempty"`
	QueuePrefix            *string `json:"queue_prefix,omitempty"`
	NoRoomQueuePrefix      *string `json:"no_room_queue_prefix,omitempty"`
	AppointmentQueuePrefix *string `json:"appointment_queue_prefix,omitempty"`
	HospitalID             uint32  `json:"hospital_id,omitempty"`
	QueueNumberType        *uint32 `json:"queue_number_type,omitempty"`
	QueueNumberIndex       *uint32 `json:"queue_number_index,omitempty" `
	QueueShowTime          *int    `json:"queue_show_time,omitempty"`
	ReasonText             string  `json:"reason_text,omitempty"`
	StationMode            int     `gorm:"default:0" json:"station_mode,omitempty"`
	Status                 *uint32 `gorm:"default:0" json:"status,omitempty"`
	OrderNo                *uint32 `json:"order_no,omitempty"`
	CreatedDate            *string `json:"created_date,omitempty"`
	UpdatedDate            *string `json:"updated_date,omitempty"`
	StatGray               *int    `json:"stat_gray,omitempty"`
	StatGreen              *int    `json:"stat_green,omitempty"`
	StatYellow             *int    `json:"stat_yellow,omitempty"`
	StatRed                *int    `json:"stat_red,omitempty"`
	ApAllowCustQueue       *int    `json:"ap_allow_cust_queue,omitempty"`
	ApAllowCustAppoint     *int    `json:"ap_allow_cust_appoint,omitempty"`
	ApSlotQuota            *int    `json:"ap_slot_quota,omitempty"`
	ApMinutesBeforeSubmit  *int    `json:"ap_minutes_before_submit,omitempty"`
	ApMinutesBeforeConfirm *int    `json:"ap_minutes_before_confirm,omitempty"`
	ApRemarkMessage        *string `json:"ap_remark_message,omitempty"`

	HospitalName *string `gorm:"-" json:"hospital_name,omitempty"`
	RoomAmount   *int    `gorm:"-" json:"room_amount,omitempty"`
	RowIndex     *int    `gorm:"-" json:"row_index,omitempty"`

	AvgTime     	*int    `gorm:"-" json:"avgQueueingTime,omitempty"`
	QueueAmount   *int    `gorm:"-" json:"queues,omitempty"`
	QueueAll     	*int    `gorm:"-" json:"queueAll,omitempty"`
	MaxQTime     	*int    `gorm:"-" json:"maxQueueingTime,omitempty"`
	MaxQNumber    *string	`gorm:"-" json:"maxQueueNumber,omitempty"`
	MinQTime     	*int    `gorm:"-" json:"minQueueingTime,omitempty"`
	MinQNumber    *string	`gorm:"-" json:"minQueueingNumber,omitempty"`
}

// HighWaittime models for json
type HighWaittime struct {
	Number      *string `gorm:"queue_number_text" json:"queue_number_text"`
	RoomID      *uint32 `gorm:"room_id" json:"room_id"`
	RoomNumber  *string `gorm:"room_name" json:"room_name"`
	WaitingTime *int32  `gorm:"queueing_time" json:"queueing_time"`
}

// MultiStation model for json
type MultiStation struct {
	StationID   uint32  `json:"station_id"`
	RowIndex    uint32  `json:"row_index"`
	StationName *string `json:"station_name"`
	RoomAmount  *int    `json:"room_amount"`
	StaffAmount *int    `json:"staff_amount"`
}

// ReorderStationModel ...
type ReorderStationModel struct {
	HospitalID string               `json:"hospitalID" binding:"required"`
	Data       []ReorderStationList `json:"data" binding:"required"`
}

// ReorderStationList :
type ReorderStationList struct {
	ID      uint32 `json:"id"`
	OrderNO uint32 `json:"order_no"`
}

// GetStationInfoByID : get station infomation by station_id
func GetStationInfoByID(sID string, station *Station) {
	db.
		Table("Station as S, Hospital as H").
		Where("H.hospital_id = S.hospital_id AND station_id = ?", sID).
		Select("S.*, H.hospital_name").
		Find(&station)
}

// GetRoomsInStation : get rooms in station
func GetRoomsInStation(sID string, rooms *[]Room) {
	db.Table("Room").Where("station_id=?", sID).Find(&rooms)
}

// GetRoomsInStationWithBrief : get rooms in station with queue infomation in room
func GetRoomsInStationWithBrief(sID string, fDate string, tDate string, rooms *[]Room) error {
	rows, err := db.Raw(`
	declare @fromdate date = ?;
	declare @todate date = ?;

	WITH rooms_in_station AS (
		SELECT S.station_id as sid, R.*
		FROM dbo.Hospital as H, dbo.Station as S, dbo.Room as R
		WHERE S.hospital_id=H.hospital_id and R.station_id=S.station_id and S.station_id=?
	),avgtime AS (
		SELECT tmp.room_id, AVG(Q.wait_min) as avg_q_time
		FROM rooms_in_station as tmp, dbo.Queue as Q
		WHERE Q.room_id=tmp.room_id
		and Q.create_date between  @fromdate and  @todate
					and Q.status!=91
					and Q.status!=0

		GROUP BY tmp.room_id
	),queues_in_room AS (
		SELECT tmp.room_id as rid, count(queue_id) as queues, avgtime.avg_q_time
		FROM rooms_in_station as tmp, dbo.Queue as Q, avgtime
		WHERE  Q.room_id=tmp.room_id
					and Q.create_date between  @fromdate and  @todate
					and avgtime.room_id = tmp.room_id
		GROUP BY tmp.room_id, avg_q_time
	), doctor AS (
		SELECT D.doctor_id, D.fullname, DW.room_id FROM dbo.Doctors AS D, dbo.Doctor_Worktime AS DW
		WHERE D.doctor_id=DW.doctor_id and DW.work_date = CONVERT(date, GETDATE()) and CONVERT(time, GETDATE()) between DW.work_start_time and DW.work_end_time
	), queues_rank as (
		SELECT room_id, queue_number_text, create_date, create_time, commit_time,call_time
		FROM (
			SELECT  *, ROW_NUMBER() OVER(PARTITION BY Queue.room_id ORDER BY Queue.call_time DESC) as rank
			FROM Queue
			where create_date BETWEEN  @fromdate and  @todate
			and call_time is not null
			and status != 91
			and status != 10
			and status != 11
		) as tmp
		WHERE tmp.rank=1
	)
	SELECT RIS.room_id as room_id, RIS.room_name as room_name, (
		SELECT 	count(Q.queue_id) as queues
		FROM Queue as Q
		WHERE Q.room_id=RIS.room_id
		and Q.create_date between  @fromdate and  @todate
			and Q.status!=10
			and Q.status!=11
	) as queues, QIR.avg_q_time as avgQueueingTime, doctor.fullname as doctor_name, queues_rank.*
	FROM rooms_in_station as RIS LEFT JOIN queues_in_room as QIR ON (RIS.room_id=QIR.rid)
		LEFT JOIN doctor ON (doctor.room_id=RIS.room_id)
		LEFT JOIN queues_rank ON (queues_rank.room_id=RIS.room_id)
	ORDER BY  RIS.order_no
	`, fDate, tDate, sID).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var tmp Room
		var lastQ LastestQ
		var highQ HighestRow
		var lowQ LowestRow
		rows.Scan(
			&tmp.RoomID, &tmp.RoomName, &tmp.Queue, &tmp.AvgQTime, &tmp.Doctor, &tmp.RoomID,
			&lastQ.NumberText, &lastQ.CreateDate, &lastQ.CreatedAt, &lastQ.CallAt, &lastQ.FinishAt,
		)
		tmp.LastestQ = &lastQ
		row := db.Raw(`
			SELECT TOP(1) queue_number_text as numberText, DATEDIFF(minute, CONVERT(Varchar(10), create_date, 112) + ' ' + CONVERT(Varchar(8), create_time), CONVERT(Varchar(10), create_date, 112) + ' ' + CONVERT(Varchar(8), call_time)) as timeWaiting FROM Queue
      WHERE room_id = ? and create_date between ? and ? and call_time  IS NOT NULL
      ORDER BY timeWaiting DESC
		`, tmp.RoomID, fDate, tDate).Row()
		row.Scan(highQ.NumberText, highQ.WaitingTime)
		row = db.Raw(`
			SELECT TOP(1) queue_number_text as numberText, DATEDIFF(minute, CONVERT(Varchar(10), create_date, 112) + ' ' + CONVERT(Varchar(8), create_time), CONVERT(Varchar(10), create_date, 112) + ' ' + CONVERT(Varchar(8), call_time)) as timeWaiting FROM Queue
	    WHERE room_id = ? and create_date between ? and ? and call_time  IS NOT NULL
	    ORDER BY timeWaiting ASC
		`, tmp.RoomID, fDate, tDate).Row()
		row.Scan(lowQ.NumberText, lowQ.WaitingTime)
		tmp.HighestRow = &highQ
		tmp.LowestRow = &lowQ
		*rooms = append(*rooms, tmp)
	}
	return nil
}

// GetMultiStationInHos : get station in hospital
func GetMultiStationInHos(hospitalID string, station *[]MultiStation) {
	db.Raw(`
		with staff as (
        	SELECT station_id, staff_id
        	From StaffConfig
        	WHERE status = 1
        	GROUP BY station_id, staff_id
        ), staffCount as (
            SELECT station_id, COUNT(staff_id) as sCount
            From staff
            GROUP BY station_id
        )
        SELECT S.station_id, S.station_name, COUNT(R.room_id) as room_amount, sc.sCount as staff_amount, ROW_NUMBER() OVER(ORDER BY S.order_no asc) AS row_index
        FROM Station as S
            LEFT JOIN Room as R ON S.station_id=R.station_id
            LEFT JOIN staffCount as sc ON S.station_id=sc.station_id
        WHERE S.hospital_id=?
        GROUP BY S.station_id, S.station_name, S.order_no, sc.sCount
	`, hospitalID).Scan(&station)
}

// GetHighestWaitingTime : get top 1 highest waiting time queue
func GetHighestWaitingTime(date string, sID uint64, model *HighWaittime) {
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

// CheckStationNameAvailable : check if stationame is available in hospital (return isAvailable)
func CheckStationNameAvailable(sName string, hID uint64) bool {
	var station Station
	db.Table("Station").Where("station_name=? AND hospital_id = ?", sName, hID).Find(&station)
	if station.StationCode != nil {
		return false
	}
	return true
}

// CheckStationCodeAvailable : check if stationcode is available (return isAvailable)
func CheckStationCodeAvailable(sCode string) bool {
	var station Station
	db.Table("Station").Where("station_code LIKE ?", sCode).Find(&station)
	if station.StationName != nil {
		return false
	}
	return true
}

// CheckStationExist : find station by id and check if station exist
func CheckStationExist(sID string) bool {
	var station Station
	db.Table("Station").Where("station_id LIKE ?", sID).Find(&station)
	if station.StationName == nil {
		return false
	}
	return true
}

// CreateStation : create station
func CreateStation(station *Station) error {
	row := db.Select("MAX(order_no)+1 as next_no").Table("Station").Row()
	row.Scan(&station.OrderNo)

	tx := db.Begin()
	if err := tx.Table("Station").Create(&station).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// EditStation : edit station info
func EditStation(sID string, station *Station) error {
	tx := db.Begin()

	if err := tx.Table("Station").Where("station_id = ?", sID).Update(&station).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// RemoveStation : remove station that match receive ID
func RemoveStation(sID string) error {
	tx := db.Begin()

	if err := tx.Table("Station").Where("station_id = ?", sID).Delete(&Station{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// // ReorderStationInHos ...
// func ReorderStationInHos(staff ReorderStationModel) error {
// 	tx := db.Begin()
// 	jobs := make(chan ReorderStationList, len(staff.Data))
// 	err := make(chan error, len(staff.Data))

// 	for w := 0; w <= 3; w++ {
// 		go reorderStationWorker(tx, jobs, err)
// 	}

// 	for _, element := range staff.Data {
// 		jobs <- element
// 	}

// 	for idx := 0; idx < len(staff.Data); idx++ {
// 		e := <-err
// 		if e != nil {
// 			tx.Rollback()
// 			return e
// 		}
// 	}
// 	db.Commit()
// 	return nil
// }

// func reorderStationWorker(tx *gorm.DB, jobs <-chan ReorderStationList, err chan<- error) {
// 	for j := range jobs {
// 		if queryErr := tx.Exec("UPDATE Station SET order_no = ? WHERE station_id = ?", j.OrderNO, j.ID).Error; err != nil {
// 			err <- queryErr
// 		}
// 		err <- nil
// 	}
// }

// ReorderStationInHos ...
func ReorderStationInHos(staff ReorderStationModel) error {
	tx := db.Begin()
	c := len(staff.Data)
	jobs := make(chan func(), c)
	err := make(chan error, c)
	pool := worker.Pool{
		Amount: 3,
		Jobs:   jobs,
		Done:   err,
	}

	for _, element := range staff.Data {
		jobs <- func() {
			if fail := tx.Exec("UPDATE Station SET order_no = ? WHERE station_id = ?", element.OrderNO, element.ID).Error; fail != nil {
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
