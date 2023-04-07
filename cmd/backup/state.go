package backup

import (
	"time"
)

type BackupEntry struct {
	Started time.Time
	Finished time.Time
	ServerHost string
	
}

type State struct {
	LastUpdated time.Time
	Backups BackupEntry
}