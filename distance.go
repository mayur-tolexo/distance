package distance

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/mayur-tolexo/distance/constant"
)

//GetPinDistanct will return distance between pincodes
func GetPinDistanct(source string, distination []string) (err error) {

	var long, lat float64
	client := &http.Client{}
	if long, lat, err = getCordinates(client, source); err == nil {
		sourceVal := fmt.Sprintf("%v,%v", long, lat)
		destVal := ""
		for _, curPin := range distination {
			if long, lat, err = getCordinates(client, curPin); err == nil {
				destVal += fmt.Sprintf("%v,%v;", long, lat)
			} else {
				break
			}
		}
		if err == nil {
			destVal = strings.TrimSuffix(destVal, ";")
			fmt.Println(sourceVal, destVal)
			getDistance(client, sourceVal, destVal)
		}
	}
	return
}

func getDistance(client *http.Client, sourceVal string, destVal string) (err error) {
	var (
		req       *http.Request
		resp      *http.Response
		bodyBytes []byte
	)
	url := fmt.Sprintf("%v?origins=%v&destinations=%v&travelMode=driving&key=%v", constant.DistanceBase, sourceVal, destVal, constant.Key)
	if req, err = http.NewRequest("GET", url, nil); err == nil {
		if resp, err = client.Do(req); err == nil {
			defer resp.Body.Close()
			if bodyBytes, err = ioutil.ReadAll(resp.Body); err == nil {
				fmt.Println(string(bodyBytes))
			}
		}
	}
	return
}

func getCordinates(client *http.Client, pincode string) (long, lat float64, err error) {
	var (
		req       *http.Request
		resp      *http.Response
		bodyBytes []byte
	)
	url := fmt.Sprintf("%v?postalCode=%v&key=%v", constant.LocationBase, pincode, constant.Key)
	if req, err = http.NewRequest("GET", url, nil); err == nil {
		if resp, err = client.Do(req); err == nil {
			defer resp.Body.Close()
			if bodyBytes, err = ioutil.ReadAll(resp.Body); err == nil {
				var data map[string]interface{}
				if err = jsoniter.Unmarshal(bodyBytes, &data); err == nil {
					if set, exists := data["resourceSets"]; exists {
						if setVal, ok := set.([]interface{}); ok {
							if len(setVal) > 0 {
								if val, ok := setVal[0].(map[string]interface{}); ok {
									if resource, exists := val["resources"]; exists {
										if rVal, ok := resource.([]interface{}); ok {
											if len(rVal) > 0 {
												if val, ok := rVal[0].(map[string]interface{}); ok {
													if point, exists := val["point"]; exists {
														if pVal, ok := point.(map[string]interface{}); ok {
															if cord, exists := pVal["coordinates"]; exists {
																if cVal, ok := cord.([]interface{}); ok {
																	long = IfToF(cVal[0])
																	lat = IfToF(cVal[1])
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

//IfToF converts interface value to float
func IfToF(value interface{}) float64 {
	var floatVal float64
	objKind := reflect.TypeOf(value).Kind()
	switch objKind {
	case reflect.Int:
		floatVal = float64(value.(int))
	case reflect.Int64:
		floatVal = float64(value.(int64))
	case reflect.Float64:
		floatVal = value.(float64)
	case reflect.Float32:
		floatVal = float64(value.(float32))
	case reflect.String:
		floatVal, _ = strconv.ParseFloat(value.(string), 64)
	}
	return floatVal
}
