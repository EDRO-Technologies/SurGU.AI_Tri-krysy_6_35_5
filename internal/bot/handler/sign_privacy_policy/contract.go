package sign_privacy_policy

import "context"

type storage interface {
	SignPrivacyPolicy(ctx context.Context, userId int64) error
}
