package api

type GetFileResponse struct {
	Content []byte
}

func (r *GetFileResponse) Write(p []byte) (n int, err error) {
	r.Content = append(r.Content, p...)

	return len(p), nil
}
