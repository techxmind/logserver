package consumer

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	pb "github.com/techxmind/logserver/interface-defs"
)

type standardName struct {
	original string
	standard string
	key      string
}

var (
	errFieldNotExists = errors.New("EventLog field not exists")
)

func newStandardName(original string) (*standardName, error) {
	var (
		e pb.EventLog
		s = &standardName{
			original: original,
		}
	)

	if !strings.Contains(original, ".") {
		s.standard = strings.ToLower(strings.Replace(original, "_", "", -1))
	} else {
		segments := strings.Split(original, ".")
		if len(segments) != 2 {
			return nil, errors.Wrap(errFieldNotExists, original)
		}
		s.standard = strings.ToLower(strings.Replace(segments[0], "_", "", -1))
		s.key = segments[1]
	}

	if _, err := s.GetValue(&e); err == errFieldNotExists {
		return nil, errors.Wrap(err, original)
	}

	return s, nil
}

func (s *standardName) Original() string {
	return s.original
}

func (s *standardName) MustGetValue(e *pb.EventLog) string {
	if v, err := s.GetValue(e); err != nil {
		return ""
	} else {
		return v
	}
}

func (s *standardName) GetValue(e *pb.EventLog) (string, error) {

	switch s.standard {
	// TPL.EVENT_LOG_GETTER.START
	case "actiontype":
		if s.key == "" {
			return e.ActionType, nil
		}
		return "", errFieldNotExists
	case "androidid":
		if s.key == "" {
			return e.AndroidId, nil
		}
		return "", errFieldNotExists
	case "appchannel":
		if s.key == "" {
			return e.AppChannel, nil
		}
		return "", errFieldNotExists
	case "apptype":
		if s.key == "" {
			return e.AppType, nil
		}
		return "", errFieldNotExists
	case "appversion":
		if s.key == "" {
			return e.AppVersion, nil
		}
		return "", errFieldNotExists
	case "carrier":
		if s.key == "" {
			return strconv.FormatInt(int64(e.Carrier), 10), nil
		}
		return "", errFieldNotExists
	case "devicebrand":
		if s.key == "" {
			return e.DeviceBrand, nil
		}
		return "", errFieldNotExists
	case "devicemodel":
		if s.key == "" {
			return e.DeviceModel, nil
		}
		return "", errFieldNotExists
	case "devicevendor":
		if s.key == "" {
			return e.DeviceVendor, nil
		}
		return "", errFieldNotExists
	case "duration":
		if s.key == "" {
			return strconv.FormatInt(int64(e.Duration), 10), nil
		}
		return "", errFieldNotExists
	case "env":
		if s.key == "" {
			return e.Env, nil
		}
		return "", errFieldNotExists
	case "event":
		if s.key == "" {
			return e.Event, nil
		}
		return "", errFieldNotExists
	case "eventid":
		if s.key == "" {
			return e.EventId, nil
		}
		return "", errFieldNotExists
	case "eventtime":
		if s.key == "" {
			return strconv.FormatInt(int64(e.EventTime), 10), nil
		}
		return "", errFieldNotExists
	case "extendinfo":
		if s.key != "" {
			if e.ExtendInfo == nil {
				return "", nil
			}
			return e.ExtendInfo[s.key], nil
		}
		if e.ExtendInfo == nil {
				return "{}", nil
			}
		if bs, err := json.Marshal(e.ExtendInfo); err == nil {
			return string(bs), nil
		}
		return "{}", nil
	case "idfa":
		if s.key == "" {
			return e.Idfa, nil
		}
		return "", errFieldNotExists
	case "imei":
		if s.key == "" {
			return e.Imei, nil
		}
		return "", errFieldNotExists
	case "ip":
		if s.key == "" {
			return e.Ip, nil
		}
		return "", errFieldNotExists
	case "ipcity":
		if s.key == "" {
			return e.IpCity, nil
		}
		return "", errFieldNotExists
	case "ipcountry":
		if s.key == "" {
			return e.IpCountry, nil
		}
		return "", errFieldNotExists
	case "ipprovince":
		if s.key == "" {
			return e.IpProvince, nil
		}
		return "", errFieldNotExists
	case "lat":
		if s.key == "" {
			return e.Lat, nil
		}
		return "", errFieldNotExists
	case "layoutid":
		if s.key == "" {
			return e.LayoutId, nil
		}
		return "", errFieldNotExists
	case "loggedtime":
		if s.key == "" {
			return strconv.FormatInt(int64(e.LoggedTime), 10), nil
		}
		return "", errFieldNotExists
	case "lon":
		if s.key == "" {
			return e.Lon, nil
		}
		return "", errFieldNotExists
	case "mac":
		if s.key == "" {
			return e.Mac, nil
		}
		return "", errFieldNotExists
	case "mid":
		if s.key == "" {
			return e.Mid, nil
		}
		return "", errFieldNotExists
	case "moduleid":
		if s.key == "" {
			return e.ModuleId, nil
		}
		return "", errFieldNotExists
	case "network":
		if s.key == "" {
			return strconv.FormatInt(int64(e.Network), 10), nil
		}
		return "", errFieldNotExists
	case "oaid":
		if s.key == "" {
			return e.Oaid, nil
		}
		return "", errFieldNotExists
	case "os":
		if s.key == "" {
			return e.Os, nil
		}
		return "", errFieldNotExists
	case "osversion":
		if s.key == "" {
			return e.OsVersion, nil
		}
		return "", errFieldNotExists
	case "pageid":
		if s.key == "" {
			return e.PageId, nil
		}
		return "", errFieldNotExists
	case "pagekey":
		if s.key == "" {
			return e.PageKey, nil
		}
		return "", errFieldNotExists
	case "platform":
		if s.key == "" {
			return e.Platform, nil
		}
		return "", errFieldNotExists
	case "pvid":
		if s.key == "" {
			return e.PvId, nil
		}
		return "", errFieldNotExists
	case "reflayoutid":
		if s.key == "" {
			return e.RefLayoutId, nil
		}
		return "", errFieldNotExists
	case "refmoduleid":
		if s.key == "" {
			return e.RefModuleId, nil
		}
		return "", errFieldNotExists
	case "refpageid":
		if s.key == "" {
			return e.RefPageId, nil
		}
		return "", errFieldNotExists
	case "refpagekey":
		if s.key == "" {
			return e.RefPageKey, nil
		}
		return "", errFieldNotExists
	case "refpvid":
		if s.key == "" {
			return e.RefPvId, nil
		}
		return "", errFieldNotExists
	case "referer":
		if s.key == "" {
			return e.Referer, nil
		}
		return "", errFieldNotExists
	case "screenheight":
		if s.key == "" {
			return strconv.FormatInt(int64(e.ScreenHeight), 10), nil
		}
		return "", errFieldNotExists
	case "screenresolution":
		if s.key == "" {
			return e.ScreenResolution, nil
		}
		return "", errFieldNotExists
	case "screensize":
		if s.key == "" {
			return e.ScreenSize, nil
		}
		return "", errFieldNotExists
	case "screenwidth":
		if s.key == "" {
			return strconv.FormatInt(int64(e.ScreenWidth), 10), nil
		}
		return "", errFieldNotExists
	case "sessionid":
		if s.key == "" {
			return e.SessionId, nil
		}
		return "", errFieldNotExists
	case "tkid":
		if s.key == "" {
			return e.Tkid, nil
		}
		return "", errFieldNotExists
	case "udid":
		if s.key == "" {
			return e.Udid, nil
		}
		return "", errFieldNotExists
	case "useragent":
		if s.key == "" {
			return e.UserAgent, nil
		}
		return "", errFieldNotExists
	// TPL.EVENT_LOG_GETTER.END
	case "eventtimestr":
		if s.key == "" {
			return time.Unix(int64(e.EventTime/1000), 0).Format("2006-01-02T15:04:05"), nil
		}
		return "", errFieldNotExists
	case "loggedtimestr":
		if s.key == "" {
			return time.Unix(int64(e.LoggedTime/1000), 0).Format("2006-01-02T15:04:05"), nil
		}
		return "", errFieldNotExists
	case "loggeddate":
		if s.key == "" {
			return time.Unix(int64(e.LoggedTime/1000), 0).Format("20060102"), nil
		}
		return "", errFieldNotExists
	default:
		return "", errFieldNotExists
	}
}
