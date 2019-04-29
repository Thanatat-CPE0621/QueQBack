package models

// StaffType - for both model db and json
type StaffType struct {
	// Staff type
	StaffTypeID                    *uint32 `json:"staff_type_id"`
	StaffTypeName                  *string `json:"staff_type_name,omitempty"`
	StaffTypeHospitalLogoPath      *string `json:"staff_type_hospital_logo_path,omitempty"`
	StaffTypeHospitalReceiptAmount *uint32 `json:"staff_type_hospital_receipt_amount,omitempty"`
	StaffTypeHospitalReceiptDesc   *string `json:"staff_type_hospital_receipt_desc,omitempty"`

	IsSuperAdmin *bool `gorm:"-" json:"is_super_admin,omitempty"`
}

func GetStaffTypeList(types *[]StaffType) error {
	if err := db.Table("StaffType").Find(&types).Error; err != nil {
		return err
	}
	return nil
}
