package logfmt

import (
	"sort"

	"github.com/alnvdl/terr"
	gologfmt "github.com/go-logfmt/logfmt"
)

type JSONFmt struct{}

func (j JSONFmt) ToLogfmt(data map[string]any) ([]byte, error) {
	var keyvals []any

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		keyvals = append(keyvals, k, data[k])
	}

	out, err := gologfmt.MarshalKeyvals(keyvals...)
	if err != nil {
		return nil, terr.Trace(err)
	}
	return out, nil
}
