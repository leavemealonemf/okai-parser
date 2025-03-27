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
	} else if pktType == "+RESP" && pktId == "GTALC" {
		packet := parseMainCfg(params)
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

func parseMainCfg(params []string) map[string]interface{} {
	packet := map[string]interface{}{
		"_ts":  time.Now().UnixMicro(),
		"imei": params[2],
		"head": map[string]interface{}{
			"protocol_version":   params[1],
			"imei":               params[2],
			"device_name":        params[3],
			"vin":                params[4],
			"qr_code":            params[5],
			"rs1":                params[6],
			"task_id":            params[7],
			"configuration_code": params[8],
		},
		"qss": map[string]interface{}{
			"apn":                    params[10],
			"apn_username":           params[11],
			"apn_password":           params[12],
			"reporting_mode":         params[13],
			"network_mode":           params[14],
			"enable_buffer":          params[15],
			"primary_server_address": params[16],
			"primary_server_port":    params[17],
			"lte_mode":               params[18],
			"region":                 params[19],
			"rs2":                    params[20],
			"heartbeat_interval":     params[21],
			"rs3":                    params[22],
			"enable_bth_unlocking":   params[23],
			"ble_broadcase_name":     params[24],
		},
		"cfg": map[string]interface{}{
			"new_password":                   params[26],
			"device_name":                    params[27],
			"gps_always_on":                  params[28],
			"filter_gps_data_time":           params[29],
			"agps_mode":                      params[30],
			"brake_cfg":                      params[31],
			"reporting_item_mask":            params[32],
			"event_mask":                     params[33],
			"instrument_type_cfg":            params[34],
			"atmosphere_light_vibration_cfg": params[35],
			"auto_lock_when_meter_conn_lost": params[36],
			"poweron_upgrade":                params[37],
			"enable_voice_playback":          params[38],
			"volume_adjustment":              params[39],
			"turn_sig_audio_play":            params[40],
		},
		"tma": map[string]interface{}{
			"marker":               params[42],
			"hourly_offset":        params[43],
			"minute_based_offset":  params[44],
			"daylight_saving_time": params[45],
		},
		"fri": map[string]interface{}{
			"mode": params[52],
			"not_report_location_when_pos_unavailable": params[53],
			"lock_sending_interval":                    params[54],
			"unlock_sending_interval":                  params[55],
			"stanby_pwr_send_interval":                 params[56],
			"use_gps_raw_data":                         params[57],
		},
		"dog": map[string]interface{}{
			"mode":                   params[62],
			"interval":               params[64],
			"date":                   params[65],
			"report_before_restart":  params[67],
			"unit":                   params[68],
			"no_network_interval":    params[69],
			"no_activation_interval": params[70],
			"send_timeout":           params[71],
		},
		"nmd": map[string]interface{}{
			"stale_time":               params[78],
			"motion_duration":          params[79],
			"sens_lvl":                 params[80],
			"sens_of_trigg_motor_lock": params[82],
		},
		"alm": map[string]interface{}{
			"vibration_alarm_duration": params[87],
			"alarm_interval":           params[88],
			"six_axis_sensor_dir":      params[90],
		},
		"ecc": map[string]interface{}{
			"motor_pwr_cfg":              params[93],
			"max_speed_limit":            params[94],
			"acceleration_mode":          params[95],
			"display_unit":               params[96],
			"acceleration_lvl":           params[97],
			"braking_force":              params[98],
			"working_mode_of_tail_light": params[99],
			"lock_force_lvl":             params[100],
			"kinetic_energy_rec_lvl":     params[101],
		},
		"led": map[string]interface{}{
			"center_console_led_settimg":   params[103],
			"headlight_auto_on_off_cfg":    params[104],
			"set_the_turn_sig":             params[105],
			"rgb_setting":                  params[106],
			"status_light_mode":            params[107],
			"status_setting":               params[108],
			"status_indicator_prompt_mode": params[109],
			"animation_lvl":                params[110],
			"mode_cfg":                     params[111],
		},
		"ipn": nil,
		"vad": map[string]interface{}{
			"enable_battery_lock":        params[123],
			"enable_electronic_bell":     params[124],
			"instrument_style_ifce_cfg":  params[125],
			"nfc_work_mode":              params[126],
			"handle_cfg":                 params[127],
			"battery_lock_type":          params[128],
			"battery_lock_alarm_pb_time": params[129],
		},
		"nfc": nil,
		"bcp": map[string]interface{}{
			"pass_type_selection":  params[139],
			"static_pass_str":      params[140],
			"mac_addr_setting":     params[141],
			"mac_addr":             params[142],
			"instument_nfc_switch": params[143],
		},
		"mel": map[string]interface{}{
			"straight_rod_lock_selection":     params[145],
			"helmet_box_lock_selection":       params[146],
			"ulock_selection":                 params[147],
			"num_of_straight_rod_lock_alarms": params[148],
			"num_of_helmet_box_lock_alarms":   params[149],
			"multi_mech_lock_enable_mask":     params[150],
		},
		"dcc": map[string]interface{}{
			"motion_data_coll":          params[153],
			"motion_data_coll_duration": params[154],
			"g_sensor_risk_duration":    params[155],
			"g_sensor_risk_sens":        params[155],
		},
		"hlm": map[string]interface{}{
			"top_light_work_mode":  params[162],
			"side_light_work_mode": params[163],
			"top_light_color_r":    params[164],
			"top_light_color_g":    params[165],
			"top_light_color_b":    params[166],
			"side_light_color_r":   params[167],
			"side_light_color_g":   params[168],
			"side_light_color_b":   params[169],
			"offline_rgb_settings": params[170],
		},
		"nal": map[string]interface{}{
			"mode":                      params[175],
			"server_conn_loss_duration": params[176],
		},
		"rmd": map[string]interface{}{
			"mode":                  params[182],
			"blacklist_operator_1":  params[184],
			"blacklist_operator_20": params[186],
		},
		"bts": map[string]interface{}{
			"discoverable_mode": params[197],
		},
		"cic": map[string]interface{}{
			"ecu_throttle":                   params[215],
			"enable_wireless_charging":       params[216],
			"push_mode":                      params[217],
			"parking_gear":                   params[218],
			"dashboard_charge_auto_off_time": params[219],
			"motor_lock":                     params[220],
			"overspeed_alarm":                params[221],
		},
		"xwm": map[string]interface{}{
			"working_mode":                 params[224],
			"reporting_dir_in_normal_mode": params[225],
			"reporting_dir_in_test_mode":   params[226],
			"customer_code":                params[229],
			"gen_time":                     params[230],
		},
	}
	return packet
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
		"_ts":                     time.Now().UnixMicro(),
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
		return fmt.Sprintf("AT+GTRTO=zk200,%s,,%d,,,,,,%s$", cmd["subcommand"], 0, tc)
	case "GTECC":
		return fmt.Sprintf("AT+GTECC=zk200,,%s,1,1,,,,,,%s$", cmd["subcommand"], tc)
	case "GTVAD":
		return fmt.Sprintf("AT+GTVAD=zk200,%s,,,,,2,,%s$", cmd["subcommand"], tc)
	case "GTQSS":
		return fmt.Sprintf("AT+GTQSS=zk200,,,,,0,,iot-socket.okai.co,14010,,,,,,,,%s$", tc)
	default:
		return fmt.Sprintf("AT+GTRTO=zk200,%s,,0,,,,,,%s$", cmd["subcommand"], tc)
	}
}
