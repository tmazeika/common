package logging

import (
	"log"
	"os"
)

const LogFlags = log.Ldate | log.Ltime | log.LUTC | log.Lshortfile

var Logger = log.New(os.Stdout, "", LogFlags)
