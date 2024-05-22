package logs

import (
	"fmt"
	"madaurus/dev/material/app/interfaces"
	"time"
)

// get the log file save path
func getLogFilePath(AppSetting *interfaces.App) string {
	return fmt.Sprintf("%s%s", AppSetting.RuntimeRootPath, AppSetting.LogSavePath)
}

// get the save name of the log file
func getLogFileName(AppSetting *interfaces.App) string {
	return fmt.Sprintf("%s%s.%s",
		AppSetting.LogSaveName,
		time.Now().Format(AppSetting.TimeFormat),
		AppSetting.LogFileExt,
	)
}
