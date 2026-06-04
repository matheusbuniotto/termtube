package ui

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseSeek parses jump targets: "1:30", "1 30", "90", "1m30s".
func ParseSeek(input string) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("empty seek position")
	}

	// :jump 1:30 or jump 1:30
	if fields := strings.Fields(input); len(fields) > 1 {
		cmd := strings.ToLower(strings.TrimPrefix(fields[0], ":"))
		if cmd == "jump" || cmd == "seek" || cmd == "g" {
			input = strings.Join(fields[1:], " ")
		}
	}

	lower := strings.ToLower(input)
	if strings.Contains(lower, "m") && strings.Contains(lower, "s") {
		return parseMinSecTokens(lower)
	}

	if strings.Contains(input, ":") {
		parts := strings.SplitN(input, ":", 2)
		m, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		s, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return 0, fmt.Errorf("invalid time %q (use m:ss)", input)
		}
		if m < 0 || s < 0 || s >= 60 {
			return 0, fmt.Errorf("invalid time %q", input)
		}
		return float64(m*60 + s), nil
	}

	fields := strings.Fields(input)
	if len(fields) == 2 {
		m, err1 := strconv.Atoi(fields[0])
		s, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return 0, fmt.Errorf("invalid time %q (use minutes seconds)", input)
		}
		if m < 0 || s < 0 || s >= 60 {
			return 0, fmt.Errorf("invalid time %q", input)
		}
		return float64(m*60 + s), nil
	}

	sec, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid time %q", input)
	}
	if sec < 0 {
		return 0, fmt.Errorf("invalid time %q", input)
	}
	return float64(sec), nil
}

func parseMinSecTokens(s string) (float64, error) {
	var min, sec int
	_, err := fmt.Sscanf(s, "%dm%ds", &min, &sec)
	if err != nil {
		_, err = fmt.Sscanf(s, "%dm %ds", &min, &sec)
	}
	if err != nil {
		return 0, fmt.Errorf("invalid time %q", s)
	}
	return float64(min*60 + sec), nil
}
