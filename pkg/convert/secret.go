// Provides functions for converting between protobuf models and regular models
package convert

import (
	"github.com/Sadere/gophkeeper/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

// Returns corresponding secret type
func ProtoToType(pbType pb.SecretType) model.SecretType {
	switch pbType {
	case pb.SecretType_SECRET_TYPE_CREDENTIAL:
		return model.CredSecret
	case pb.SecretType_SECRET_TYPE_TEXT:
		return model.TextSecret
	case pb.SecretType_SECRET_TYPE_BLOB:
		return model.BlobSecret
	case pb.SecretType_SECRET_TYPE_CARD:
		return model.CardSecret
	default:
		return model.UnknownSecret
	}
}

// Returns corresponding protobuf secret type
func TypeToProto(sType string) pb.SecretType {
	switch sType {
	case string(model.CredSecret):
		return pb.SecretType_SECRET_TYPE_CREDENTIAL
	case string(model.TextSecret):
		return pb.SecretType_SECRET_TYPE_TEXT
	case string(model.BlobSecret):
		return pb.SecretType_SECRET_TYPE_BLOB
	case string(model.CardSecret):
		return pb.SecretType_SECRET_TYPE_CARD
	default:
		return pb.SecretType_SECRET_TYPE_UNSPECIFIED
	}
}

// Converts secret model to protobuf counterpart
func SecretToProto(secret *model.Secret) *pb.Secret {
	pbSecret := &pb.Secret{
		Id:        secret.ID,
		CreatedAt: timestamppb.New(secret.CreatedAt),
		UpdatedAt: timestamppb.New(secret.UpdatedAt),
		Metadata:  secret.Metadata,
		Type:      pb.SecretType_SECRET_TYPE_UNSPECIFIED,
	}

	pbSecret.Type = TypeToProto(secret.SType)

	switch secret.SType {
	case string(model.CredSecret):
		pbSecret.Content = &pb.Secret_Credential{
			Credential: &pb.Credential{
				Login:    secret.Creds.Login,
				Password: secret.Creds.Password,
			},
		}
	case string(model.TextSecret):
		pbSecret.Content = &pb.Secret_Text{
			Text: &pb.Text{
				Text: secret.Text.Content,
			},
		}
	case string(model.BlobSecret):
		pbSecret.Content = &pb.Secret_Blob{
			Blob: &pb.Blob{
				FileName: secret.Blob.FileName,
			},
		}
	case string(model.CardSecret):
		pbSecret.Content = &pb.Secret_Card{
			Card: &pb.Card{
				Number:   secret.Card.Number,
				ExpYear:  secret.Card.ExpYear,
				ExpMonth: secret.Card.ExpMonth,
				Cvv:      secret.Card.Cvv,
			},
		}
	}

	return pbSecret
}

// Converts protobuf model to regular model
func ProtoToSecret(pbSecret *pb.Secret) *model.Secret {
	secret := &model.Secret{
		ID:        pbSecret.Id,
		CreatedAt: pbSecret.CreatedAt.AsTime(),
		UpdatedAt: pbSecret.UpdatedAt.AsTime(),
		Metadata:  pbSecret.Metadata,
		SType:     string(model.UnknownSecret),
	}

	secret.SType = string(ProtoToType(pbSecret.Type))

	switch secret.SType {
	case string(model.CredSecret):
		pbCred := pbSecret.Content.(*pb.Secret_Credential)
		secret.Creds = &model.Credentials{
			Login:    pbCred.Credential.GetLogin(),
			Password: pbCred.Credential.GetPassword(),
		}
	case string(model.TextSecret):
		pbText := pbSecret.Content.(*pb.Secret_Text)
		secret.Text = &model.Text{
			Content: pbText.Text.GetText(),
		}
	case string(model.BlobSecret):
		pbBlob := pbSecret.Content.(*pb.Secret_Blob)
		secret.Blob = &model.Blob{
			FileName: pbBlob.Blob.GetFileName(),
		}
	case string(model.CardSecret):
		pbCard := pbSecret.Content.(*pb.Secret_Card)
		secret.Card = &model.Card{
			Number:   pbCard.Card.GetNumber(),
			ExpYear:  pbCard.Card.GetExpYear(),
			ExpMonth: pbCard.Card.GetExpMonth(),
			Cvv:      pbCard.Card.GetCvv(),
		}
	}

	return secret
}

// Converts protobuf secret preview to regular secret preview
func ProtoToPreview(pbPreview *pb.SecretPreview) *model.SecretPreview {
	return &model.SecretPreview{
		ID:        pbPreview.Id,
		CreatedAt: pbPreview.CreatedAt.AsTime(),
		UpdatedAt: pbPreview.UpdatedAt.AsTime(),
		Metadata:  pbPreview.Metadata,
		SType:     string(ProtoToType(pbPreview.Type)),
	}
}

// Converts secret preview to protobuf struct
func PreviewToProto(preview *model.SecretPreview) *pb.SecretPreview {
	return &pb.SecretPreview{
		Id:        preview.ID,
		CreatedAt: timestamppb.New(preview.CreatedAt),
		UpdatedAt: timestamppb.New(preview.UpdatedAt),
		Metadata:  preview.Metadata,
		Type:      TypeToProto(preview.SType),
	}
}
