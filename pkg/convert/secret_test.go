package convert

import (
	"reflect"
	"testing"
	"time"

	"github.com/Sadere/gophkeeper/pkg/model"
	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSecretToProto(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		secret *model.Secret
		want   *pb.Secret
	}{
		{
			name: "cred",
			secret: &model.Secret{
				ID:        111,
				CreatedAt: now,
				UpdatedAt: now,
				Metadata:  "metadata",
				SType:     string(model.CredSecret),
				Creds: &model.Credentials{
					Login:    "login",
					Password: "password",
				},
			},
			want: &pb.Secret{
				Id:        111,
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
				Metadata:  "metadata",
				Type:      pb.SecretType_SECRET_TYPE_CREDENTIAL,
				Content: &pb.Secret_Credential{
					Credential: &pb.Credential{
						Login:    "login",
						Password: "password",
					},
				},
			},
		},
		{
			name: "text",
			secret: &model.Secret{
				ID:        111,
				CreatedAt: now,
				UpdatedAt: now,
				Metadata:  "metadata",
				SType:     string(model.TextSecret),
				Text: &model.Text{
					Content: "Text content",
				},
			},
			want: &pb.Secret{
				Id:        111,
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
				Metadata:  "metadata",
				Type:      pb.SecretType_SECRET_TYPE_TEXT,
				Content: &pb.Secret_Text{
					Text: &pb.Text{
						Text: "Text content",
					},
				},
			},
		},
		{
			name: "blob",
			secret: &model.Secret{
				ID:        111,
				CreatedAt: now,
				UpdatedAt: now,
				Metadata:  "metadata",
				SType:     string(model.BlobSecret),
				Blob: &model.Blob{
					FileName: "file.txt",
				},
			},
			want: &pb.Secret{
				Id:        111,
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
				Metadata:  "metadata",
				Type:      pb.SecretType_SECRET_TYPE_BLOB,
				Content: &pb.Secret_Blob{
					Blob: &pb.Blob{
						FileName: "file.txt",
					},
				},
			},
		},
		{
			name: "card",
			secret: &model.Secret{
				ID:        111,
				CreatedAt: now,
				UpdatedAt: now,
				Metadata:  "metadata",
				SType:     string(model.CardSecret),
				Card: &model.Card{
					Number:   "111",
					ExpYear:  2014,
					ExpMonth: 12,
					Cvv:      444,
				},
			},
			want: &pb.Secret{
				Id:        111,
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
				Metadata:  "metadata",
				Type:      pb.SecretType_SECRET_TYPE_CARD,
				Content: &pb.Secret_Card{
					Card: &pb.Card{
						Number:   "111",
						ExpYear:  2014,
						ExpMonth: 12,
						Cvv:      444,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualSecret := SecretToProto(tt.secret)

			if !reflect.DeepEqual(tt.want, actualSecret) {
				t.Errorf("unexpected proto secret want = %v got = %v", tt.want, actualSecret)
			}
		})
	}
}

func TestProtoToSecret(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name   string
		secret *pb.Secret
		want   *model.Secret
	}{
		{
			name: "cred",
			secret: &pb.Secret{
				Id:        111,
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
				Metadata:  "metadata",
				Type:      pb.SecretType_SECRET_TYPE_CREDENTIAL,
				Content: &pb.Secret_Credential{
					Credential: &pb.Credential{
						Login:    "login",
						Password: "password",
					},
				},
			},
			want: &model.Secret{
				ID:        111,
				CreatedAt: now,
				UpdatedAt: now,
				Metadata:  "metadata",
				SType:     string(model.CredSecret),
				Creds: &model.Credentials{
					Login:    "login",
					Password: "password",
				},
			},
		},
		{
			name: "text",
			secret: &pb.Secret{
				Id:        111,
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
				Metadata:  "metadata",
				Type:      pb.SecretType_SECRET_TYPE_TEXT,
				Content: &pb.Secret_Text{
					Text: &pb.Text{
						Text: "Text content",
					},
				},
			},
			want: &model.Secret{
				ID:        111,
				CreatedAt: now,
				UpdatedAt: now,
				Metadata:  "metadata",
				SType:     string(model.TextSecret),
				Text: &model.Text{
					Content: "Text content",
				},
			},
		},
		{
			name: "blob",
			secret: &pb.Secret{
				Id:        111,
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
				Metadata:  "metadata",
				Type:      pb.SecretType_SECRET_TYPE_BLOB,
				Content: &pb.Secret_Blob{
					Blob: &pb.Blob{
						FileName: "file.txt",
					},
				},
			},
			want: &model.Secret{
				ID:        111,
				CreatedAt: now,
				UpdatedAt: now,
				Metadata:  "metadata",
				SType:     string(model.BlobSecret),
				Blob: &model.Blob{
					FileName: "file.txt",
				},
			},
		},
		{
			name: "card",
			secret: &pb.Secret{
				Id:        111,
				CreatedAt: timestamppb.New(now),
				UpdatedAt: timestamppb.New(now),
				Metadata:  "metadata",
				Type:      pb.SecretType_SECRET_TYPE_CARD,
				Content: &pb.Secret_Card{
					Card: &pb.Card{
						Number:   "111",
						ExpYear:  2014,
						ExpMonth: 12,
						Cvv:      444,
					},
				},
			},
			want: &model.Secret{
				ID:        111,
				CreatedAt: now,
				UpdatedAt: now,
				Metadata:  "metadata",
				SType:     string(model.CardSecret),
				Card: &model.Card{
					Number:   "111",
					ExpYear:  2014,
					ExpMonth: 12,
					Cvv:      444,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualSecret := ProtoToSecret(tt.secret)

			if !reflect.DeepEqual(tt.want, actualSecret) {
				t.Errorf("unexpected secret want = %v got = %v", tt.want, actualSecret)
			}
		})
	}
}
