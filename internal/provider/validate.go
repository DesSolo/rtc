package provider

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/DesSolo/rtc/internal/models"
)

var valueTypesValidators = map[models.ValueType]func(string) error{
	models.ValueTypeUnknown: func(_ string) error {
		return errors.New("unknown value type")
	},
	models.ValueTypeString: func(_ string) error {
		// sting value is always valid
		return nil
	},
	models.ValueTypeBool: func(value string) error {
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("is not a bool: %s", value)
		}

		return nil
	},
	models.ValueTypeInt: func(value string) error {
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("is not a int: %s", value)
		}

		return nil
	},
	models.ValueTypeInt64: func(value string) error {
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Errorf("is not a int64: %s", value)
		}

		return nil
	},

	models.ValueTypeUint: func(value string) error {
		if _, err := strconv.ParseUint(value, 10, 64); err != nil {
			return fmt.Errorf("is not a uint: %s", value)
		}

		return nil
	},

	models.ValueTypeUint64: func(value string) error {
		if _, err := strconv.ParseUint(value, 10, 64); err != nil {
			return fmt.Errorf("is not a uint64: %s", value)
		}

		return nil
	},
	models.ValueTypeFloat: func(value string) error {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("is not a float: %s", value)
		}

		return nil
	},
	models.ValueTypeFloat64: func(value string) error {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("is not a float64: %s", value)
		}

		return nil
	},
}

func validateNewValue(config *models.Config, newValue []byte) error {
	if !config.Metadata.Writable {
		return fmt.Errorf("%w: value is not writable", ErrNotValid)
	}

	if err := validateByValueType(config.ValueType, newValue); err != nil {
		return fmt.Errorf("%w: value type: %w", ErrNotValid, err)
	}

	// TODO: view validate

	return nil
}

func validateUpsert(configs []*models.Config) error {
	// etcd maximum items in one transaction
	// TODO: add chunked wrapper
	if len(configs) > 128 {
		return fmt.Errorf("%w: to many configs", ErrNotValid)
	}

	for _, config := range configs {
		if err := validateByValueType(config.ValueType, config.Value); err != nil {
			return fmt.Errorf("%w: value type: %w", ErrNotValid, err)
		}

		// TODO: view validate
	}

	return nil
}

func validateByValueType(valueType models.ValueType, value []byte) error {
	validator, ok := valueTypesValidators[valueType]
	if !ok {
		return fmt.Errorf("unsupported: %s", valueType)
	}

	if err := validator(string(value)); err != nil {
		return fmt.Errorf("validator: %w", err)
	}

	return nil
}
