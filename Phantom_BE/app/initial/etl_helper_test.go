package initial

import (
	"reflect"
	"testing"
	"PhantomBE/global"
)

func TestParseOpeningHours(t *testing.T) {
	raw := "Mon, Wed 08:00 - 12:00 / Tue 14:00 - 18:00"
	expected := []global.OpeningHour{
		{DayOfWeek: "Monday", OpenTime: "08:00", CloseTime: "12:00"},
		{DayOfWeek: "Wednesday", OpenTime: "08:00", CloseTime: "12:00"},
		{DayOfWeek: "Tuesday", OpenTime: "14:00", CloseTime: "18:00"},
	}

	result := parseOpeningHours(raw)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
