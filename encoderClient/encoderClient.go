/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package encoderClient

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ea7kir/qLog"
)

/*
My working encoder settings as they appear in s Safari

iMac has 2nd IP 192.168.3.3

HEV-10 IP address is 192.168.3.1 and sending UDP stream to 192.168.3.10

Secondary Stream
Audio Encoding Type: AAC
Audio Bitrate(bps): 64000
Video Encoding Type: H.265
Video Encoding Size: 1280*720
Video Bitrate(Kbps): 350
RTSP URL1: Disabled
RTSP URL2(TS): Disabled
Misc Stream/Secondary Stream UDP(unicast/multicast) URL: udp://192.168.3.10:8282
SRT: Disabled
*/

type (
	// API
	HeConfig struct {
		Codecs                string
		AudioBitRate          string
		VideoBitRate          string
		Spare1                string
		Spare2                string
		StreamIP              string // "udp://192.168.3.10:8282"
		StreamPort            string
		ConfigIP              string // 192.168.3.1"
		audio_bitrate         string // "64000" or "32000"
		audio_sample_rate     string // "44100" or "48000"
		audio_bits_per_sample string // "16"
		chn                   string // The main stream is 0, and the sub stream is 1.
		bps                   string // video encoding bit rate, in bps. Range is [32-16384]
		fps                   string // video encoding frame rate.
		res_w                 string // The horizontal resolution of the encoded video.
		res_h                 string // The vertical resolution of the encoded video.
		_type                 string // video encoding format, H.264 is 0, H.265 is 1
		gop                   string // The range is [1-600]
		profile               string // baseline is 0, main is 1, high is 2
		rc_mode               string // bit rate control mode, CBR is 0, VBR is 1, AVBR is 2, FixQP is 3
		qp1, qp2, qp3         string // bit rate control parameters, related to rc_mode value.
		// When rc_mode = 1 string //
		// qp1 string // max_qp, the value range is [1, 51]
		// qp2 string // min_qp, the value range is [0, 50]
		// qp3 string //min_iqp
	}
)

var (
	arg HeConfig
)

// API
func Initialize(cfg HeConfig) {
	arg = cfg
	// settings not used by the GUI
	arg.audio_sample_rate = "44100" // or "48000"
	arg.audio_bits_per_sample = "16"

	arg.chn = "0" // The main stream is 0, and the sub stream is 1.
	// arg.bps = ""       // video encoding bit rate, in bps. Range is [32-16384]
	arg.fps = "25"     // video encoding frame rate.
	arg.res_w = "1280" // The horizontal resolution of the encoded video.
	arg.res_h = "720"  // The vertical resolution of the encoded video.
	// arg._type = ""     // video encoding format, H.264 is 0, H.265 is 1
	arg.gop = "50"    // The range is [1-600]
	arg.profile = "1" // baseline is 0, main is 1, high is 2
	arg.rc_mode = "0" // bit rate control mode, CBR is 0, VBR is 1, AVBR is 2, FixQP is 3
	arg.qp1 = "0"     // bit rate control parameters, related to rc_mode value.
	arg.qp2 = "0"     // bit rate control parameters, related to rc_mode value.
	arg.qp3 = "0"     // bit rate control parameters, related to rc_mode value.
	// When rc_mode = 3:
	// qp1: iqp, the value range is [0, 50]
	// qp2: pqp, the value range is [0, 50]
	// qp3: bqp
}

func SetParams(cfg *HeConfig) error {
	const (
		PORT = "55555"
		// SUCCESS = "#8001,23,06,OK!"
		// FAIL    = "#8001,23,06,ERR!"
	)
	var (
		cmdStr string
		codec  string
	)

	// NETWORK

	url := fmt.Sprintf("%s:%s", arg.ConfigIP, PORT)

	qLog.Info("Connecting to: %s", url)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		qLog.Error("Failed to connect to: %s", url)
		return err
	}
	qLog.Info("Connected to: %v", url)
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
		qLog.Error("Failed to set timeout: %s", err)
		return err
	}

	// AUDIO

	codec = strings.Fields(cfg.Codecs)[1] // extract audio codec from eg "H.265 AAA"
	switch codec {
	case "ACC":
		// update from txControl
		arg.audio_bitrate = cfg.AudioBitRate
		// Command format: @0001,23,06,00,01,44100,16,bps!
		cmdStr = fmt.Sprintf("@0001,23,06,00,01,%v,%v,%v!",
			arg.audio_sample_rate,
			arg.audio_bits_per_sample,
			arg.audio_bitrate)
	case "G711u":
		// Command format: @0001,23,06,00,00,8000,16,0!
		cmdStr = fmt.Sprintf("@0001,23,06,00,00,%v,%v,0!",
			arg.audio_sample_rate,
			arg.audio_bits_per_sample)
	}

	if err := sendToEncoder(conn, cmdStr, "AUDIO"); err != nil {
		return err
	}

	// VIDEO

	// update from txControl
	codec = strings.Fields(cfg.Codecs)[0] // extract video codec from eg "H.265 AAA"
	switch codec {
	case "H264":
		arg._type = "0"
	case "H265":
		arg._type = "1"
	}
	arg.bps = cfg.VideoBitRate
	// Command format: @0001,22,06,chn,bps,fps,res_w,res_h,type,gop,pro<ile,rc_mode,qp1,qp2,qp3!
	cmdStr = fmt.Sprintf("@0001,22,06,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v!",
		arg.chn,
		arg.bps,
		arg.fps,
		arg.res_w,
		arg.res_h,
		arg._type,
		arg.gop,
		arg.profile,
		arg.rc_mode,
		arg.qp1,
		arg.qp2,
		arg.qp3)

	if err := sendToEncoder(conn, cmdStr, "VIDEO"); err != nil {
		return err
	}

	return nil
}

func sendToEncoder(conn net.Conn, cmdStr string, what string) error {
	const SUCCESS_A = "#8001,23,06,OK!"
	const SUCCESS_V = "#8001,22,06,OK!"
	const FAIL = "#8001,23,06,ERR!"
	var err error
	qLog.Info("cmdStr is: %s", cmdStr)
	// send
	_, err = conn.Write([]byte(cmdStr))
	if err != nil {
		return fmt.Errorf("failed to write %s: %s", what, err)
	}
	// receive
	buf := bufio.NewReader(conn)
	result, err := buf.ReadString('!')
	if err != nil {
		return fmt.Errorf("failed to read %s result: %s", what, err)
	}
	switch result {
	case FAIL:
		return fmt.Errorf("failed to send %s to encoder >%v<", what, result)
	case SUCCESS_A, SUCCESS_V:
		qLog.Info("HEV-10 %s configured ok >%v<", what, result)
	default:
		return fmt.Errorf("undefine %s result: >%v<", what, result)
	}
	/*
		#8001,23,06,OK!
		#8001,22,06,OK!
		      #  8  0  0  0  ,  2  3  ,  0  6  ,  O  K  !
			>[35 56 48 48 49 44 50 51 44 48 54 44 79 75 33]<
			  #  8  0  0  0  ,  2  2  ,  0  6  ,  O  K  !
			>[35 56 48 48 49 44 50 50 44 48 54 44 79 75 33]<
	*/
	return nil
}
