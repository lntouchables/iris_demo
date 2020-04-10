package libs

import "time"

func TimeFormat(time *time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
