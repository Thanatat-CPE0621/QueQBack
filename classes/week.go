package classes

// Week : week class model
type Week struct {
  Sunday    int32 `json:"Sunday"`
  Monday    int32 `json:"Monday"`
  Tuesday   int32 `json:"Tuesday"`
  Wednesday int32 `json:"Wednesday"`
  Thursday  int32 `json:"Thursday"`
  Friday    int32 `json:"Friday"`
  Saturday  int32 `json:"Saturday"`
}

// SetDayInWeek : set date in week day struct
func (w *Week) SetDayInWeek (day int, amount int32) {
  switch day {
      case 1: w.Monday += amount
      case 2: w.Tuesday += amount
      case 3: w.Wednesday += amount
      case 4: w.Thursday += amount
      case 5: w.Friday += amount
      case 6: w.Saturday += amount
      case 7: w.Sunday += amount
  }
}
