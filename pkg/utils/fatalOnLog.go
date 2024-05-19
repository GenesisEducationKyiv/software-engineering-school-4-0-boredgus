package utils

type Fatal interface {
	Fatalf(template string, args ...interface{})
}

func FatalOnError(err error, logger Fatal, msg string) {
	if err != nil {
		logger.Fatalf("%s: %v", msg, err.Error())
	}
}
