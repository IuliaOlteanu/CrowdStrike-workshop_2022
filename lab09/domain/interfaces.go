package domain

import "context"

type FilesMetadataStorage interface {
	SaveFile(ctx context.Context, metadata *FileMetadata) error
	RetrieveFile(ctx context.Context, sha256 string) (*FileMetadata, error)
}

type BookMetadataStorage interface {
	SaveBook(ctx context.Context, metadata *BookMetadata) error
	RetrieveBook(ctx context.Context, sha256 string) (*BookMetadata, error)
}
