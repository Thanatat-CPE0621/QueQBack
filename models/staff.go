package models

import (
	"fmt"
)

// Staff - for both model db and json
type Staff struct {
	// Staff
	StaffID    										uint32     			`gorm:"primary_key;AUTO_INCREMENT" json:"staff_id"`
	StaffName   									*string    			`json:"staff_name"`
	StaffType   									*string    			`json:"staff_type,omitempty"`
	LoginName  										string     			`json:"login_name,omitempty"`
	LoginHashPassword  						string     			`json:"login_hash_password,omitempty"`
	StationID  										int		     			`json:"station_id,omitempty"`
	RoomID  											int		     			`json:"room_id,omitempty"`
	LangCode  										*string     		`json:"lang_code,omitempty"`
	Status  											int		     			`gorm:"default:1" json:"status,omitempty"`
	StaffHospitalLogoPath  				*string     		`json:"staff_hospital_logo_path,omitempty"`
	StaffHospitalPrintLogoPath  	*string     		`json:"staff_hospital_print_logo_path,omitempty"`
	UserToken  										*string     		`json:"user_token,omitempty"`
	last_login_date  							*string     		`json:"last_login_date,omitempty"`
	last_login_time  							*string     		`json:"last_login_time,omitempty"`
	staff_parameter  							*string     		`json:"staff_parameter,omitempty"`
	order_no  										*uint32     		`json:"order_no,omitempty"`
	created_date  								*string		     	`gorm:"default:CURRENT_TIMESTAMP" json:"created_date,omitempty"`
	updated_date  								*string		     	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_date,omitempty"`
	app_version  									*string     		`json:"app_version,omitempty"`

	// Staff type
	// staff_type_id, staff_type_name, staff_type_hospital_logo_path, staff_type_hospital_receipt_amount, staff_type_hospital_receipt_desc
	StaffTypeID    									uint32     			`gorm:"-" json:"staff_type_id"`
	StaffTypeName    								*string    			`gorm:"-" json:"staff_type_name"`
	StaffTypeHospitalLogoPath   		*string    			`gorm:"-" json:"staff_type_hospital_logo_path"`
	StaffTypeHospitalReceiptAmount	*uint32    			`gorm:"-" json:"staff_type_hospital_receipt_amount"`
	StaffTypeHospitalReceiptDesc   	*string    			`gorm:"-" json:"staff_type_hospital_receipt_desc"`

	HospitalID    									uint32     			`gorm:"-" json:"hospital_id"`
	HospitalName    								*string     		`gorm:"-" json:"hospital_name"`
	StationName    									*string     		`gorm:"-" json:"station_name"`
}

// GetHospitalIDbyStaffToken - get a user
func GetHospitalIDbyStaffToken (staff *Staff, staffToken string) error {
	if err := db.Where("user_token = ?", staffToken).Table("Staff").First(&staff).Error; err != nil {
		return err
	}
	return nil
}

func GetStaffList (staffs *[]Staff, size int, page int, hID string, rID string) error {
	where := "ST.station_id=S.station_id and S.hospital_id=H.hospital_id and ST.staff_type=STT.staff_type_id "
	if hID == "" {
		where += fmt.Sprintf("and H.hospital_id like %v ", hID)
	}
	if rID == "" {
		where += fmt.Sprintf("and STT.staff_type_id like %v", rID)
	}
	if err := db.
		Where(where).
		Table("Staff as ST, Hospital as H, Station as S, StaffType as STT").
		Select(" ST.*, STT.*, H.hospital_id, H.hospital_name, S.station_name").
		Scan(&staffs).
		Offset(page).
		Limit(size).
	Error; err != nil {
		return err
	}
	return nil
}
