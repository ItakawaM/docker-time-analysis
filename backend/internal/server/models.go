package server

type UploadResponse struct {
	Message    string `json:"message"`
	ParsedRows int    `json:"parsed_rows"`
}
