package scalars

import (
	"encoding/json"
	"fmt"
	"io"
)

type JSON json.RawMessage

func (j *JSON) UnmarshalGQL(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("scalars.JSON: cannot marshal value: %w", err)
	}
	*j = JSON(b)
	return nil
}

func (j JSON) MarshalGQL(w io.Writer) {
	_, _ = w.Write(j)
}
