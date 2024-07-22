package customtypes

import (
	"encoding/json"
	"strconv"
)

type Int64Slice []int64

func (slice Int64Slice) MarshalJSON() ([]byte, error) {
	strSlice := make([]string, len(slice))
	for i, intValue := range slice {
		strSlice[i] = strconv.FormatInt(intValue, 10)
	}

	return json.Marshal(strSlice)
}

func (slice *Int64Slice) UnmarshalJSON(data []byte) error {
	var rawSlice []json.RawMessage
	err := json.Unmarshal(data, &rawSlice)
	if err != nil {
		return err
	}

	intSlice := make([]int64, len(rawSlice))
	for i, raw := range rawSlice {
		var strValue string
		err = json.Unmarshal(raw, &strValue)
		if err != nil {
			return err
		}

		intValue, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return err
		}

		intSlice[i] = intValue
	}

	*slice = Int64Slice(intSlice)
	return nil
}
