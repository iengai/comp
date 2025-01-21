package domain

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	otpExpiration      = 300 * time.Second
	otpResendInterval  = 60 * time.Second
	otpLockoutDuration = 600 * time.Second

	otpMaxAttempts = 5
)

var (
	ErrOTPExpired         = errors.New("OTP expired")
	ErrOTPTooManyAttempts = errors.New("OTP too many attempts")
	ErrOTPInvalidCode     = errors.New("OTP invalid")
	ErrOTPCodeUsed        = errors.New("OTP used")

	ErrOTPResendNeedInterval = errors.New("OTP resend needs interval")
)

type (
	Code string
	OTP  struct {
		ID             uuid.UUID
		ExpiredAt      time.Time
		Code           Code
		Used           bool
		FailedAttempts int
		LockoutAt      time.Time
		SentAt         time.Time
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	OTPRepository interface {
		Get(ctx context.Context, id uuid.UUID) (*OTP, error)
		Save(ctx context.Context, otp *OTP) error
	}
)

func NewCode(now time.Time) Code {
	r := rand.New(rand.NewSource(now.UnixNano()))
	return Code(fmt.Sprintf("%06d", r.Intn(1000000)))
}

func (c Code) ToString() string {
	return string(c)
}

func NewOTPForSending(id uuid.UUID, now time.Time) *OTP {
	return &OTP{
		ID:        id,
		ExpiredAt: now.Add(otpExpiration),
		SentAt:    now,
		Code:      NewCode(now),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (o *OTP) Resend(now time.Time) error {
	if now.Before(o.SentAt.Add(otpResendInterval)) {
		return ErrOTPResendNeedInterval
	}
	if o.isLockedOut(now) {
		return ErrOTPTooManyAttempts
	}
	o.Code = NewCode(now)
	o.ExpiredAt = now.Add(otpExpiration)
	o.SentAt = now
	o.Used = false
	o.FailedAttempts = 0
	o.UpdatedAt = now
	return nil
}

func (o *OTP) Verify(code Code, now time.Time) error {
	if o.isExpired(now) {
		return ErrOTPExpired
	}
	if o.isLockedOut(now) {
		return ErrOTPTooManyAttempts
	}
	if o.Used {
		return ErrOTPCodeUsed
	}
	if code != o.Code {
		o.FailedAttempts++
		if o.FailedAttempts >= otpMaxAttempts {
			o.lockOut(now)
		}
		return ErrOTPInvalidCode
	}
	o.Used = true
	o.FailedAttempts = 0
	o.UpdatedAt = now
	return nil
}

func (o *OTP) isExpired(now time.Time) bool {
	return o.ExpiredAt.After(now)
}

func (o *OTP) lockOut(now time.Time) {
	if o.isLockedOut(now) {
		return
	}
	o.LockoutAt = now
	o.FailedAttempts = 0
}

func (o *OTP) isLockedOut(now time.Time) bool {
	return now.After(o.LockoutAt.Add(otpLockoutDuration))
}
