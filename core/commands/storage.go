package commands

import (
	"fmt"
	cmds "github.com/ipfs/go-ipfs-cmds"
	"strings"
)

const (
	uploadPriceOptionName = "price"
	replicationFactorOptionName = "replication-factor"
	hostSelectModeOptionName = "host-select-mode"
	hostSelectionOptionName = "host-selection"
)
var StorageCmd = &cmds.Command{
	Subcommands: map[string]*cmds.Command{
		"upload": storageUploadCmd,
	},
}

var storageUploadCmd = &cmds.Command{
	Arguments: []cmds.Argument{
		cmds.StringArg("file-hash", true, false, "add hash of file to upload").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.Int64Option(uploadPriceOptionName, "p", "max price per GB of storage in BTT"),
		cmds.Int64Option(replicationFactorOptionName, "replica", "replication factor for the file with erasure coding built-in").WithDefault(int64(3)),
		cmds.StringOption(hostSelectModeOptionName, "mode", "based on mode to select the host and upload automatically").WithDefault("score"),
		cmds.StringOption(hostSelectionOptionName, "list", "use only these hosts in order on CUSTOM mode"),
	},
	PreRun: func(req *cmds.Request, env cmds.Environment) error {
		price, found := req.Options[uploadPriceOptionName].(int64)
		if found {
			if price < 0 {
				return fmt.Errorf("cannot input a negative price")
			}
		} else {
			// TODO: default value to be set as top-selected nodes' price
			// currently set price with 10000
			req.Options[uploadPriceOptionName] = int64(10)
		}

		mode, _ := req.Options[hostSelectModeOptionName].(string)
		list, found := req.Options[hostSelectionOptionName].(string)
		if strings.EqualFold(mode, "custom") {
			if !found {
				return fmt.Errorf("custom mode needs input host lists")
			} else {
				req.Options[hostSelectionOptionName] = strings.Split(list, ",")
			}
		} else {
			req.Options[hostSelectionOptionName] = nil
		}
		return nil
	},
	Run: func(req *cmds.Request, emitter cmds.ResponseEmitter, environment cmds.Environment) error {
		for _, hash := range req.Arguments {
			err := emitter.Emit(hash)
			if err != nil {
				return err
			}
		}
		price, _ := req.Options[uploadPriceOptionName].(int64)
		replica, _ := req.Options[replicationFactorOptionName].(int64)
		mode, _ := req.Options[hostSelectModeOptionName].(string)
		list, _ := req.Options[hostSelectionOptionName].([]string)
		return emitter.Emit(fmt.Sprintf("price %d, replica %d, mode %s, list %s", price, replica, mode, list))
	},
}
