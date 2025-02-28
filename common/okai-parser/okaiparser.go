package okaiparser

import (
	"encoding/json"
	"fmt"
	okaiparsetools "okai/common/okai-parse-tools"
	"okai/common/utils"
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
			"protocolVer": params[1],
			"imei":        params[2],
			"dvceName":    params[3],
			"vin":         params[4],
			"qr":          params[5],
			"rs1":         params[6],
			"rs2":         params[7],
			"reportType":  params[8],
			"ecuErrCode":  params[9],
			"gpsAccuracy": params[10],
			"speed":       params[11],
			"azimut":      params[12],
			"alt":         params[13],
			"lon":         params[14],
			"lat":         params[15],
			"gpsUtcTime":  params[16],
			// gns here
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
			// ecu info here
			"eScootBatteryPercentage": params[34],
			"genTime":                 params[35],
			"totalCount":              params[36],
		}

		// rawGNSSInfo := params[17]
		rawEcuInfo := params[33]
		ecuInfo := parseEcu(rawEcuInfo)
		packet["ecuInfo"] = ecuInfo

		data, _ := json.Marshal(packet)
		return string(data), nil

	}
	return "", nil
}

func HeadInfo(head string) (string, string) {
	parts := strings.Split(head, ":")
	return parts[0], parts[1]
}

func parseEcu(raw string) map[string]interface{} {
	info := okaiparsetools.SplitParams(raw, "&")
	ecuInfo := map[string]interface{}{
		"speedKmh":                 info[0],
		"currentMileageKm":         info[1],
		"remainingMileageKm":       info[2],
		"totalMileageKm":           info[3],
		"headlightStatus":          info[4],
		"tailLightStatus":          info[5],
		"ridingTime":               info[6],
		"firmwareVer":              info[7],
		"hardwareVer":              info[8],
		"electricity":              info[9],
		"charging":                 info[10],
		"batteryLockStatus":        info[12],
		"batteryCompartmentStatus": info[13],
		"rs1":                      info[14],
	}

	batteryStatus := parseBatteryStatusInfo(info[11])
	ecuInfo["batteryStatus"] = batteryStatus

	return ecuInfo
}

func parseBatteryStatusInfo(batteryRaw string) map[string]interface{} {
	data := utils.HexToBytes(batteryRaw)
	// Группа 1: Статусы MOSFET (заряд/разряд)
	mosStatus := data[0]
	chargingMOS := (mosStatus >> 0) & 1
	dischargingMOS := (mosStatus >> 1) & 1

	// Группа 2–6: Температуры и здоровье батареи
	batteryHealth := data[1]   // Состояние батареи, %
	maxTemp := int8(data[2])   // Макс. температура
	minTemp := int8(data[3])   // Мин. температура
	mosTemp := int8(data[4])   // Температура MOSFET
	otherTemp := int8(data[5]) // Прочие температуры

	btStatus := map[string]interface{}{
		"chargingMOS":    chargingMOS,
		"dischargingMOS": dischargingMOS,
		"batteryHealth":  batteryHealth,
		"maxTemp":        maxTemp,
		"minTemp":        minTemp,
		"mosTemp":        mosTemp,
		"otherTemp":      otherTemp,
	}

	return btStatus
}
