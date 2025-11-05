package logger

import "go.uber.org/zap"

func New(isDevelopment bool) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error

	if isDevelopment {
		// Logger untuk development: human-readable, level debug
		cfg := zap.NewDevelopmentConfig()
		logger, err = cfg.Build()
	} else {
		// Logger untuk production: format JSON, level info
		cfg := zap.NewProductionConfig()
		logger, err = cfg.Build()
	}

	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
