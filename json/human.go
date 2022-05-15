package json

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var suffix = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// BytesToHuman returns size in bytes as human readable with a suffix included
// in the string. Supported: B, KB, MG, GB, TB, PB and EB
func BytesToHuman(i int64) string {
	expo := math.Floor(math.Log10(float64(i)) / math.Log10(1024.0))
	hum := float64(i) / math.Pow(1024.0, expo)
	return fmt.Sprintf("%.2f%s", hum, suffix[int(expo)])
}

// ToBytes takes a human readable size and convert it to bytes.
// Supported: B, KB, MB, GB, TB, PB and EB
// e.g 10.023MB
// TODO: add bigInt support
func ToBytes(s string) (float64, error) {
	s = strings.ToUpper(s)
	var pow float64
	re := regexp.MustCompile(`^\d*[.,]?\d+`)
	fact, err := strconv.ParseFloat(re.FindString(s), 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse string: %v", err)
	}
	// checking if suffix exists
	suffre := regexp.MustCompile(`\D+$`)
	actualSuffix := suffre.FindString(s)
	if !arrContains(suffix, actualSuffix) {
		return 0.0, fmt.Errorf("suffix '%s' not supported", actualSuffix)
	}
	// finding power of
	for index, val := range suffix {
		if strings.Contains(s, val) {
			pow = float64(index)
		}
	}
	return fact * math.Pow(1024.0, pow), nil
}

// arrContains checks whether str is present in array or not
func arrContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
