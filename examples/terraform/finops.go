package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Report struct {
	TotalMonthlyCost string `json:"totalMonthlyCost"`
}

type Policy struct {
	Threshold float64 `yaml:"threshold"`
}

func main() {
	file, err := os.Open("report.json")
	if err != nil {
		log.Fatalf("❌ Error opening file: %v", err)
	}
	defer file.Close()

	var r Report
	if err := json.NewDecoder(file).Decode(&r); err != nil {
		log.Fatalf("❌ Error decoding JSON: %v", err)
	}

	tmc, err := strconv.ParseFloat(r.TotalMonthlyCost, 64)
	if err != nil {
		log.Fatalf("❌ Error parsing cost '%s': %v", r.TotalMonthlyCost, err)
	}

	policyData, err := os.ReadFile("finops-policy.yaml")
	if err != nil {
		log.Fatalf("❌ Error reading policy file: %v", err)
	}
	var p Policy
	if err := yaml.Unmarshal(policyData, &p); err != nil {
		log.Fatalf("❌ Error decoding YAML: %v", err)
	}

	// 3. Compare using dynamic threshold
	if tmc > p.Threshold {
		log.Fatalf("❌ Limit exceeded: Cost ($%.2f) is higher than threshold ($%.2f)", tmc, p.Threshold)
	}

	fmt.Printf("✅ Success: Total Monthly Cost Found: $%.2f (Threshold: $%.2f)\n", tmc, p.Threshold)
}
