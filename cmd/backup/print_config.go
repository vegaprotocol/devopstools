package backup

import (
	"fmt"

	"github.com/spf13/cobra"
)

const exampleConfig = `# Address where the /statistics endpoint is available to obtain the node height
core_rest_url = "localhost:3003"

# Parent file system We want to bacukp. We also backup all inherited pools
#
# Let's say We have the following zfs pools:
#       NAME
#       vega_pool
#       vega_pool/home
#       vega_pool/home/network-history
#       vega_pool/home/postgresql
#       vega_pool/home/tendermint_home
#
# and we provide` + "`file_system = \"vega_pool\"`" + `, We backup all of the above file system
file_system = "vega_pool"

# Path where the state is saved
state_file = "./backups.json"

# https://github.com/someone1/zfsbackup-go
zfsbackup_binary = "/home/daniel/go/bin/zfsbackup-go"
 
# Configuration for the S3 destination of the backups
[destination]
# The S3 endpoint
endpoint = "fra1.digitaloceanspaces.com"

# S3 region
region = "fra1"

# The S3 bucket name
bucket = "vega-internal-tm-postgres-backups"

# The S3 bucket path. It composes into the followint S3 path: s3://<destination.bucket>/<destination.path>
path = "n00.devnet1.vega.xyz"

# Environment variables names. You must export the environment variables, you define here.
access_key_id_env_name = "AWS_ACCESS_KEY_ID"
access_key_secret_env_name = "AWS_SECRET_ACCESS_KEY"`

var printConfigCmd = &cobra.Command{
	Use:   "print-config",
	Short: "Prints example config",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(exampleConfig)
	},
}

func init() {
	createArgs.BackupArgs = &backupArgs

	BackupCmd.AddCommand(printConfigCmd)
}
