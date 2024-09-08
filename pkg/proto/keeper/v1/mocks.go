package keeperv1

//go:generate mockgen -destination mocks/mock_file.go -package keeperv1 github.com/Sadere/gophkeeper/pkg/proto/keeper/v1 SecretsService_UploadFileV1Server,SecretsService_DownloadFileV1Server
