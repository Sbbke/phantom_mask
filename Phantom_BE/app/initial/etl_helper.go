package initial

import(
	"regexp"
	"strings"
	"github.com/charmbracelet/log"
	// "os"
	"PhantomBE/global"
	// "encoding/json"
	"time"
	// "gorm.io/datatypes"

	"golang.org/x/text/cases"
    "golang.org/x/text/language"
)


func parseOpeningHours(raw string) []global.OpeningHour{
	var result []global.OpeningHour
	log.Info("Analyzed days", "raw", raw)
	pattern := regexp.MustCompile(`(?i)((?:mon|tue|wed|thu|thur|fri|sat|sun)[\w,\s-]*?)\s*:?(\d{1,2}[:.]?\d{0,2}\s*(?:am|pm)?)\s*[-–to]+\s*(\d{1,2}[:.]?\d{0,2}\s*(?:am|pm)?)`)
	segments := strings.Split(raw, "/")

	for _, segment := range segments {
	segment = strings.TrimSpace(segment)
	log.Info("Analyzed segment", "segment", segment)
	matches := pattern.FindAllStringSubmatch(segment, -1)
		for _, match := range matches {
			dayExpr := match[1]
			startTime := parseFlexibleTime(match[2])
			endTime := parseFlexibleTime(match[3])
			days := expandDays(dayExpr)
			for _, day := range days {
				result = append(result, global.OpeningHour{
					DayOfWeek: day,
					OpenTime:  startTime,
					CloseTime: endTime,
				})
			}
		}
	}
	return result
}

func parseFlexibleTime(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, ".", ":") // 支援 08.00
	if !strings.Contains(s, ":") {
		s += ":00"
	}
	t, err := time.Parse("3:04pm", s)
	if err != nil {
		t, err = time.Parse("15:04", s)
		if err != nil {
			return "00:00"
		}
	}
	return t.Format("15:04")
}
func expandDays(dayExpr string) []string {
	var result []string
	var titleCaser = cases.Title(language.English)
	parts := strings.Split(dayExpr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			if len(bounds) == 2 {
				start := strings.TrimSpace(bounds[0])
				end := strings.TrimSpace(bounds[1])
				startFull := global.ShortToFullDay[titleCaser.String(strings.ToLower(start))]
				endFull := global.ShortToFullDay[titleCaser.String(strings.ToLower(end))]
				startIdx := indexOf(global.Days, startFull)
				endIdx := indexOf(global.Days, endFull)
				if startIdx >= 0 && endIdx >= 0 {
					for i := startIdx; ; i = (i + 1) % len(global.Days) {
						result = append(result, global.Days[i])
						if i == endIdx {
							break
						}
					}
				}
			}
		} else {
			full := global.ShortToFullDay[titleCaser.String(strings.ToLower(part))]
			if full != "" {
				result = append(result, full)
			}
		}
	}
	return result
}

func indexOf(slice []string, val string) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

