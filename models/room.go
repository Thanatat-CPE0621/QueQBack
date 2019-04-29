package models

import ()

// Room models
type Room struct {
  RoomID                   uint32        `json:"room_id"`
  RoomCode                 *string       `json:"room_code,omitempty"`
  RoomName                 *string       `json:"room_name"`
  RoomName2                *string       `json:"room_name2,omitempty"`
  RoomName3                *string       `json:"room_name3,omitempty"`
  RoomName4                *string       `json:"room_name4,omitempty"`
  RoomName5                *string       `json:"room_name5,omitempty"`
  RoomName6                *string       `json:"room_name6,omitempty"`
  RoomName7                *string       `json:"room_name7,omitempty"`
  RoomName8                *string       `json:"room_name8,omitempty"`
  StationID                uint32       `json:"station_id,omitempty"`
  Status                   *string       `json:"status,omitempty"`
  ReasonText               string        `json:"reason_text,omitempty"`
  OrderNo                  *uint32       `json:"order_no,omitempty"`
  CreatedDate              *string       `json:"created_date,omitempty"`
  UpdatedDate              *uint32       `json:"updated_date,omitempty"`
  AmountMobileMsg        *int          `json:"amount_mobile_msg,omitempty"`
  AmountNotif             *int          `json:"amount_notif,omitempty"`
}
