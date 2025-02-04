package states

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// getForceReleaseCmd returns command for force-release
// usage: force-release [flags]
func getForceReleaseCmd(cli *clientv3.Client, basePath string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "force-release",
		Short: "Force release the collections from QueryCoord",
		Run: func(cmd *cobra.Command, args []string) {
			// basePath = 'by-dev/meta/'
			// queryCoord prefix = 'queryCoord-'
			prefix := path.Join(basePath, "queryCoord-")
			now := time.Now()
			err := backupEtcd(cli, prefix, fmt.Sprintf("bw_etcd_querycoord.%s.bak.gz", now.Format("060102-150405")))
			if err != nil {
				fmt.Printf("backup etcd failed, error: %v, stop doing force-release\n", err)
			}

			// remove all keys start with [basePath]/queryCoord-
			_, err = cli.Delete(context.Background(), "queryCoord-", clientv3.WithPrefix())
			if err != nil {
				fmt.Printf("failed to remove queryCoord etcd kv, err: %v\n", err)
			}
			// release all collections from online querynodes

			// maybe? kill session of queryCoord?
		},
	}

	return cmd
}
