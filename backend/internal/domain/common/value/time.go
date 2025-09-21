package value

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"financetracker/internal/domain/common"
)

// Time HH:mm形式の時刻を表現する値オブジェクト
type Time struct {
	hour   int
	minute int
}

// HH:mm形式の時刻の正規表現
var timeRegex = regexp.MustCompile(`^([0-1][0-9]|2[0-3]):([0-5][0-9])$`)

// NewTime 新しいTimeインスタンスを作成
func NewTime(hour, minute int) (*Time, error) {
	if err := validateTime(hour, minute); err != nil {
		return nil, err
	}

	return &Time{
		hour:   hour,
		minute: minute,
	}, nil
}

// NewTimeFromString HH:mm形式の文字列からTimeインスタンスを作成
func NewTimeFromString(timeStr string) (*Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	
	if !timeRegex.MatchString(timeStr) {
		return nil, common.NewValidationError("time", timeStr, 
			"time must be in HH:mm format (e.g., 09:30, 14:15)")
	}

	parts := strings.Split(timeStr, ":")
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, common.NewValidationError("time", timeStr, "invalid hour format")
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, common.NewValidationError("time", timeStr, "invalid minute format")
	}

	return NewTime(hour, minute)
}

// NewTimeFromGoTime Go標準のtime.Timeから時刻部分のみを抽出してTimeインスタンスを作成
func NewTimeFromGoTime(t time.Time) *Time {
	return &Time{
		hour:   t.Hour(),
		minute: t.Minute(),
	}
}

// Hour 時を取得
func (t Time) Hour() int {
	return t.hour
}

// Minute 分を取得
func (t Time) Minute() int {
	return t.minute
}

// String HH:mm形式の文字列表現
func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d", t.hour, t.minute)
}

// ToGoTime 指定された日付と組み合わせてGo標準のtime.Timeを作成
func (t Time) ToGoTime(date time.Time) time.Time {
	return time.Date(
		date.Year(), date.Month(), date.Day(),
		t.hour, t.minute, 0, 0,
		date.Location(),
	)
}

// AddMinutes 指定した分数を加算（日をまたぐ場合は0時から継続）
func (t Time) AddMinutes(minutes int) *Time {
	totalMinutes := t.hour*60 + t.minute + minutes
	
	// 負の値の場合は前日からの計算
	for totalMinutes < 0 {
		totalMinutes += 24 * 60 // 1日分を加算
	}
	
	// 24時間を超える場合は翌日への繰り越し
	totalMinutes = totalMinutes % (24 * 60)
	
	newHour := totalMinutes / 60
	newMinute := totalMinutes % 60
	
	return &Time{
		hour:   newHour,
		minute: newMinute,
	}
}

// SubtractMinutes 指定した分数を減算
func (t Time) SubtractMinutes(minutes int) *Time {
	return t.AddMinutes(-minutes)
}

// MinutesSinceMidnight 0時からの経過分数を取得
func (t Time) MinutesSinceMidnight() int {
	return t.hour*60 + t.minute
}

// IsBefore 他の時刻より前かどうかを判定
func (t Time) IsBefore(other Time) bool {
	return t.MinutesSinceMidnight() < other.MinutesSinceMidnight()
}

// IsAfter 他の時刻より後かどうかを判定
func (t Time) IsAfter(other Time) bool {
	return t.MinutesSinceMidnight() > other.MinutesSinceMidnight()
}

// Equals 時刻の同一性判定
func (t Time) Equals(other Time) bool {
	return t.hour == other.hour && t.minute == other.minute
}

// IsBusinessHours 営業時間内かどうかを判定（9:00-18:00）
func (t Time) IsBusinessHours() bool {
	businessStart := Time{hour: 9, minute: 0}
	businessEnd := Time{hour: 18, minute: 0}
	
	return !t.IsBefore(businessStart) && t.IsBefore(businessEnd)
}

// Format 指定されたフォーマットで時刻を文字列化
func (t Time) Format(format string) string {
	switch format {
	case "12h":
		hour12 := t.hour
		ampm := "AM"
		if t.hour >= 12 {
			ampm = "PM"
			if t.hour > 12 {
				hour12 = t.hour - 12
			}
		}
		if hour12 == 0 {
			hour12 = 12
		}
		return fmt.Sprintf("%d:%02d %s", hour12, t.minute, ampm)
	default:
		return t.String() // デフォルトは24時間形式
	}
}

// validateTime 時刻のバリデーション
func validateTime(hour, minute int) error {
	if hour < 0 || hour > 23 {
		return common.NewValidationError("hour", hour, "hour must be between 0 and 23")
	}
	
	if minute < 0 || minute > 59 {
		return common.NewValidationError("minute", minute, "minute must be between 0 and 59")
	}
	
	return nil
}