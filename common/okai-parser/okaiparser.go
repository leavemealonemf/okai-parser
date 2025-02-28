package okaiparser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// +RESP:GTFRI,OK043A,868070043228349,zk200,,,,,0,0000000000000000000,,,,,,,,,0250,0099,04E9,08C41A65,26&99,2,41,0,52322,4022,87,0,,0,,0.0&0.00&42.50&263.13&1&1&0&0000000D0000000D011A0000&000000D10500000000060000&00000000&0&02641C1B1A1AFFFFFF7D&1&1&00000000000000,85,20250228065940,00A6$

func ParseParams(params []string) (string, error) {
	if len(params) == 0 {
		return "", fmt.Errorf("Parse params failed. Zero len")
	}

	pktType, pktId := HeadInfo(params[0])

	if pktType == "+RESP" && pktId == "GTFRI" {
		packet := map[string]interface{}{
			"protocolVer":               params[1],
			"imei":                      params[2],
			"dvceName":                  params[3],
			"vin":                       params[4],
			"qr":                        params[5],
			"rs1":                       params[6],
			"rs2":                       params[7],
			"reportType":                params[8],
			"ecuErrCode":                params[9],
			"gpsAccuracy":               params[10],
			"speed":                     params[11],
			"azimut":                    params[12],
			"alt":                       params[13],
			"lon":                       params[14],
			"lat":                       params[15],
			"gpsUtcTime":                params[16],
			"rawGNSSInfo":               params[17],
			"mcc":                       params[18],
			"mnc":                       params[19],
			"lac":                       params[20],
			"cellID":                    params[21],
			"csq":                       params[22],
			"networkType":               params[23],
			"status":                    params[24],
			"powerSupply":               params[25],
			"mainSupplyVolt":            params[26],
			"backupBatteryVolt":         params[27],
			"percentageOfBackupBattery": params[28],
			"ecuErrType":                params[29],
			"rs3":                       params[30],
			"ecuLockStatus":             params[31],
			"rs4":                       params[32],
			"ecuInfo":                   params[33],
			"eScootBatteryPercentage":   params[34],
			"genTime":                   params[35],
			"totalCount":                params[36],
		}

		data, _ := json.Marshal(packet)
		return string(data), nil

	}
	return "", nil
}

func HeadInfo(head string) (string, string) {
	parts := strings.Split(head, ":")
	return parts[0], parts[1]
}

func parseEcu() {
	//
}
