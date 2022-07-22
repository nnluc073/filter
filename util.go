package filter

import (
	"net/url"

	"github.com/nnluc073/filter/filter"
	"github.com/nnluc073/filter/models"

	"github.com/valyala/fasthttp"
)

func fMapping(values []string) interface{} {
	var fData = make([]*filter.Filter, len(values))
	for k, v := range values {
		if data, err := filter.ParseFilter(v); err == nil {
			fData[k] = data

		}
	}
	return fData
}

func jMapping(values []string) interface{} {
	var fData = make([]*filter.Join, len(values))
	for k, v := range values {
		if data, err := filter.ParseJoin(v); err == nil {
			fData[k] = data
		}
	}
	return fData
}

func sMapping(values []string) interface{} {
	var fData = make([]*filter.Sort, len(values))
	for k, v := range values {
		if data, err := filter.ParseSort(v); err == nil {
			fData[k] = data
		}
	}
	return fData
}

func flatten(dst map[string]interface{}, values url.Values) {
	for field, value := range values {
		switch field {
		case "filter":
			dst[field] = fMapping(value)
		case "or":
			dst[field] = fMapping(value)
		case "join":
			dst[field] = jMapping(value)
		case "sort":
			dst[field] = sMapping(value)
		default:
			if len(values) > 1 {
				dst[field] = values
			} else {
				dst[field] = value[0]
			}
		}
	}

}

func ParseQuery(request *fasthttp.Request) (*models.FilterReq, error) {

	filterData := make(map[string]interface{})

	queryParams, err := url.ParseQuery(string(request.URI().QueryString()))
	if err == nil {
		flatten(filterData, queryParams)
		return &models.FilterReq{Data: filterData}, nil
	}
	return nil, err
}
