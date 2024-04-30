package renderer

import (
	"encoding/json"
	"fmt"
)

func (r Renderer) PrintReport() {
	report := r.wgpuInstance.GenerateReport()
	buf, _ := json.MarshalIndent(report, "", "  ")
	fmt.Print(string(buf))
}
