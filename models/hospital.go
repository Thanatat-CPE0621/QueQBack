package models

import ()

type Hospital struct {
  HospitalID                  uint32      `json:"hospital_id"`
  HospitalName                *string     `json:"hospital_name"`
  HospitalName2               *string     `json:"hospital_name2,omitempty"`
  HospitalName3               *string     `json:"hospital_name3,omitempty"`
  HospitalName4               *string     `json:"hospital_name4,omitempty"`
  HospitalName5               *string     `json:"hospital_name5,omitempty"`
  HospitalName6               *string     `json:"hospital_name6,omitempty"`
  HospitalName7               *string     `json:"hospital_name7,omitempty"`
  HospitalName8               *string     `json:"hospital_name8,omitempty"`
  HospitalUID                 *string     `json:"hospital_uid,omitempty"`
  HospitalLogoPath            *string     `json:"hospital_logo_path,omitempty"`
  HospitalPrintLogoPath       *string     `json:"hospital_print_logo_path,omitempty"`
  HospitalMobileLogoPath      *string     `json:"hospital_mobile_logo_path,omitempty"`
  Status                      uint32      `json:"status,omitempty"`
  OrderNo                     *uint32     `json:"order_no,omitempty"`
  CreatedDate                 string      `json:"created_date,omitempty"`
  UpdatedDate                 string      `json:"updated_date,omitempty"`
  ApProvince                  *string     `json:"ap_province,omitempty"`
  ApLatitude                  *float64    `json:"ap_latitude,omitempty"`
  ApLongitude                 *float64    `json:"ap_longitude,omitempty"`
  ApDistanceFlag              *uint32     `json:"ap_distance_flag,omitempty"`
  ApDistanceLimit             *uint32     `json:"ap_distance_limit,omitempty"`
  ApOpenCloseRules            *string     `json:"ap_open_close_rules,omitempty"`
  ApBoxTypeCode               *string     `json:"ap_box_type_code,omitempty"`
}
