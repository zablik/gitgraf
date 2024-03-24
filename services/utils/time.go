package utils

import (
	"log"
	"time"
)

func MeasureTime(msg string, function func()) {
	start := time.Now()
	function()
	elapsed := time.Since(start)
	log.Println(msg, "> Execution time:", elapsed)
}
