package i18n

import (
	"embed"
	"fmt"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFS embed.FS

var (
	bundle      *i18n.Bundle
	localizers  map[string]*i18n.Localizer
	defaultLang string
	mu          sync.RWMutex
	initialized bool
)

// Config holds the i18n settings.
type Config struct {
	DefaultLanguage    string   // default language (e.g. "ko")
	SupportedLanguages []string // supported languages (e.g. ["ko", "en"])
}

// Init initializes the i18n system.
func Init(cfg *Config) error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
	}

	// Apply defaults
	if cfg.DefaultLanguage == "" {
		cfg.DefaultLanguage = "ko"
	}
	if len(cfg.SupportedLanguages) == 0 {
		cfg.SupportedLanguages = []string{"ko", "en"}
	}

	defaultLang = cfg.DefaultLanguage

	// Create the bundle (sets the default language)
	bundle = i18n.NewBundle(language.Korean)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// Load the locale file for each supported language
	localizers = make(map[string]*i18n.Localizer)

	for _, lang := range cfg.SupportedLanguages {
		filename := fmt.Sprintf("locales/%s.yaml", lang)
		data, err := localeFS.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to load locale file (%s): %w", filename, err)
		}

		if _, err := bundle.ParseMessageFileBytes(data, filename); err != nil {
			return fmt.Errorf("failed to parse locale file (%s): %w", filename, err)
		}

		localizers[lang] = i18n.NewLocalizer(bundle, lang)
	}

	initialized = true
	return nil
}

// Translate translates a message key into the given language.
// lang: language code (e.g. "ko", "en")
// key: message key (e.g. "auth.signup.success")
// args: template variables (optional)
func Translate(lang, key string, args ...map[string]interface{}) string {
	mu.RLock()
	defer mu.RUnlock()

	if !initialized {
		return key
	}

	// Get the localizer (fall back to the default language if missing)
	localizer, ok := localizers[lang]
	if !ok {
		localizer = localizers[defaultLang]
	}
	if localizer == nil {
		return key
	}

	// Set template data
	var templateData map[string]interface{}
	if len(args) > 0 && args[0] != nil {
		templateData = args[0]
	}

	// Try to translate
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: templateData,
	})
	if err != nil {
		// Fall back to the default language
		if lang != defaultLang {
			if defaultLocalizer := localizers[defaultLang]; defaultLocalizer != nil {
				if fallbackMsg, fallbackErr := defaultLocalizer.Localize(&i18n.LocalizeConfig{
					MessageID:    key,
					TemplateData: templateData,
				}); fallbackErr == nil {
					return fallbackMsg
				}
			}
		}
		// If the fallback also fails, return the key itself
		return key
	}

	return msg
}

// T is shorthand for Translate.
func T(lang, key string, args ...map[string]interface{}) string {
	return Translate(lang, key, args...)
}

// GetSupportedLanguages returns the supported languages.
func GetSupportedLanguages() []string {
	mu.RLock()
	defer mu.RUnlock()

	languages := make([]string, 0, len(localizers))
	for lang := range localizers {
		languages = append(languages, lang)
	}
	return languages
}

// GetDefaultLanguage returns the default language.
func GetDefaultLanguage() string {
	mu.RLock()
	defer mu.RUnlock()
	return defaultLang
}

// IsInitialized reports whether i18n has been initialized.
func IsInitialized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return initialized
}

// Reset clears the i18n state (for tests).
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	bundle = nil
	localizers = nil
	defaultLang = ""
	initialized = false
}

// ContextKeyLocale is the key under which the locale is stored in the request context.
const ContextKeyLocale = "locale"

// LocaleGetter reads the locale from a context.
// A minimal interface is used to avoid a circular dependency on gin.Context.
type LocaleGetter interface {
	Get(key any) (value any, exists bool)
}

// GetLocaleFromContext returns the current locale from the context.
func GetLocaleFromContext(ctx LocaleGetter) string {
	if locale, exists := ctx.Get(ContextKeyLocale); exists {
		if l, ok := locale.(string); ok {
			return l
		}
	}
	return GetDefaultLanguage()
}
