package i18n

import (
	"testing"
)

func TestInit(t *testing.T) {
	// Reset state before the test
	Reset()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "init with explicit config",
			config: &Config{
				DefaultLanguage:    "ko",
				SupportedLanguages: []string{"ko", "en"},
			},
			wantErr: false,
		},
		{
			name: "init with empty config (defaults)",
			config: &Config{
				DefaultLanguage:    "",
				SupportedLanguages: []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset() // reset state before each test
			err := Init(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !IsInitialized() {
				t.Error("IsInitialized() returned false after Init()")
			}
		})
	}
}

func TestTranslate(t *testing.T) {
	// Initialize before the test
	Reset()
	if err := Init(&Config{
		DefaultLanguage:    "ko",
		SupportedLanguages: []string{"ko", "en"},
	}); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	tests := []struct {
		name     string
		lang     string
		key      string
		args     map[string]interface{}
		expected string // if empty, the exact-match check is skipped
		checkNot string // result must not equal this value
	}{
		{
			name:     "Korean signup success message",
			lang:     "ko",
			key:      "auth.signup.success",
			expected: "회원가입이 완료되었습니다",
		},
		{
			name:     "English signup success message",
			lang:     "en",
			key:      "auth.signup.success",
			expected: "Signup completed successfully",
		},
		{
			name:     "Korean error message",
			lang:     "ko",
			key:      "auth.signup.email_required",
			expected: "이메일은 필수 입력 항목입니다",
		},
		{
			name:     "English error message",
			lang:     "en",
			key:      "auth.signup.email_required",
			expected: "Email is required",
		},
		{
			name:     "unsupported language falls back to the default language",
			lang:     "ja",
			key:      "auth.signup.success",
			expected: "회원가입이 완료되었습니다", // falls back to the default language (ko)
		},
		{
			name:     "unknown key returns the key itself",
			lang:     "ko",
			key:      "nonexistent.key",
			expected: "nonexistent.key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.args != nil {
				result = Translate(tt.lang, tt.key, tt.args)
			} else {
				result = Translate(tt.lang, tt.key)
			}

			if tt.expected != "" && result != tt.expected {
				t.Errorf("Translate(%q, %q) = %q, want %q", tt.lang, tt.key, result, tt.expected)
			}

			if tt.checkNot != "" && result == tt.checkNot {
				t.Errorf("Translate(%q, %q) = %q, should not be %q", tt.lang, tt.key, result, tt.checkNot)
			}
		})
	}
}

func TestT(t *testing.T) {
	// T is shorthand for Translate, so both must return the same result
	Reset()
	if err := Init(&Config{
		DefaultLanguage:    "ko",
		SupportedLanguages: []string{"ko", "en"},
	}); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	key := "auth.signup.success"
	translateResult := Translate("ko", key)
	tResult := T("ko", key)

	if translateResult != tResult {
		t.Errorf("T() = %q, Translate() = %q, results differ", tResult, translateResult)
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	Reset()
	if err := Init(&Config{
		DefaultLanguage:    "ko",
		SupportedLanguages: []string{"ko", "en"},
	}); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	languages := GetSupportedLanguages()

	if len(languages) != 2 {
		t.Errorf("GetSupportedLanguages() length = %d, want 2", len(languages))
	}

	hasKo := false
	hasEn := false
	for _, lang := range languages {
		if lang == "ko" {
			hasKo = true
		}
		if lang == "en" {
			hasEn = true
		}
	}

	if !hasKo {
		t.Error("GetSupportedLanguages() is missing ko")
	}
	if !hasEn {
		t.Error("GetSupportedLanguages() is missing en")
	}
}

func TestGetDefaultLanguage(t *testing.T) {
	Reset()
	if err := Init(&Config{
		DefaultLanguage:    "ko",
		SupportedLanguages: []string{"ko", "en"},
	}); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	defaultLang := GetDefaultLanguage()
	if defaultLang != "ko" {
		t.Errorf("GetDefaultLanguage() = %q, want %q", defaultLang, "ko")
	}
}

func TestTranslateBeforeInit(t *testing.T) {
	Reset()

	// Translate before Init must return the key itself
	result := Translate("ko", "auth.signup.success")
	if result != "auth.signup.success" {
		t.Errorf("Translate() before Init = %q, want %q", result, "auth.signup.success")
	}
}

func TestAllMessageKeysExist(t *testing.T) {
	Reset()
	if err := Init(&Config{
		DefaultLanguage:    "ko",
		SupportedLanguages: []string{"ko", "en"},
	}); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Key message IDs (error and success messages)
	messageKeys := []string{
		// Auth
		"auth.signup.email_required",
		"auth.signup.email_already_exists",
		"auth.signup.invalid_email",
		"auth.signup.password_required",
		"auth.signup.invalid_password",
		"auth.signup.name_required",
		"auth.signup.nickname_required",
		"auth.signup.invalid_nickname",
		"auth.signup.nickname_discriminator_full",
		"auth.signup.success",
		"auth.login.invalid_credentials",
		"auth.login.account_inactive",
		"auth.login.success",
		"auth.logout.failed",
		"auth.logout.success",
		"auth.token.required",
		"auth.token.invalid",
		"auth.token.expired",
		"auth.token.malformed",
		"auth.token.refresh_invalid",
		"auth.token.refresh_expired",
		"auth.refresh.success",
		"auth.withdraw.user_not_found",
		"auth.withdraw.password_mismatch",
		"auth.withdraw.already_withdrawn",
		"auth.withdraw.success",
		"auth.reauth.required",
		"auth.reauth.invalid",
		"auth.reauth.password_not_set",
		"auth.reauth.invalid_request",
		"auth.reauth.success",
		"auth.email.already_verified",
		"auth.email.verification_not_found",
		"auth.email.verification_expired",
		"auth.email.verification_invalid",
		"auth.email.send_failed",
		"auth.email.verification_sent",
		"auth.email.verified",
		"auth.password.reset_not_found",
		"auth.password.reset_expired",
		"auth.password.reset_invalid",
		"auth.password.reset_email_sent",
		"auth.password.reset_success",
		"auth.internal_error",

		// Health
		"health.check.success",
		"health.check.partial_failure",
	}

	languages := []string{"ko", "en"}

	for _, key := range messageKeys {
		for _, lang := range languages {
			result := Translate(lang, key)
			// A failed lookup returns the key itself
			if result == key {
				t.Errorf("[%s] no translation for message key %q", lang, key)
			}
		}
	}
}
