package util

import (
	"fmt"
	"time"
)

type Measurement struct {
	Start time.Time
}

func (m *Measurement) StartMeasure() {
	m.Start = time.Now()
}
func (m *Measurement) EndMeasure(key string) {
	duration := time.Since(m.Start)
	fmt.Println(key, "Time taken: ", duration)
}
