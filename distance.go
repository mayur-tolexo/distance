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
func GetPinDistanct(source string, distination []string) (distance []float64, err error) {

	var (
		long, lat float64
	)
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
			// fmt.Println(sourceVal, destVal)
			distance, err = getDistance(client, sourceVal, destVal)
		}
	}
	return
}

func getDistance(client *http.Client, sourceVal string, destVal string) (distance []float64, err error) {
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
				var data Matrix
				if err = jsoniter.Unmarshal(bodyBytes, &data); err == nil {
					if len(data.Res) > 0 {
						resp := data.Res[0]
						if len(resp.Res) > 0 {
							resource := resp.Res[0]
							if val, ok := resource["results"].([]interface{}); ok {
								if len(val) > 0 {
									for _, curVal := range val {
										if dis, ok := curVal.(map[string]interface{}); ok {
											distance = append(distance, FloatRound(IfToF(dis["travelDistance"]), 4))
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

//FloatRound : Float round off
func FloatRound(v float64, decimals int) float64 {
	var pow float64 = 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int((v*pow)+0.5)) / pow
}
