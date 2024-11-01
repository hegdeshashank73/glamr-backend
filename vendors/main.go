package vendors

import (
	"time"

	"github.com/hegdeshashank73/glamr-backend/utils"
)

func Setup() {
	st := time.Now()
	defer utils.LogTimeTaken("init.initSetup", st)

	initAWS()
	initDB()
}
