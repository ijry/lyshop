package i18n

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// supported locales
var messages = map[string]map[string]string{
	"zh-CN": zhCN,
	"en":    en,
}

const defaultLocale = "zh-CN"

// T returns the translated message for the given key.
// It reads the locale from the gin.Context (set by LocaleMiddleware).
func T(c *gin.Context, key string) string {
	lang := Lang(c)
	if m, ok := messages[lang]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	// fallback to default locale
	if m, ok := messages[defaultLocale]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return key
}

// Lang extracts the current locale from gin.Context.
// Returns default locale if c is nil.
func Lang(c *gin.Context) string {
	if c == nil {
		return defaultLocale
	}
	if lang, ok := c.Get("lang"); ok {
		if s, ok := lang.(string); ok && s != "" {
			return s
		}
	}
	return defaultLocale
}

// LocaleMiddleware parses the Accept-Language header or lang query param
// and stores the resolved locale in gin.Context.
func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		if lang == "" {
			lang = parseAcceptLanguage(c.GetHeader("Accept-Language"))
		}
		if _, ok := messages[lang]; !ok {
			lang = defaultLocale
		}
		c.Set("lang", lang)
		c.Next()
	}
}

// parseAcceptLanguage extracts the first supported locale from the header.
func parseAcceptLanguage(header string) string {
	if header == "" {
		return defaultLocale
	}
	for _, part := range strings.Split(header, ",") {
		tag := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		// exact match
		if _, ok := messages[tag]; ok {
			return tag
		}
		// prefix match: "en-US" → "en"
		prefix := strings.SplitN(tag, "-", 2)[0]
		if _, ok := messages[prefix]; ok {
			return prefix
		}
	}
	return defaultLocale
}
