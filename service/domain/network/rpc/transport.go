package rpc

import (
	"encoding/json"
	"fmt"
)

type RequestBody struct {
	Name []string        `json:"name"`
	Type string          `json:"type"`
	Args json.RawMessage `json:"args"`
}

const (
	transportStringForProcedureTypeSource = "source"
	transportStringForProcedureTypeDuplex = "duplex"
	transportStringForProcedureTypeAsync  = "async"
)

func decodeProcedureType(str string) ProcedureType {
	switch str {
	case transportStringForProcedureTypeSource:
		return ProcedureTypeSource
	case transportStringForProcedureTypeDuplex:
		return ProcedureTypeDuplex
	case transportStringForProcedureTypeAsync:
		return ProcedureTypeAsync
	default:
		return ProcedureTypeUnknown
	}
}

func encodeProcedureType(t ProcedureType) (string, error) {
	switch t {
	case ProcedureTypeSource:
		return transportStringForProcedureTypeSource, nil
	case ProcedureTypeDuplex:
		return transportStringForProcedureTypeDuplex, nil
	case ProcedureTypeAsync:
		return transportStringForProcedureTypeAsync, nil
	default:
		return "", fmt.Errorf("unknown procedure type %T", t)
	}
}