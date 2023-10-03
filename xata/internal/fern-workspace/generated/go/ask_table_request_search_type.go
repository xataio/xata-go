// This file was auto-generated by Fern from our API Definition.

package api

import (
	json "encoding/json"
	fmt "fmt"
	strconv "strconv"
)

// The type of search to use. If set to `keyword` (the default), the search can be configured by passing
// a `search` object with the following fields. For more details about each, see the Search endpoint documentation.
// All fields are optional.
//   - fuzziness  - typo tolerance
//   - target - columns to search into, and weights.
//   - prefix - prefix search type.
//   - filter - pre-filter before searching.
//   - boosters - control relevancy.
//
// If set to `vector`, a `vectorSearch` object must be passed, with the following parameters. For more details, see the Vector
// Search endpoint documentation. The `column` and `contentColumn` parameters are required.
//   - column - the vector column containing the embeddings.
//   - contentColumn - the column that contains the text from which the embeddings where computed.
//   - filter - pre-filter before searching.
type AskTableRequestSearchType uint8

const (
	AskTableRequestSearchTypeKeyword AskTableRequestSearchType = iota + 1
	AskTableRequestSearchTypeVector
)

func (a AskTableRequestSearchType) String() string {
	switch a {
	default:
		return strconv.Itoa(int(a))
	case AskTableRequestSearchTypeKeyword:
		return "keyword"
	case AskTableRequestSearchTypeVector:
		return "vector"
	}
}

func (a AskTableRequestSearchType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", a.String())), nil
}

func (a *AskTableRequestSearchType) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch raw {
	case "keyword":
		value := AskTableRequestSearchTypeKeyword
		*a = value
	case "vector":
		value := AskTableRequestSearchTypeVector
		*a = value
	}
	return nil
}
