package convert

import (
	"github.com/Sadere/gophkeeper/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

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

func SecretToProto(secret *model.Secret) *pb.Secret {
	pbSecret := &pb.Secret{
		Id:        secret.ID,
		CreatedAt: timestamppb.New(secret.CreatedAt),
		UpdatedAt: timestamppb.New(secret.UpdatedAt),
		Metadata:  secret.Metadata,
		Type:      pb.SecretType_SECRET_TYPE_UNSPECIFIED,
	}

	pbSecret.Type = TypeToProto(secret.SType)

	return pbSecret
}

func ProtoToPreview(pbPreview *pb.SecretPreview) *model.SecretPreview {
	return &model.SecretPreview{
		ID:        pbPreview.Id,
		CreatedAt: pbPreview.CreatedAt.AsTime(),
		UpdatedAt: pbPreview.UpdatedAt.AsTime(),
		Metadata:  pbPreview.Metadata,
		SType:     string(ProtoToType(pbPreview.Type)),
	}
}

func PreviewToProto(preview *model.SecretPreview) *pb.SecretPreview {
	return &pb.SecretPreview{
		Id:        preview.ID,
		CreatedAt: timestamppb.New(preview.CreatedAt),
		UpdatedAt: timestamppb.New(preview.UpdatedAt),
		Metadata:  preview.Metadata,
		Type:      TypeToProto(preview.SType),
	}
}
