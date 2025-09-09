package ctl

import "rtc/internal/ctl/client"

func convertConfigsToReq(configs upsertConfigYaml) []*client.UpsertConfigRequest {
	result := make([]*client.UpsertConfigRequest, 0, len(configs))
	for key, options := range configs {
		result = append(result, &client.UpsertConfigRequest{
			Key:       key,
			Value:     options.Value,
			ValueType: options.Type,
			Usage:     options.Usage,
			Group:     options.Group,
			Writable:  options.Writable,
		})
	}

	return result
}
