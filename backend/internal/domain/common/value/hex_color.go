package value

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"financetracker/internal/domain/common"
)

// HexColor HEX形式の色を表現する値オブジェクト
type HexColor struct {
	value string
}

// RGB RGB値を表現する構造体
type RGB struct {
	R int
	G int
	B int
}

// HEX形式の色の正規表現（#で始まる6桁の16進数）
var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// NewHexColor 新しいHexColorインスタンスを作成
func NewHexColor(value string) (*HexColor, error) {
	color := strings.TrimSpace(value)

	if err := validateHexColor(color); err != nil {
		return nil, err
	}

	return &HexColor{
		value: strings.ToUpper(color), // 大文字で正規化
	}, nil
}

// NewHexColorFromRGB RGB値からHexColorを作成
func NewHexColorFromRGB(r, g, b int) (*HexColor, error) {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return nil, common.NewValidationError("rgb", fmt.Sprintf("R:%d G:%d B:%d", r, g, b),
			"RGB values must be between 0 and 255")
	}

	hexValue := fmt.Sprintf("#%02X%02X%02X", r, g, b)
	return &HexColor{
		value: hexValue,
	}, nil
}

// Value HEX色の値を取得
func (h HexColor) Value() string {
	return h.value
}

// ToRGB RGB値に変換
func (h HexColor) ToRGB() (RGB, error) {
	// #を除去
	hex := h.value[1:]

	// RGBそれぞれの16進数を10進数に変換
	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return RGB{}, common.NewValidationError("hex_color", h.value, "failed to parse red component")
	}

	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return RGB{}, common.NewValidationError("hex_color", h.value, "failed to parse green component")
	}

	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return RGB{}, common.NewValidationError("hex_color", h.value, "failed to parse blue component")
	}

	return RGB{
		R: int(r),
		G: int(g),
		B: int(b),
	}, nil
}

// Equals 色の同一性判定（大文字小文字は区別しない）
func (h HexColor) Equals(other HexColor) bool {
	return strings.EqualFold(h.value, other.value)
}

// String 文字列表現
func (h HexColor) String() string {
	return h.value
}

// IsLight 明るい色かどうかを判定（テキストの可読性判定に使用）
func (h HexColor) IsLight() bool {
	rgb, err := h.ToRGB()
	if err != nil {
		return false
	}

	// 輝度計算（簡易版）
	// 輝度 = 0.299*R + 0.587*G + 0.114*B
	luminance := 0.299*float64(rgb.R) + 0.587*float64(rgb.G) + 0.114*float64(rgb.B)
	return luminance > 127.5 // 中間値より明るいかどうか
}

// GetContrastColor コントラストの良い文字色を取得（白または黒）
func (h HexColor) GetContrastColor() string {
	if h.IsLight() {
		return "#000000" // 明るい背景には黒文字
	}
	return "#FFFFFF" // 暗い背景には白文字
}

// validateHexColor HEX色のバリデーション
func validateHexColor(color string) error {
	// 空文字チェック
	if color == "" {
		return common.NewValidationError("hex_color", color, "hex color is required")
	}

	// フォーマットチェック
	if !hexColorRegex.MatchString(color) {
		return common.NewValidationError("hex_color", color,
			"hex color must be in format #RRGGBB (e.g., #3B82F6)")
	}

	return nil
}

// 定義済みの色定数
var (
	Black   = &HexColor{value: "#000000"}
	White   = &HexColor{value: "#FFFFFF"}
	Red     = &HexColor{value: "#FF0000"}
	Green   = &HexColor{value: "#00FF00"}
	Blue    = &HexColor{value: "#0000FF"}
	Yellow  = &HexColor{value: "#FFFF00"}
	Cyan    = &HexColor{value: "#00FFFF"}
	Magenta = &HexColor{value: "#FF00FF"}
)
