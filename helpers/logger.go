package helpers

import (
	"github.com/woaitsAryan/stuneckt-task/initializers"
)

func LogDatabaseError(customString string, err error, path string) {
	if appError, ok := err.(*AppError); ok {
		initializers.Logger.Warnw(customString, "Message", appError.Message, "Path", path, "Error", appError.Err)
	} else {
		initializers.Logger.Warnw(customString, "Message", err.Error(), "Path", path, "Error", err)
	}
}

func LogServerError(customString string, err error, path string) {
	if appError, ok := err.(*AppError); ok {
		initializers.Logger.Errorw(customString, "Message", appError.Message, "Path", path, "Error", appError.Err)
	} else {
		initializers.Logger.Errorw(customString, "Message", err.Error(), "Path", path, "Error", err)
	}
}
