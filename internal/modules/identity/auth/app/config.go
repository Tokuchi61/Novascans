package app

import "time"

type ServiceConfig struct {
	AppEnv                    string
	AccessTokenSecret         string
	AccessTokenTTL            time.Duration
	RefreshTokenTTL           time.Duration
	EmailVerificationTokenTTL time.Duration
	PasswordResetTokenTTL     time.Duration
}

func (cfg ServiceConfig) ExposeDebugTokens() bool {
	return cfg.AppEnv != "production"
}
