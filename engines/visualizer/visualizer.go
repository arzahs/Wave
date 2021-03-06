package visualizer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hkparker/Wave/helpers"
	"github.com/hkparker/Wave/models"
	"strings"
	"sync"
)

const (
	DEVICE_ISAP           = "IsAP"
	DEVICE_MAC            = "MAC"
	DEVICE_NULLPROBE      = "NullProbe"
	DEVICE_PROBE          = "ProbedFor"
	DEVICE_POWERSTATE     = "PowerState"
	DEVICE_POWERSTATE_ON  = "online"
	DEVICE_POWERSTATE_OFF = "offline"

	TYPE_NULL_PROBE_REQUEST = "NullProbeRequest"
	TYPE_PROBE_REQUEST      = "ProbeRequest"
	TYPE_UPDATE_AP          = "UpdateAccessPoint"

	EVENT = "Event"
	SSID  = "SSID"
	TYPE  = "type"
)

type VisualEvent map[string]string

var VisualEvents = make(chan VisualEvent, 0)
var VendorBytes = make(map[string]string)

var Devices = make(map[string]models.Device)
var DevicesMux sync.Mutex
var Networks = make(map[string]models.Network)
var NetworksMux sync.Mutex
var Associations = make(map[string][]string)
var AssociationsMux sync.Mutex

func init() {
	loadMetadata()
}

func loadMetadata() {
	prefix_path := "engines/visualizer/metadata/nmap-mac-prefixes"
	vendor_data, err := helpers.Asset(prefix_path)
	if err != nil {
		log.WithFields(log.Fields{
			"at":    "visualizer.loadMetadata",
			"error": err.Error(),
		}).Error("unable to load vendor bytes")
		return
	}
	lines := strings.Split(string(vendor_data), "\n")
	for _, line := range lines {
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}
		raw_mac := line[0:6]
		name := line[7:]
		mac := raw_mac[0:2] + ":" +
			raw_mac[2:4] + ":" +
			raw_mac[4:6]
		VendorBytes[strings.ToUpper(mac)] = name
	}
}

func Insert(frame models.Wireless80211Frame) {
	DevicesMux.Lock()
	defer DevicesMux.Unlock()
	NetworksMux.Lock()
	defer NetworksMux.Unlock()
	AssociationsMux.Lock()
	defer AssociationsMux.Unlock()
	//updateKnownDevices(frame)

	if len(frame.Type) < 4 {
		log.WithFields(log.Fields{}).Warn("frame type too small")
		return
	}
	switch frame.Type[:4] {
	case "Mgmt":
		insertMgmt(frame)
	case "Data":
		insertData(frame)
	case "Ctrl":
		insertCtrl(frame)
	//case "Rese":
	default:
		log.WithFields(log.Fields{
			"at":         "visualizer.Insert",
			"frame.Type": frame.Type,
		}).Warn("unknown frame type")
	}
}

func insertMgmt(frame models.Wireless80211Frame) {
	switch frame.Type {
	case "MgmtAssociationReq":
	case "MgmtAssociationResp":
	case "MgmtReassociationReq":
	case "MgmtReassociationResp":
	case "MgmtProbeReq":
		updateKnownDevices(frame)
		updateProbeRequests(frame)
	case "MgmtProbeResp":
	case "MgmtMeasurementPilot":
	case "MgmtBeacon":
		updateKnownDevices(frame)
		updateAccessPoints(frame)
	case "MgmtATIM":
	case "MgmtDisassociation":
		updateKnownDevices(frame)
		animateDeauth(frame)
	case "MgmtAuthentication":
	case "MgmtDeauthentication":
		updateKnownDevices(frame)
		animateDeauth(frame)
	case "MgmtAction":
	case "MgmtActionNoAck":
	default:
		log.WithFields(log.Fields{
			"at":   "visualizer.insertMgmt",
			"type": frame.Type,
		}).Warn("unknown frame type")
	}
}

func insertData(frame models.Wireless80211Frame) {
	//updateTx()
	switch frame.Type {
	case "Data":
		updateKnownDevices(frame)
		updateAssociation(frame)
	case "DataCFAck":
		updateKnownDevices(frame)
		updateAssociation(frame)
	case "DataCFPoll":
		updateKnownDevices(frame)
		updateAssociation(frame)
	case "DataCFAckPoll":
		updateKnownDevices(frame)
		updateAssociation(frame)
	case "DataNull":
		updateKnownDevices(frame)
		updateDataNull(frame)
	case "DataCFAckNoData":
	case "DataCFPollNoData":
	case "DataCFAckPollNoData":
	case "DataQOSData":
		updateKnownDevices(frame)
		updateAssociation(frame)
	case "DataQOSDataCFAck":
		updateKnownDevices(frame)
		updateAssociation(frame)
	case "DataQOSDataCFPoll":
		updateKnownDevices(frame)
		updateAssociation(frame)
	case "DataQOSDataCFAckPoll":
		updateKnownDevices(frame)
		updateAssociation(frame)
	case "DataQOSNull":
	case "DataQOSCFPollNoData":
	case "DataQOSCFAckPollNoData":
	default:
		log.WithFields(log.Fields{
			"at":   "visualizer.insertData",
			"type": frame.Type,
		}).Warn("unknown frame type")
	}
}

func insertCtrl(frame models.Wireless80211Frame) {
	switch frame.Type {
	case "CtrlWrapper":
	case "CtrlBlockAckReq":
	case "CtrlBlockAck":
	case "CtrlPowersavePoll":
	case "CtrlRTS":
	case "CtrlCTS":
	case "CtrlAck":
	case "CtrlCFEnd":
	case "CtrlCFEndAck":
	default:
		log.WithFields(log.Fields{
			"at":   "visualizer.insertCtrl",
			"type": frame.Type,
		}).Warn("unknown frame type")
	}
}

func CatchupEvents() []VisualEvent {
	catchup_events := make([]VisualEvent, 0)
	for _, device := range Devices {
		catchup_events = append(catchup_events, device.VisualData())
	}
	for _, network := range Networks {
		ssid_set := network.VisualData()
		for _, network_event := range ssid_set {
			catchup_events = append(catchup_events, VisualEvent(network_event))
		}
	}
	for mac, associations := range Associations {
		for _, other := range associations {
			catchup_events = append(
				catchup_events,
				VisualEvent{
					TYPE:   "NewAssociation",
					"MAC1": mac,
					"MAC2": other,
				},
			)
		}
	}
	catchup_events = append(
		catchup_events,
		VisualEvent{
			TYPE: "CacheCleared",
		},
	)
	// add other resources, create other events
	return catchup_events
}
