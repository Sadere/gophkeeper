package keeperv1

import (
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// Server stream mocks
//go:generate mockgen -destination mocks/mock_file.go -package keeperv1 github.com/Sadere/gophkeeper/pkg/proto/keeper/v1 SecretsService_UploadFileV1Server,SecretsService_DownloadFileV1Server

// Client mocks

//go:generate mockgen -destination mocks/mock_client_auth.go -package keeperv1 github.com/Sadere/gophkeeper/pkg/proto/keeper/v1 AuthServiceClient
//go:generate mockgen -source secrets_grpc.pb.go -destination mocks/mock_client_secrets.go -package keeperv1 SecretsServiceClient
//go:generate mockgen -source notification_grpc.pb.go -destination mocks/mock_client_notification.go -package keeperv1 NotificationServiceClient

// Client stream mocks

//go:generate mockgen -destination mocks/mock_client_stream.go -package keeperv1 github.com/Sadere/gophkeeper/pkg/proto/keeper/v1 UploadFile_ClientStream,DownloadFile_ClientStream

type UploadFile_ClientStream interface {
	grpc.ClientStreamingClient[UploadFileRequestV1, emptypb.Empty]
}

type DownloadFile_ClientStream interface {
	grpc.ServerStreamingClient[DownloadFileResponseV1]
}