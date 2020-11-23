package eventlog

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/techxmind/ip2location"

	pb "github.com/techxmind/logserver/interface-defs"
)

// Fill system and common property into EventLog
//
func Fill(ctx context.Context, common *pb.EventLogCommon, logs []*pb.EventLog) {
	var (
		loggedTime = time.Now().UnixNano() / int64(1e6)
		ua         = getStringValueFromCtx(ctx, "user-agent")
		ip         = getStringValueFromCtx(ctx, "remote-ip")
		referer    = getStringValueFromCtx(ctx, "referer")
		ipCountry  string
		ipProvince string
		ipCity     string
	)

	if ip != "" {
		loc, err := ip2location.Get(ip)
		if err == nil {
			ipCountry = loc.Country
			ipProvince = loc.Province
			ipCity = loc.City
		}
	}

	// strip base64 tail padding(=)
	if common != nil && common.Udid != "" && strings.HasSuffix(common.Udid, "=") {
		common.Udid = strings.TrimRight(common.Udid, "=")
	}

	for _, log := range logs {
		log.LoggedTime = loggedTime
		log.UserAgent = ua
		log.Ip = ip
		log.IpCountry = ipCountry
		log.IpProvince = ipProvince
		log.IpCity = ipCity
		log.Referer = referer

		// strip base64 tail padding(=)
		if log.Udid != "" && strings.HasSuffix(log.Udid, "=") {
			log.Udid = strings.TrimRight(log.Udid, "=")
		}

		// TPL.EVENT_LOG_FILL.START
		if log.AndroidId == "" && common != nil && common.AndroidId != "" {
			log.AndroidId = common.AndroidId
		}

		if log.AppChannel == "" && common != nil && common.AppChannel != "" {
			log.AppChannel = common.AppChannel
		}

		if log.AppType == "" && common != nil && common.AppType != "" {
			log.AppType = common.AppType
		}

		if log.AppVersion == "" && common != nil && common.AppVersion != "" {
			log.AppVersion = common.AppVersion
		}

		if log.Carrier == 0 && common != nil && common.Carrier != 0 {
			log.Carrier = common.Carrier
		}

		if log.DeviceBrand == "" && common != nil && common.DeviceBrand != "" {
			log.DeviceBrand = common.DeviceBrand
		}

		if log.DeviceModel == "" && common != nil && common.DeviceModel != "" {
			log.DeviceModel = common.DeviceModel
		}

		if log.DeviceVendor == "" && common != nil && common.DeviceVendor != "" {
			log.DeviceVendor = common.DeviceVendor
		}

		if log.Env == "" && common != nil && common.Env != "" {
			log.Env = common.Env
		}

		if log.Idfa == "" && common != nil && common.Idfa != "" {
			log.Idfa = common.Idfa
		}

		if log.Imei == "" && common != nil && common.Imei != "" {
			log.Imei = common.Imei
		}

		if log.Lat == "" && common != nil && common.Lat != "" {
			log.Lat = common.Lat
		}

		if log.Lon == "" && common != nil && common.Lon != "" {
			log.Lon = common.Lon
		}

		if log.Mac == "" && common != nil && common.Mac != "" {
			log.Mac = common.Mac
		}

		if log.Mid == "" && common != nil && common.Mid != "" {
			log.Mid = common.Mid
		}

		if log.Network == 0 && common != nil && common.Network != 0 {
			log.Network = common.Network
		}

		if log.Oaid == "" && common != nil && common.Oaid != "" {
			log.Oaid = common.Oaid
		}

		if log.Os == "" && common != nil && common.Os != "" {
			log.Os = common.Os
		}

		if log.OsVersion == "" && common != nil && common.OsVersion != "" {
			log.OsVersion = common.OsVersion
		}

		if log.Platform == "" && common != nil && common.Platform != "" {
			log.Platform = common.Platform
		}

		if log.ScreenHeight == 0 && common != nil && common.ScreenHeight != 0 {
			log.ScreenHeight = common.ScreenHeight
		}

		if log.ScreenResolution == "" && common != nil && common.ScreenResolution != "" {
			log.ScreenResolution = common.ScreenResolution
		}

		if log.ScreenSize == "" && common != nil && common.ScreenSize != "" {
			log.ScreenSize = common.ScreenSize
		}

		if log.ScreenWidth == 0 && common != nil && common.ScreenWidth != 0 {
			log.ScreenWidth = common.ScreenWidth
		}

		if log.Tkid == "" && common != nil && common.Tkid != "" {
			log.Tkid = common.Tkid
		}

		if log.Udid == "" && common != nil && common.Udid != "" {
			log.Udid = common.Udid
		}

		// TPL.EVENT_LOG_FILL.END

		if log.ScreenResolution == "" && log.ScreenWidth > 0 {
			log.ScreenResolution = fmt.Sprintf("%dx%d", log.ScreenWidth, log.ScreenHeight)
		}
	}
}

// FillSingle wrapper Fill func for single log case
//
func FillSingle(ctx context.Context, log *pb.EventLog) {
	Fill(ctx, nil, []*pb.EventLog{log})
}

func getStringValueFromCtx(ctx context.Context, name string) string {
	if v, ok := ctx.Value(name).(string); ok {
		return v
	}
	return ""
}
