/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package encoderClient

import (
	"bufio"
	"fmt"
	"net"
	"q100transmitter/logger"
	"strings"
	"time"
)

type (
	// API
	HeConfig struct {
		Codecs       string
		AudioBitRate string
		VideoBitRate string
		Spare1       string
		Spare2       string

		// settings not used by the GUI
		audio_codec           string // "ACC"
		audio_bitrate         string // "64000" or "32000"
		audio_sample_rate     string // "44100" or "48000"
		audio_bits_per_sample string // "16"
		// video_codec           string // "H.265"
		// video_size            string // "1280x720"
		// video_bitrate         string // "330"
		// spare_1               string
		// spare_2               string
		StreamUrl string // "udp://192.168.3.10:8282"
		ConfigIP  string // 192.168.3.1"

		// exclusive to HEV-10 commands
		chn           string // The main stream is 0, and the sub stream is 1.
		bps           string // video encoding bit rate, in bps. Range is [32-16384]
		fps           string // video encoding frame rate.
		res_w         string // The horizontal resolution of the encoded video.
		res_h         string // The vertical resolution of the encoded video.
		_type         string // video encoding format, H.264 is 0, H.265 is 1
		gop           string // The range is [1-600]
		profile       string // baseline is 0, main is 1, high is 2
		rc_mode       string // bit rate control mode, CBR is 0, VBR is 1, AVBR is 2, FixQP is 3
		qp1, qp2, qp3 string // bit rate control parameters, related to rc_mode value.
		// When rc_mode = 1 string //
		// qp1 string // max_qp, the value range is [1, 51]
		// qp2 string // min_qp, the value range is [0, 50]
		// qp3 string //min_iqp
	}
)

var (
	arg = HeConfig{}
)

// API
func Initialize(cfg *HeConfig) {
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

	arg.StreamUrl = cfg.StreamUrl
	arg.ConfigIP = cfg.ConfigIP

}

// API
//
// setarams is called from tuner. The function will send the params to the HEV-10 encoder.
func SetParams(cfg *HeConfig) {
	// audio settings used by the GUI
	arg.audio_codec = strings.Fields(cfg.Codecs)[1]
	arg.audio_bitrate = cfg.AudioBitRate

	var hevAudioCmdStr string
	var hevVideoCmdStr string

	switch arg.audio_codec {
	case "ACC":
		// Command format: @0001,23,06,00,01,44100,16,bps!
		hevAudioCmdStr = fmt.Sprintf("@0001,23,06,00,01,%v,%v,%v!",
			arg.audio_sample_rate,
			arg.audio_bits_per_sample,
			arg.audio_bitrate)
	case "G711u":
		// Command format: @0001,23,06,00,00,8000,16,0!
		hevAudioCmdStr = fmt.Sprintf("@0001,23,06,00,00,%v,%v,0!",
			arg.audio_sample_rate,
			arg.audio_bits_per_sample)
	}

	// video settings used by the GUI
	// chn: The main stream is 0, and the sub stream is 1.

	// video encoding bit rate, in bps. Range is [32-16384]
	arg.bps = cfg.VideoBitRate

	// fps: video encoding frame rate.

	// res_w: The horizontal resolution of the encoded video.

	// res_h: The vertical resolution of the encoded video.

	// type: video encoding format, H.264 is 0, H.265 is 1
	switch strings.Fields(cfg.Codecs)[0] {
	case "H264":
		arg._type = "0"
	case "H265":
		arg._type = "1"
	}
	// Command format: @0001,22,06,chn,bps,fps,res_w,res_h,type,gop,pro<ile,rc_mode,qp1,qp2,qp3!
	hevVideoCmdStr = fmt.Sprintf("@0001,22,06,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v!",
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

	sendToEncoder(hevAudioCmdStr, hevVideoCmdStr)

	// gop: The range is [1-600]

	// profile: baseline is 0, main is 1, high is 2

	// rc_mode: bit rate control mode, CBR is 0, VBR is 1, AVBR is 2, FixQP is 3

	// qp1, qp2, qp3: bit rate control parameters, related to rc_mode value.

	// When rc_mode = 3:
	// qp1: iqp, the value range is [0, 50]
	// qp2: pqp, the value range is [0, 50]
	// qp3: bqp
}

/*
My Working Encoder Settings

iMac has 2nd IP 192.168.3.3

HEV-10 IP address is 192.168.3.1 and sending UDP stream to 192.168.3.10

Secondary Stream
Audio Encoding Type: AAC
Audio Bitrate(bps): 64000
Video Encoding Type: H.265
Video Encoding Size: 1280*720
Video Bitrate(Kbps): 330
RTSP URL1: Disabled
RTSP URL2(TS): Disabled
Misc Stream/Secondary Stream UDP(unicast/multicast) URL: udp://192.168.3.10:8282
SRT: Disabled

*/

func sendToEncoder(audioCmd, videoCmd string) {
	// Windows Serial Port Utility...
	// Port: TCP/UDP
	// Mode: TCP Client
	// Host: 192.168.1.251
	// Port: 55555
	// using either arg.Url or arg.IP_Address
	// Command execution success: #8001,23,06,OK!
	// Command execution fails: #8001,23,06,ERR!
	const (
		PORT    = "55555"
		SUCCESS = "#8001,23,06,OK!"
		FAIL    = "#8001,23,06,ERR!"
	)

	url := fmt.Sprintf("%s:%s", arg.ConfigIP, PORT)
	logger.Info.Printf("Connecting to: %s", url)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		logger.Error.Printf("Failed to connect to: %s", url)
		return
	}
	logger.Info.Printf("Connected to: %v", url)
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(50 * time.Millisecond)); err != nil {
		logger.Error.Printf("Failed to set timeout: %s", err)
		return
	}

	// AUDIO

	logger.Info.Printf("will send audio cmd %s to HEV-10", audioCmd)

	// send
	_, err = conn.Write([]byte(audioCmd))
	if err != nil {
		println("Write failed:", err.Error())
		logger.Error.Printf("Failed to write: %s", err)
		return
	}
	// receive
	audioBuffer := bufio.NewReader(conn)
	audioResult, err := audioBuffer.ReadString('!')
	if err != nil {
		logger.Error.Printf("Failed to read result")
		return
	}
	switch audioResult {
	case FAIL:
		logger.Error.Printf("Failed to send audioCmd")
		return
	case SUCCESS:
		logger.Info.Printf("HEV-10 audio configured ok")
	default:
		logger.Error.Printf("Undefine result: %v", audioResult)
		return
	}

	// VIDEO

	logger.Info.Printf("will send video cmd %s to HEV-10", videoCmd)

	// send
	_, err = conn.Write([]byte(videoCmd))
	if err != nil {
		println("Write failed:", err.Error())
		logger.Error.Printf("Failed to write: %s", err)
		return
	}
	// receive
	videoBuffer := bufio.NewReader(conn)
	videoResult, err := videoBuffer.ReadString('!')
	if err != nil {
		logger.Error.Printf("Failed to read result")
		return
	}
	switch videoResult {
	case FAIL:
		logger.Error.Printf("Failed to send videoCmd")
		return
	case SUCCESS:
		logger.Info.Printf("HEV-10 video configured ok")
	default:
		logger.Error.Printf("Undefine result: >%v<", videoResult)
		return
	}
}
