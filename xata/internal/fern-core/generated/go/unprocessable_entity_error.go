// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"

	core "github.com/xataio/xata-go/xata/internal/fern-core/generated/go/core"
)

type UnprocessableEntityError struct {
	*core.APIError
	Body any
}

func (u *UnprocessableEntityError) UnmarshalJSON(data []byte) error {
	var body any
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	u.StatusCode = 422
	u.Body = body
	return nil
}

func (u *UnprocessableEntityError) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Body)
}
