package nameserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type record struct {
	Domain     string `json:"domain"`
	RecordType string `json:"record_type"` // CNAME, A, MX etc...
	Value      []byte `json:"value"`       // IP address, CNAME value
	Ttl        int64  `json:"ttl"`
}

type recordStd struct {
	Domain     string `json:"domain"`
	RecordType string `json:"record_type"` // CNAME, A, MX etc...
	Value      string `json:"value"`       // IP address, CNAME value
	Ttl        int64  `json:"ttl"`
}

// AddRecordsFromJson
func AddRecordsFromJson(filePath string, logger *slog.Logger, dht dht) error {
	fp := filepath.Clean(filePath)

	rawJSON, err := os.ReadFile(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// file does not exist
			// no record import is needed
			logger.Info("no existing record file found")
			return nil
		}
		return fmt.Errorf("os.ReadFile: %w", err)
	}

	type x struct {
		Records []recordStd `json:"records"`
	}

	var rs x
	err = json.Unmarshal(rawJSON, &rs)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	for _, r := range rs.Records {
		key := fmt.Sprintf("%s:%s", r.Domain, r.RecordType)

		rr := &record{
			Domain:     r.Domain,
			RecordType: r.RecordType,
			Value:      []byte(r.Value),
			Ttl:        r.Ttl,
		}
		if err := dht.setValue(key, rr); err != nil {
			return fmt.Errorf("failed to set value %s %v", key, err)
		}
		logger.Info("import dns record", slog.Any("dns_record", r), slog.Any("record", rr))
	}

	return nil
}
