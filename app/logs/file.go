package logs

import (
	"fmt"
	"madaurus/dev/material/app/interfaces"
	"time"
)

// get the log file save path
func getLogFilePath(AppSettings *interfaces.App) string {
	return fmt.Sprintf("%s%s", AppSettings.RuntimeRootPath, AppSettings.LogSavePath)
}

// get the save name of the log file
func getLogFileName(AppSettings *interfaces.App) string {
	return fmt.Sprintf("%s%s.%s",
		AppSettings.LogSaveName,
		time.Now().Format(AppSettings.TimeFormat),
		AppSettings.LogFileExt,
	)
}
