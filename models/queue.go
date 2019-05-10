package models

// QueueInterface : queue interface
type QueueInterface map[string]interface{}

// GetperiodQWeek : get hourly period queues in week
func (q *QueueInterface) GetperiodQWeek (hID, year, week string) error {
  rows, err := db.Raw(`
		SELECT SUBSTRING(create_time, 1, 2) as Hour, COUNT(queue_id) as Queues
        FROM Queue as Q, Hospital as H, Station as S
        WHERE Q.station_id=S.station_id AND S.hospital_id=H.hospital_id AND H.hospital_id=? AND SUBSTRING(create_date, 1, 4)=?
			AND DATEPART(wk, create_date)=?
        GROUP BY SUBSTRING(create_time, 1, 2)
	`, hID, year, week).Rows()
	if err != nil {
		return err
	}
	tmp := make(map[string]interface{})
	for rows.Next() {
		var hour string; var queues int32
		rows.Scan(&hour, &queues)
		tmp[hour] = queues
	}
	*q = tmp

	return nil
}

// GetperiodQMonth : get hourly period queues in month
func (q *QueueInterface) GetperiodQMonth (hID, date string) error {
  rows, err := db.Raw(`
		SELECT SUBSTRING(create_time, 1, 2) as Hour, COUNT(queue_id) as Queues
    FROM Queue as Q, dbo.Hospital as H, dbo.Station as S
    WHERE Q.station_id=S.station_id and S.hospital_id=H.hospital_id and H.hospital_id=? and SUBSTRING(create_date, 1, 7)=?
    GROUP BY SUBSTRING(create_time, 1, 2)
	`, hID, date).Rows()
	if err != nil {
		return err
	}
	tmp := make(map[string]interface{})
	for rows.Next() {
		var hour string; var queues int32
		rows.Scan(&hour, &queues)
		tmp[hour] = queues
	}
	*q = tmp

	return nil
}
