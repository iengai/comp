package domain

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	otpExpirationSeconds = 300
)

var (
	ErrOTPExpired = errors.New("OTP expired")
)

type (
	Code string
	OTP  struct {
		ID        uuid.UUID
		ExpiredAt time.Time
		Code      Code
		Used      bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func NewCode(now time.Time) Code {
	r := rand.New(rand.NewSource(now.UnixNano()))
	return Code(fmt.Sprintf("%06d", r.Intn(1000000)))
}

func (c Code) ToString() string {
	return string(c)
}

func NewOTP(id uuid.UUID, now time.Time) *OTP {
	return &OTP{
		ID:        id,
		ExpiredAt: now.Add(time.Duration(otpExpirationSeconds) * time.Second),
		Code:      NewCode(now),
		Used:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (o *OTP) Expired(now time.Time) bool {
	return o.ExpiredAt.After(now)
}
