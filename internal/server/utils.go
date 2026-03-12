package server

import (
	"encoding/base32"
	"log"
	"strings"
)

func decodeBase32Label(label string) []byte {
	if label == "" {
		return nil
	}

	decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	encoded := strings.ToUpper(label)
	decoded, err := decoder.DecodeString(encoded)
	if err != nil {
		log.Printf("Failed to decode base32 label '%s': %v", label, err)
		return nil
	}

	return decoded
}

func extractPayload(qname string, domains []string) (domain string, payload []byte) {
	qname = strings.TrimSuffix(qname, ".")
	for _, d := range domains {
		d = strings.TrimSuffix(d, ".")
		if strings.HasSuffix(qname, d) {
			labels := strings.Split(qname, ".")
			decoded := decodeBase32Label(labels[0])
			return d, decoded
		}
	}

	log.Printf("Query %s does not match any domain", qname)
	return "", nil
}
