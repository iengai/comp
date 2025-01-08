package infra

import (
	"context"

	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/iengai/comp/user/internal"
	"github.com/iengai/comp/user/internal/domain"
)

const (
	otpTTL = 24 * time.Hour
)

type (
	OTPRepository interface {
		Get(ctx context.Context, id uuid.UUID) (*domain.OTP, error)
		Save(ctx context.Context, otp *domain.OTP) error
	}

	otpRepository struct {
		db    *dynamodb.Client
		table string
	}

	OTP struct {
		PK        string    `dynamodbav:"PK"`
		SK        string    `dynamodbav:"SK"`
		ExpiredAt time.Time `dynamodbav:"expired_at"`
		Code      string    `dynamodbav:"code"`
		Used      bool      `dynamodbav:"used"`
		CreatedAt time.Time `dynamodbav:"created_at"`
		UpdatedAt time.Time `dynamodbav:"updated_at"`
		TTL       int64     `dynamodbav:"ttl"`
	}
)

func (o OTP) IsExpired(now time.Time) bool {
	return o.ExpiredAt.After(now)
}

func (o otpRepository) Get(ctx context.Context, id uuid.UUID) (*domain.OTP, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(o.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: o.pk()},
			"SK": &types.AttributeValueMemberS{Value: o.newSK(id)},
		},
	}
	res, err := o.db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}
	if res.Item == nil {
		return nil, internal.ErrNotFound
	}
	var e OTP
	err = attributevalue.UnmarshalMap(res.Item, &e)
	if err != nil {
		return nil, err
	}
	return o.toDomain(&e), nil
}

func (o otpRepository) toDomain(data *OTP) *domain.OTP {
	return &domain.OTP{
		ID:        uuid.MustParse(data.SK),
		ExpiredAt: data.ExpiredAt,
		Code:      domain.Code(data.Code),
		Used:      data.Used,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func (o otpRepository) Save(ctx context.Context, otp *domain.OTP) error {
	data := o.toData(otp)
	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		return err
	}
	_, err = o.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &o.table,
		Item:      item,
	})
	return err
}

func (o otpRepository) toData(d *domain.OTP) *OTP {
	return &OTP{
		PK:        o.pk(),
		SK:        o.newSK(d.ID),
		ExpiredAt: d.ExpiredAt,
		Code:      string(d.Code),
		Used:      d.Used,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		TTL:       d.ExpiredAt.Add(otpTTL).Unix(),
	}
}

func (o otpRepository) newSK(id uuid.UUID) string {
	return id.String()
}

func (o otpRepository) pk() string {
	return "otp"
}

func NewOTPRepository(table string, db *dynamodb.Client) OTPRepository {
	return &otpRepository{
		table: table,
		db:    db,
	}
}
