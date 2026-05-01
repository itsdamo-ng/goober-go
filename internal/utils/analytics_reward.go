package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ValidateAnalyticsRewardName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("analytics_reward: name is required")
	}
	if len(name) < 2 {
		return fmt.Errorf("analytics_reward: name must be at least 2 characters")
	}
	if len(name) > 500 {
		return fmt.Errorf("analytics_reward: name must not exceed 500 characters")
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) && r != '-' && r != '_' && r != '.' {
			return fmt.Errorf("analytics_reward: name contains invalid character: %c", r)
		}
	}
	return nil
}

func ValidateAnalyticsRewardEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("analytics_reward: email is required")
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return fmt.Errorf("analytics_reward: invalid email format")
	}
	if parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("analytics_reward: invalid email format")
	}
	if !strings.Contains(parts[1], ".") {
		return fmt.Errorf("analytics_reward: invalid email domain")
	}
	if len(email) > 254 {
		return fmt.Errorf("analytics_reward: email exceeds maximum length")
	}
	return nil
}

func ValidateAnalyticsRewardPhone(phone string) error {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return fmt.Errorf("analytics_reward: phone is required")
	}
	digits := 0
	for _, r := range phone {
		if unicode.IsDigit(r) {
			digits++
		} else if r != '+' && r != '-' && r != ' ' && r != '(' && r != ')' {
			return fmt.Errorf("analytics_reward: phone contains invalid character: %c", r)
		}
	}
	if digits < 7 || digits > 15 {
		return fmt.Errorf("analytics_reward: phone must have 7-15 digits, got %d", digits)
	}
	return nil
}

func SlugifyAnalyticsReward(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	var result []rune
	prevDash := false
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result = append(result, r)
			prevDash = false
		} else if !prevDash && len(result) > 0 {
			result = append(result, '-')
			prevDash = true
		}
	}
	s := string(result)
	s = strings.TrimRight(s, "-")
	return s
}

func TruncateAnalyticsRewardText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return text[:maxLen]
	}
	return text[:maxLen-3] + "..."
}

func FormatAnalyticsRewardAmount(amount float64, currency string) string {
	if currency == "" {
		currency = "USD"
	}
	formatted := strconv.FormatFloat(amount, 'f', 2, 64)
	parts := strings.SplitN(formatted, ".", 2)
	whole := parts[0]
	var withCommas []byte
	for i, c := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			withCommas = append(withCommas, ',')
		}
		withCommas = append(withCommas, byte(c))
	}
	return fmt.Sprintf("%s %s.%s", currency, string(withCommas), parts[1])
}

func FormatAnalyticsRewardDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func FormatAnalyticsRewardDateShort(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func FormatAnalyticsRewardDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

func ParseAnalyticsRewardID(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("analytics_reward: ID is required")
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("analytics_reward: invalid ID format: %s", s)
	}
	if id <= 0 {
		return 0, fmt.Errorf("analytics_reward: ID must be positive")
	}
	return id, nil
}

func GenerateAnalyticsRewardCode(prefix string, id int64) string {
	if prefix == "" {
		prefix = "analytics_reward"
	}
	return fmt.Sprintf("%s_%06d", strings.ToUpper(prefix), id)
}

func NormalizeAnalyticsRewardCountry(country string) string {
	country = strings.ToUpper(strings.TrimSpace(country))
	if len(country) > 2 {
		common := map[string]string{
			"UNITED STATES": "US", "UNITED KINGDOM": "GB",
			"GERMANY": "DE", "FRANCE": "FR", "JAPAN": "JP",
			"CHINA": "CN", "INDIA": "IN", "BRAZIL": "BR",
			"CANADA": "CA", "AUSTRALIA": "AU", "MEXICO": "MX",
			"SPAIN": "ES", "ITALY": "IT", "NETHERLANDS": "NL",
		}
		if code, ok := common[country]; ok {
			return code
		}
	}
	return country
}

func SanitizeAnalyticsRewardInput(input string) string {
	input = strings.TrimSpace(input)
	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"&", "&amp;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(input)
}

func CompareAnalyticsRewardVersions(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func MergeAnalyticsRewardStringMaps(base, override map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		result[k] = v
	}
	return result
}

func AnalyticsRewardContainsAny(s string, substrs ...string) bool {
	s = strings.ToLower(s)
	for _, sub := range substrs {
		if strings.Contains(s, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
