package logs

import (
	"fmt"
	"madaurus/dev/material/app/startup"
	"time"
)

// get the log file save path
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", startup.AppSetting.RuntimeRootPath, startup.AppSetting.LogSavePath)
}

// get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		startup.AppSetting.LogSaveName,
		time.Now().Format(startup.AppSetting.TimeFormat),
		startup.AppSetting.LogFileExt,
	)
}
