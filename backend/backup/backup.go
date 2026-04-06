package backup

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gorm.io/gorm"
)

type Scheduler struct {
	db        *gorm.DB
	dbPath    string
	backupDir string
	interval  time.Duration
	retain    int
	done      chan struct{}
}

func NewScheduler(db *gorm.DB, dbPath, backupDir string, interval time.Duration, retain int) *Scheduler {
	os.MkdirAll(backupDir, 0755)
	return &Scheduler{
		db: db, dbPath: dbPath, backupDir: backupDir,
		interval: interval, retain: retain, done: make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go func() {
		// Run first backup after 1 minute
		timer := time.NewTimer(1 * time.Minute)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				s.runBackup()
				timer.Reset(s.interval)
			case <-s.done:
				return
			}
		}
	}()
	slog.Info("backup scheduler started", "interval", s.interval.String(), "retain", s.retain)
}

func (s *Scheduler) Close() {
	close(s.done)
}

func (s *Scheduler) runBackup() {
	backupName := fmt.Sprintf("cvbuilder_%s.db", time.Now().Format("20060102_150405"))
	backupPath := filepath.Join(s.backupDir, backupName)

	// Use SQLite VACUUM INTO for safe hot backup
	result := s.db.Exec("VACUUM INTO ?", backupPath)
	if result.Error != nil {
		slog.Error("backup failed", "error", result.Error)
		return
	}
	slog.Info("backup completed", "path", backupPath)

	// Clean old backups
	s.cleanOldBackups()
}

func (s *Scheduler) cleanOldBackups() {
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return
	}
	var backups []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".db" {
			backups = append(backups, filepath.Join(s.backupDir, e.Name()))
		}
	}
	sort.Strings(backups)
	for len(backups) > s.retain {
		os.Remove(backups[0])
		slog.Info("removed old backup", "path", backups[0])
		backups = backups[1:]
	}
}
