package exception

import (
	"product-service/pkg/logger"

	"go.uber.org/zap"
)

func FailOnError(err error, msg string) {
	if err != nil {
		logger.Error(msg, zap.Error(err))
	}
}
