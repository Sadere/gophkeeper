package keeperv1

// Server stream mocks
//go:generate mockgen -destination mocks/mock_file.go -package keeperv1 github.com/Sadere/gophkeeper/pkg/proto/keeper/v1 SecretsService_UploadFileV1Server,SecretsService_DownloadFileV1Server

// Client mocks

//go:generate mockgen -destination mocks/mock_client_auth.go -package keeperv1 github.com/Sadere/gophkeeper/pkg/proto/keeper/v1 AuthServiceClient
//go:generate mockgen -source secrets_grpc.pb.go -destination mocks/mock_client_secrets.go -package keeperv1 SecretsServiceClient
//go:generate mockgen -source notification_grpc.pb.go -destination mocks/mock_client_notification.go -package keeperv1 NotificationServiceClient
