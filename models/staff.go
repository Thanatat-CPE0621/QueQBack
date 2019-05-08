package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/paiduay/queq-hospital-api/utils"
)

// Staff - for both model db and json
type Staff struct {
	// Staff
	StaffID                    uint32  `gorm:"primary_key;AUTO_INCREMENT" json:"staff_id,omitempty"`
	StaffName                  *string `json:"staff_name" binding:"required"`
	StaffType                  *uint32 `json:"staff_type,omitempty" binding:"required"`
	LoginName                  *string `json:"login_name,omitempty" binding:"required"`
	LoginHashPassword          *string `json:"login_hash_password,omitempty"`
	StationID                  *int    `json:"station_id,omitempty" binding:"required"`
	RoomID                     *int    `json:"room_id,omitempty" binding:"required"`
	LangCode                   *string `json:"lang_code,omitempty"`
	Status                     *int    `gorm:"default:1" json:"status,omitempty"`
	StaffHospitalLogoPath      *string `json:"staff_hospital_logo_path,omitempty"`
	StaffHospitalPrintLogoPath *string `json:"staff_hospital_print_logo_path,omitempty"`
	UserToken                  *string `gorm:"default:NULL" json:"user_token,omitempty"`
	LastLoginDate              *string `gorm:"default:NULL" json:"last_login_date,omitempty"`
	LastLoginTime              *string `gorm:"default:NULL" json:"last_login_time,omitempty"`
	StaffParameter             *string `json:"staff_parameter,omitempty"`
	OrderNo                    *uint32 `json:"order_no,omitempty"`
	CreatedDate                *string `gorm:"default:CURRENT_TIMESTAMP" json:"created_date,omitempty"`
	UpdatedDate                *string `gorm:"default:CURRENT_TIMESTAMP" json:"updated_date,omitempty"`
	AppVersion                 *string `gorm:"default:NULL" json:"app_version,omitempty"`

	// Staff type
	StaffTypeName                  *string `gorm:"-" json:"staff_type_name,omitempty"`
	StaffTypeHospitalLogoPath      *string `gorm:"-" json:"staff_type_hospital_logo_path,omitempty"`
	StaffTypeHospitalReceiptAmount *uint32 `gorm:"-" json:"staff_type_hospital_receipt_amount,omitempty"`
	StaffTypeHospitalReceiptDesc   *string `gorm:"-" json:"staff_type_hospital_receipt_desc,omitempty"`

	// Staff config
	StaffConfig []StaffConfig `gorm:"-" json:"staff_config,omitempty"`

	// Decorator
	HospitalID   *uint32 `gorm:"-" json:"hospital_id,omitempty"`
	HospitalName *string `gorm:"-" json:"hospital_name,omitempty"`
	StationName  *string `gorm:"-" json:"station_name,omitempty"`
	RoomName     *string `gorm:"-" json:"room_name,omitempty"`
}

// StaffInfomation ....
type StaffInfomation struct {
	User            Staff     `json:"user"`
	Role            StaffType `json:"role,omitempty"`
	Hospital        Hospital  `json:"hospital,omitempty"`
	Station         Station   `json:"station,omitempty"`
	Room            Room      `json:"room,omitempty"`
	QueueNumberType *uint32   `json:"queueNumberType,omitempty"`
	ReceiptAmount   *string   `json:"receiptAmount,omitempty"`
	StationName     *string   `json:"stationName,omitempty"`
}

// ReorderStaffModel ...
type ReorderStaffModel struct {
	HospitalID string             `json:"hospitalID" binding:"required"`
	Data       []ReorderStaffList `json:"data" binding:"required"`
}

// ReorderStaffList :
type ReorderStaffList struct {
	ID      uint32 `json:"id"`
	OrderNO uint32 `json:"order_no"`
}

// GetHospitalIDbyStaffToken - get a users
func GetHospitalIDbyStaffToken(staff *Staff, staffToken string) error {
	if err := db.Where("user_token = ?", staffToken).Table("Staff").First(&staff).Error; err != nil {
		return err
	}
	return nil
}

// GetStaffList - get staff list in hospital
func GetStaffList(staffs *[]Staff, size int, page int, hID string, rID string) error {
	where := "ST.station_id=S.station_id and S.hospital_id=H.hospital_id and ST.staff_type=STT.staff_type_id "
	if hID != "" {
		where += fmt.Sprintf("and H.hospital_id like %v ", hID)
	}
	if rID != "" {
		where += fmt.Sprintf("and STT.staff_type_id like %v", rID)
	}
	rows, err := db.
		Table("Staff as ST, Hospital as H, Station as S, StaffType as STT").
		Select(" ST.*, STT.*, H.hospital_id, H.hospital_name, S.station_name").
		Where(where).
		Order("ST.order_no").
		Offset((page - 1) * size).
		Limit(size).
		Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var staff Staff
		db.ScanRows(rows, &staff)
		*staffs = append(*staffs, staff)
	}
	return nil
}

// GetAdminStaffList - get admin list in hospital
func GetAdminStaffList(staffs *[]Staff, size int, page int, hID string) error {
	where := "ST.station_id=S.station_id and S.hospital_id=H.hospital_id and ST.staff_type=STT.staff_type_id "
	if hID != "" {
		where += fmt.Sprintf("and H.hospital_id like %v ", hID)
	}

	rows, err := db.
		Table("Staff as ST, Hospital as H, Station as S, StaffType as STT").
		Select(" ST.*, STT.*, H.hospital_id, H.hospital_name, S.station_name").
		Where(where).
		Order("ST.order_no").
		Offset((page - 1) * size).
		Limit(size).
		Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var staff Staff
		db.ScanRows(rows, &staff)
		*staffs = append(*staffs, staff)
	}
	return nil
}

// GetAdminStaffListInStation - get admin list in station
func GetAdminStaffListInStation(staffs *[]Staff, size int, page int, sID string) error {
	where := "ST.station_id = S.station_id and S.hospital_id = H.hospital_id and ST.staff_type = STT.staff_type_id "
	if sID != "" {
		where += fmt.Sprintf("and S.station_id like %v ", sID)
	}

	rows, err := db.
		Table("Staff as ST, Hospital as H, Station as S, StaffType as STT").
		Select("ST.staff_id, ST.staff_name, STT.staff_type_name, H.hospital_name, S.station_name").
		Where(where).
		Order("ST.order_no").
		Offset((page - 1) * size).
		Limit(size).
		Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var staff Staff
		db.ScanRows(rows, &staff)
		*staffs = append(*staffs, staff)
	}
	return nil
}

// GetAdminStaffListInRoom - get admin list in station
func GetAdminStaffListInRoom(staffs *[]Staff, size int, page int, rID string) error {
	where := "ST.station_id = S.station_id and S.hospital_id = H.hospital_id and ST.staff_type = STT.staff_type_id "
	if rID != "" {
		where += fmt.Sprintf("and R.room_id like %v ", rID)
	}

	rows, err := db.
		Table("Staff as ST, Hospital as H, Station as S, StaffType as STT, Room as R").
		Select("ST.staff_id, ST.staff_name, STT.staff_type_name, H.hospital_name, R.room_name").
		Where(where).
		Order("ST.order_no").
		Offset((page - 1) * size).
		Limit(size).
		Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var staff Staff
		db.ScanRows(rows, &staff)
		*staffs = append(*staffs, staff)
	}
	return nil
}

// GetStaffInfomation : get staff infomation (staff, stafftype, hospital, station)
func GetStaffInfomation(staff *StaffInfomation, staffID string) error {
	row := db.Raw(`
		WITH SAH as (
	    SELECT S.*, H.hospital_name FROM dbo.Station as S, dbo.Hospital as H WHERE S.hospital_id = H.hospital_id
    )
    SELECT
			ST.staff_id, ST.staff_name, ST.login_name, ST.user_token,
			ST.created_date, ST.updated_date, ST.last_login_date, ST.last_login_time,
			STT.staff_type_id, STT.staff_type_name,
			SAH.hospital_id, SAH.hospital_name, SAH.station_id, SAH.station_name,
			R.room_id, R.room_name
    FROM Staff as ST
        LEFT JOIN SAH ON ST.station_id=SAH.station_id
        LEFT JOIN dbo.Room AS R ON R.room_id=ST.room_id
        LEFT JOIN dbo.StaffType AS STT ON STT.staff_type_id=ST.staff_type
    WHERE ST.staff_id=?
	`, staffID).Row()
	row.Scan(&staff.User.StaffID, &staff.User.StaffName, &staff.User.LoginName, &staff.User.UserToken,
		&staff.User.CreatedDate, &staff.User.UpdatedDate, &staff.User.LastLoginDate, &staff.User.LastLoginTime,
		&staff.Role.StaffTypeID, &staff.Role.StaffTypeName,
		&staff.Hospital.HospitalID, &staff.Hospital.HospitalName, &staff.Station.StationID, &staff.Station.StationName,
		&staff.Room.RoomID, &staff.Room.RoomName)
	db.Where("staff_id=?", staffID).Table("StaffConfig").Find(&staff.User.StaffConfig)
	return nil
}

// FindDuplicateStaff : find duplicate login_name in table
func FindDuplicateStaff(staff *Staff) bool {
	var tmpStaff Staff
	db.Where("login_name=?", staff.LoginName).Table("Staff").Scan(&tmpStaff)
	if tmpStaff.LoginName != nil {
		return true
	}
	return false
}

// CheckUserExist : Check if login_name exist in station (return isLoginNameAvailable)
func CheckUserExist(loginName string, stationID uint64) bool {
	var tmpStaff Staff
	db.Where("login_name = ? AND station_id = ?", loginName, stationID).Table("Staff").Scan(&tmpStaff)
	if tmpStaff.LoginName != nil {
		return true
	}
	return false
}

// CheckUserExistbyID : Check if staff exist in staff's table (return isExist)
func CheckUserExistbyID(staffID uint64) bool {
	var tmpStaff Staff
	db.Where("staff_id = ?", staffID).Table("Staff").Scan(&tmpStaff)
	if tmpStaff.LoginName != nil {
		return true
	}
	return false
}

// StaffLogin : return is error, is found user, error
func StaffLogin(staff *Staff, data *StaffInfomation) (bool, bool, error) {
	var superAdmin uint32 = 53
	var password *string

	password = utils.PasswordHashing(*staff.LoginHashPassword)
	row := db.
		Table("Staff as St, Hospital as H, Station as S, StaffType as Stt").
		Select(`
		H.hospital_id, H.hospital_name, H.hospital_logo_path, H.hospital_print_logo_path,
		St.login_name, St.staff_name, Stt.staff_type_id, Stt.staff_type_name
	`).
		Where("St.staff_type=Stt.staff_type_id and St.station_id=S.station_id and S.hospital_id=H.hospital_id and St.status = 1 and St.login_name=? AND St.login_hash_password=?", staff.LoginName, *password).
		Row()
	row.Scan(&data.Hospital.HospitalID, &data.Hospital.HospitalName, &data.Hospital.HospitalLogoPath, &data.Hospital.HospitalPrintLogoPath,
		&data.User.LoginName, &data.User.StaffName, &data.Role.StaffTypeID, &data.Role.StaffTypeName)
	if data.User.StaffName == nil {
		return false, false, nil
	}

	token := utils.TokenGenerating(*data.User.LoginName, *staff.LoginHashPassword)
	tx := db.Begin()
	if err := tx.Exec(`
			UPDATE Staff
	    SET user_token=?, last_login_date=?, last_login_time=?
	    WHERE login_name=? AND login_hash_password=?
	`, token, time.Now().Format("2006-01-02"), time.Now().Format("15:04:05"), staff.LoginName, *password).Error; err != nil {
		tx.Rollback()
		return true, false, err
	}
	tx.Commit()
	isSuper := *data.Role.StaffTypeID == superAdmin
	data.User.UserToken = &token
	data.Role.IsSuperAdmin = &isSuper
	return false, true, nil
}

// CreateStaff : create staff
func CreateStaff(staff *Staff) error {
	staff.LoginHashPassword = utils.PasswordHashing(*staff.LoginHashPassword)
	tx := db.Begin()
	if err := tx.Table("Staff").Create(&staff).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := CreateStaffConfigArray(staff.StaffID, staff.StaffConfig, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// EditStaff : edit staff
func EditStaff(staff *Staff) error {
	tx := db.Begin()
	staff.LoginHashPassword = utils.PasswordHashing(*staff.LoginHashPassword)
	if err := tx.Table("Staff").Where("staff_id = ?", staff.StaffID).Save(&staff).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := UpdateStaffConfigArray(*staff, staff.StaffConfig, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// ReorderStaffInHos ...
func ReorderStaffInHos(staff ReorderStaffModel) error {
	tx := db.Begin()
	jobs := make(chan ReorderStaffList, len(staff.Data))
	err := make(chan error, len(staff.Data))

	for w := 0; w <= 3; w++ {
		go reorderWorker(tx, jobs, err)
	}

	for _, element := range staff.Data {
		jobs <- element
	}

	for idx := 0; idx < len(staff.Data); idx++ {
		e := <-err
		if e != nil {
			tx.Rollback()
			return e
		}
	}
	db.Commit()
	return nil
}

func reorderWorker(tx *gorm.DB, jobs <-chan ReorderStaffList, err chan<- error) {
	for j := range jobs {
		if queryErr := tx.Exec("UPDATE Staff SET order_no = ? WHERE staff_id = ?", j.OrderNO, j.ID).Error; err != nil {
			err <- queryErr
		}
		err <- nil
	}
}
