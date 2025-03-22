package okaiparser

import (
	"fmt"
	okaiparsetools "okai/common/okai-parse-tools"
	"okai/common/utils"
	"strings"
	"time"
)

func ParseParams(params []string) (string, string, map[string]interface{}, error) {
	if len(params) == 0 {
		return "", "", nil, fmt.Errorf("Parse params failed. Zero len")
	}

	pktType, pktId := HeadInfo(params[0])

	if (pktType == "+RESP" || pktType == "+BUFF") && pktId == "GTFRI" {
		packet := parseBasePacket(params)
		return pktType, pktId, packet, nil
	} else if pktType == "+RESP" && pktId == "GTNCN" {
		packet := parseBasePacket(params)
		return pktType, pktId, packet, nil
	} else if pktType == "+ACK" && pktId == "GTRTO" {
		packet := map[string]interface{}{
			"cmdID": params[6],
		}
		return pktType, pktId, packet, nil
	} else if pktType == "+ACK" && pktId == "GTECC" {
		packet := map[string]interface{}{
			"cmdID": params[5],
		}
		return pktType, pktId, packet, nil
	}

	// fmt.Println("----------------------------")
	// fmt.Printf("TYPE: %s | ID: %s\nparams: %s\n", pktType, pktId, params)

	return pktType, pktId, nil, nil
}

func HeadInfo(head string) (string, string) {
	parts := strings.Split(head, ":")
	return parts[0], parts[1]
}

func parseBasePacket(params []string) map[string]interface{} {
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
		"online":                  true,
		"_ts":                     time.Now().Unix(),
		"charge":                  params[34],
	}

	// rawGNSSInfo := params[17]
	rawEcuInfo := params[33]
	ecuInfo := parseEcu(rawEcuInfo)
	packet["ecuInfo"] = ecuInfo
	packet["gnssInfo"] = params[17]
	fmt.Println("RAW GNSS_INFO", params[17])

	return packet
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

func CommandBuilder(cmd map[string]string, tc string) string {

	switch cmd["head"] {
	case "GTRTO":
		return fmt.Sprintf("AT+GTRTO=zk200,%s,,%x,,,,,,%s$", cmd["subcommand"], 0, tc)
	case "GTECC":
		return fmt.Sprintf("AT+GTECC=zk200,,%s,,1,,,,,,%s$", cmd["subcommand"], tc)
	default:
		return fmt.Sprintf("AT+GTRTO=zk200,%s,,0,,,,,,%s$", cmd["subcommand"], tc)
	}
}
