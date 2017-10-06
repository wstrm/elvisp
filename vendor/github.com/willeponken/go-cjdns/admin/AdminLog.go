package admin

import "errors"

const (
	// CJDNS log levels:
	KEYS     = "KEYS"     // Not compiled in by default, contains private keys and other secret information.
	DEBUG    = "DEBUG"    // Default level, contains lots of information which is probably not useful unless you are diagnosing an ongoing problem.
	INFO     = "INFO"     // Shows starting and stopping of various components and general purpose information.
	WARN     = "WARN"     // Generally this means some system has undergone a minor failure, this includes failures due to network disturbance.
	ERROR    = "ERROR"    // This means there was a (possibly temporary) failure of a system within cjdns.
	CRITICAL = "CRITICAL" // This means something is broken such that the cjdns core will likely have to exit immedietly.
)

// LogMessage represents a log entry returned from CJDNS.
type LogMessage struct {
	File    string `bencode:"file"`
	Level   string `bencode:"level"`
	Line    int    `bencode:"line"`
	Message string `bencode:"message"`
	Time    int64  `bencode:"time"`
}

func (m *LogMessage) String() string { return m.Message }

// Subscribes you to receive logging updates based on the parameters that are set.
// Set file to "" to log from all files, set line to -1 lo log from any line.
func (a *Conn) AdminLog_subscribe(level, file string, line int, c chan<- *LogMessage) (streamId string, err error) {
	var pack *packet
	req := request{AQ: "AdminLog_subscribe"}
	if file != "" {
		if line != -1 {
			args := new(struct {
				Line  int    `bencode:"line"`
				File  string `bencode:"file"`
				Level string `bencode:"level"`
			})
			args.Line = line
			args.File = file
			args.Level = level
			req.Args = args

		} else {
			args := new(struct {
				File  string `bencode:"file"`
				Level string `bencode:"level"`
			})
			args.File = file
			args.Level = level
			req.Args = args
		}
	} else {
		args := new(struct {
			Level string `bencode:"level"`
		})
		args.Level = level
		req.Args = args
	}

	if pack, err = a.sendCmd(&req); err != nil {
		return
	}
	res := new(struct {
		StreamId string `bencode:"streamId"`
		Error    string `bencode:"error"`
	})
	if err = pack.Decode(res); err != nil {
		return
	}
	streamId = res.StreamId
	if res.Error != "none" {
		err = errors.New(res.Error)
		return
	}
	a.registerLogChan(res.StreamId, c)
	return
}

// Removes the logging subscription so that you no longer recieve log info.
func (a *Conn) AdminLog_unsubscribe(streamId string) error {
	args := new(struct {
		StreamId string `bencode:"streamId"`
	})
	args.StreamId = streamId

	pack, err := a.sendCmd(&request{AQ: "AdminLog_unsubscribe", Args: args})
	if err != nil {
		return err
	}
	res := new(struct{ Error string })
	pack.Decode(res)
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}
