package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// GenerateCaptcha 生成数字验证码，返回 (captchaId, svgBody, error)
func GenerateCaptcha(rds *redis.Redis, prefix string, expire int) (string, string, error) {
	// 生成 4 位数字验证码
	code, err := randomDigits(4)
	if err != nil {
		return "", "", err
	}

	id, err := randomHex(16)
	if err != nil {
		return "", "", err
	}

	// 存入 Redis
	key := prefix + id
	if err := rds.Setex(key, code, expire); err != nil {
		return "", "", err
	}

	// 生成 SVG
	svg := generateSVG(code)
	return id, svg, nil
}

// VerifyCaptcha 校验图形验证码，校验成功后删除
func VerifyCaptcha(rds *redis.Redis, prefix, id, answer string) bool {
	if id == "" || answer == "" {
		return false
	}
	key := prefix + id
	val, err := rds.Get(key)
	if err != nil || val == "" {
		return false
	}
	// 校验成功后删除（一次性）
	_, _ = rds.Del(key)
	return strings.EqualFold(val, answer)
}

func randomDigits(n int) (string, error) {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		d, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		sb.WriteString(fmt.Sprintf("%d", d))
	}
	return sb.String(), nil
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func generateSVG(code string) string {
	// 生成简单的 SVG 验证码图片
	width := 120
	height := 40
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, height, width, height))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#f0f0f0"/>`, width, height))

	// 绘制干扰线
	sb.WriteString(`<line x1="10" y1="5" x2="110" y2="35" stroke="#ccc" stroke-width="1"/>`)
	sb.WriteString(`<line x1="5" y1="30" x2="115" y2="10" stroke="#ddd" stroke-width="1"/>`)

	// 绘制验证码字符
	colors := []string{"#e74c3c", "#3498db", "#2ecc71", "#9b59b6", "#e67e22"}
	for i, ch := range code {
		x := 20 + i*25
		y := 28
		color := colors[i%len(colors)]
		rotate := (i*7)%30 - 15
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="24" font-weight="bold" fill="%s" transform="rotate(%d %d %d)">%c</text>`,
			x, y, color, rotate, x, y, ch))
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}
