package domain

type FileMetadata struct {
	Sha256          string `json:"sha256"`
	FileName        string `json:"file_name"`
	CompilerVersion int    `json:"compiler_version"`
	Size            int64  `json:"size"`
}

type BookMetadata struct {
	Sha256          string `json:"sha256"`
	Title			string `json:title`
	Author			string `json:author`

}