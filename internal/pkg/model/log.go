package model

import "time"

type Log struct {
	JobName      string
	Command      string
	Output       string
	Err          string
	PlanTime     time.Time
	ScheduleTime time.Time
	StartTime    time.Time
	EndTime      time.Time
}
