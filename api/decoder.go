package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/schema"
)

var (
	schemaDecoder *schema.Decoder
	once          sync.Once
)

func getSchemaDecoder() *schema.Decoder {
	once.Do(func() {
		schemaDecoder = schema.NewDecoder()
		schemaDecoder.SetAliasTag("json")
		schemaDecoder.IgnoreUnknownKeys(true)
	})
	return schemaDecoder
}

type Decoder struct {
	r        *http.Request
	bodyErr  error
	queryErr error
}

func NewDecoder(r *http.Request) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Body(target any) *Decoder {
	if d.bodyErr != nil || d.r.Body == nil {
		return d
	}

	err := json.NewDecoder(d.r.Body).Decode(target)
	if err != nil && !errors.Is(err, io.EOF) {
		d.bodyErr = err
	}
	return d
}

func (d *Decoder) Query(target any) *Decoder {
	if d.queryErr != nil {
		return d
	}
	d.queryErr = getSchemaDecoder().Decode(target, d.r.URL.Query())
	return d
}

func (d *Decoder) Err() error {
	if d.bodyErr != nil {
		return d.bodyErr
	}
	return d.queryErr
}

func (d *Decoder) HasError() bool {
	return d.bodyErr != nil || d.queryErr != nil
}
