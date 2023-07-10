package server

type (
	SvrConfig struct {
		IP_Address string
		IP_Port    int16
	}
	SvrData struct {
		Status string
	}
)

var (
	svrChannel chan SvrData
	svrData    = SvrData{
		Status: "data from server goes here",
	}
	ipAddress string
	ipPort    int16
)

func Initialize(cfg SvrConfig, ch chan SvrData) {
	ipAddress = cfg.IP_Address
	ipPort = cfg.IP_Port
	svrChannel = ch

	svrChannel <- svrData
}
