package check_privacy_policy

import "context"

type storage interface {
	PrivacyPolicySigned(ctx context.Context, userId int64) (bool, error)
}
