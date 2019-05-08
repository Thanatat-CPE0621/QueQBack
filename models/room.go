package models

// Room models
type Room struct {
	RoomID          uint32  `gorm:"primary_key;AUTO_INCREMENT" json:"room_id"`
	RoomCode        *string `json:"room_code,omitempty"`
	RoomName        *string `json:"room_name"`
	RoomName2       *string `json:"room_name2,omitempty"`
	RoomName3       *string `json:"room_name3,omitempty"`
	RoomName4       *string `json:"room_name4,omitempty"`
	RoomName5       *string `json:"room_name5,omitempty"`
	RoomName6       *string `json:"room_name6,omitempty"`
	RoomName7       *string `json:"room_name7,omitempty"`
	RoomName8       *string `json:"room_name8,omitempty"`
	StationID       uint32  `json:"station_id,omitempty"`
	Status          *string `json:"status,omitempty"`
	ReasonText      string  `json:"reason_text,omitempty"`
	OrderNo         *uint32 `json:"order_no,omitempty"`
	CreatedDate     *string `json:"created_date,omitempty"`
	UpdatedDate     *uint32 `json:"updated_date,omitempty"`
	AmountMobileMsg *int    `json:"amount_mobile_msg,omitempty"`
	AmountNotif     *int    `json:"amount_notif,omitempty"`

	Queue    *int32   `gorm:"-" json:"queues,omitempty"`
	AvgQTime *int32		`gorm:"-" json:"avgQueueingTime,omitempty"`
	Doctor   *string  `gorm:"-" json:"doctor_name,omitempty"`

	LastestQ 		*LastestQ `gorm:"-" json:"lastestQueue,omitempty"`
	HighestRow 	*HighestRow `gorm:"-" json:"highestWaitingTimeQueue,omitempty"`
	LowestRow 	*LowestRow `gorm:"-" json:"lowestWaitingTimeQueue,omitempty"`
}

// LastestQ : Last Queue in Room Model
type LastestQ struct {
	NumberText *string `json:"numberText,omitempty"`
	CreateDate *string `json:"createDate,omitempty"`
	CreatedAt  *string `json:"createdAt,omitempty"`
	CallAt     *string `json:"callAt,omitempty"`
	FinishAt   *string `json:"finishedAt,omitempty"`
}

// HighestRow : Highest Waiting Time Queue in Room
type HighestRow struct {
	NumberText  *string  `json:"numberText"`
	WaitingTime *float64 `json:"timeWaiting"`
}

// LowestRow : Lowest Waiting Time Queue in Room
type LowestRow struct {
	NumberText  *string  `json:"numberText"`
	WaitingTime *float64 `json:"timeWaiting"`
}
