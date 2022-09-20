package cleanup

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/serverlock"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/annotations"
	"github.com/grafana/grafana/pkg/services/dashboardsnapshots"
	dashver "github.com/grafana/grafana/pkg/services/dashboardversion"
	"github.com/grafana/grafana/pkg/services/loginattempt"
	"github.com/grafana/grafana/pkg/services/ngalert/image"
	"github.com/grafana/grafana/pkg/services/queryhistory"
	"github.com/grafana/grafana/pkg/services/shorturls"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	tempuser "github.com/grafana/grafana/pkg/services/temp_user"
	"github.com/grafana/grafana/pkg/setting"
)

func ProvideService(cfg *setting.Cfg, serverLockService *serverlock.ServerLockService,
	shortURLService shorturls.Service, sqlstore *sqlstore.SQLStore, queryHistoryService queryhistory.Service,
	dashboardVersionService dashver.Service, dashSnapSvc dashboardsnapshots.Service, deleteExpiredImageService *image.DeleteExpiredService,
	loginAttemptService loginattempt.Service, tempUserService tempuser.Service, tracer tracing.Tracer) *CleanUpService {
	s := &CleanUpService{
		Cfg:                       cfg,
		ServerLockService:         serverLockService,
		ShortURLService:           shortURLService,
		QueryHistoryService:       queryHistoryService,
		store:                     sqlstore,
		log:                       log.New("cleanup"),
		dashboardVersionService:   dashboardVersionService,
		dashboardSnapshotService:  dashSnapSvc,
		deleteExpiredImageService: deleteExpiredImageService,
		loginAttemptService:       loginAttemptService,
		tempUserService:           tempUserService,
		tracer:                    tracer,
	}
	return s
}

type CleanUpService struct {
	log                       log.Logger
	tracer                    tracing.Tracer
	store                     sqlstore.Store
	Cfg                       *setting.Cfg
	ServerLockService         *serverlock.ServerLockService
	ShortURLService           shorturls.Service
	QueryHistoryService       queryhistory.Service
	dashboardVersionService   dashver.Service
	dashboardSnapshotService  dashboardsnapshots.Service
	deleteExpiredImageService *image.DeleteExpiredService
	loginAttemptService       loginattempt.Service
	tempUserService           tempuser.Service
}

type cleanUpJob struct {
	name string
	fn   func(context.Context, log.Logger)
}

func (j cleanUpJob) String() string {
	return strconv.Quote(j.name)
}

func (srv *CleanUpService) Run(ctx context.Context) error {
	srv.cleanUpTmpFiles(ctx, srv.log)

	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			srv.clean(ctx)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (srv *CleanUpService) clean(ctx context.Context) {
	const timeout = 9 * time.Minute
	start := time.Now()
	ctx, span := srv.tracer.Start(ctx, "cleanup background job")
	defer span.End()
	tID := tracing.TraceIDFromContext(ctx, false)
	logger := srv.log.New("traceID", tID)
	ctx, cancelFn := context.WithTimeout(ctx, timeout)
	defer cancelFn()

	cleanupJobs := []cleanUpJob{
		{"clean up temporary files", srv.cleanUpTmpFiles},
		{"delete expired snapshots", srv.deleteExpiredSnapshots},
		{"delete expired dashboard versions", srv.deleteExpiredDashboardVersions},
		{"delete expired images", srv.deleteExpiredImages},
		{"cleanup old annotations", srv.cleanUpOldAnnotations},
		{"expire old user invites", srv.expireOldUserInvites},
		{"delete stale short URLs", srv.deleteStaleShortURLs},
		{"delete stale query history", srv.deleteStaleQueryHistory},
		{"delete old login attempts", srv.deleteOldLoginAttempts},
	}

	logger.Debug("Starting cleanup jobs", "jobs", fmt.Sprintf("%v", cleanupJobs))

	for _, j := range cleanupJobs {
		if ctx.Err() != nil {
			logger.Error("Cancelled cleanup job", "error", ctx.Err(), "duration", time.Since(start))
			return
		}
		ctx, span := srv.tracer.Start(ctx, j.name)
		j.fn(ctx, logger)
		span.End()
	}

	logger.Info("Completed cleanup jobs", "duration", time.Since(start))
}

func (srv *CleanUpService) cleanUpOldAnnotations(ctx context.Context, logger log.Logger) {
	cleaner := annotations.GetAnnotationCleaner()
	affected, affectedTags, err := cleaner.CleanAnnotations(ctx, srv.Cfg)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		logger.Error("failed to clean up old annotations", "error", err)
	} else {
		logger.Debug("Deleted excess annotations", "annotations affected", affected, "annotation tags affected", affectedTags)
	}
}

func (srv *CleanUpService) cleanUpTmpFiles(ctx context.Context, logger log.Logger) {
	folders := []string{
		srv.Cfg.ImagesDir,
		srv.Cfg.CSVsDir,
	}

	for _, f := range folders {
		ctx, span := srv.tracer.Start(ctx, fmt.Sprintf("removing old temporary items from directory: %s", strconv.Quote(f)))
		srv.cleanUpTmpFolder(ctx, logger, f)
		span.End()
	}
}

func (srv *CleanUpService) cleanUpTmpFolder(ctx context.Context, logger log.Logger, folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		logger.Error("Problem reading dir", "folder", folder, "error", err)
		return
	}

	var toDelete []fs.DirEntry
	var now = time.Now()

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			logger.Error("Problem reading file", "folder", folder, "file", file, "error", err)
			continue
		}

		if srv.shouldCleanupTempFile(info.ModTime(), now) {
			toDelete = append(toDelete, file)
		}
	}

	for _, file := range toDelete {
		fullPath := path.Join(folder, file.Name())
		err := os.Remove(fullPath)
		if err != nil {
			logger.Error("Failed to delete temp file", "file", file.Name(), "error", err)
		}
	}

	logger.Debug("Found old rendered file to delete", "folder", folder, "deleted", len(toDelete), "kept", len(files))
}

func (srv *CleanUpService) shouldCleanupTempFile(filemtime time.Time, now time.Time) bool {
	if srv.Cfg.TempDataLifetime == 0 {
		return false
	}

	return filemtime.Add(srv.Cfg.TempDataLifetime).Before(now)
}

func (srv *CleanUpService) deleteExpiredSnapshots(ctx context.Context, logger log.Logger) {
	cmd := dashboardsnapshots.DeleteExpiredSnapshotsCommand{}
	if err := srv.dashboardSnapshotService.DeleteExpiredSnapshots(ctx, &cmd); err != nil {
		logger.Error("Failed to delete expired snapshots", "error", err.Error())
	} else {
		logger.Debug("Deleted expired snapshots", "rows affected", cmd.DeletedRows)
	}
}

func (srv *CleanUpService) deleteExpiredDashboardVersions(ctx context.Context, logger log.Logger) {
	cmd := dashver.DeleteExpiredVersionsCommand{}
	if err := srv.dashboardVersionService.DeleteExpired(ctx, &cmd); err != nil {
		logger.Error("Failed to delete expired dashboard versions", "error", err.Error())
	} else {
		logger.Debug("Deleted old/expired dashboard versions", "rows affected", cmd.DeletedRows)
	}
}

func (srv *CleanUpService) deleteExpiredImages(ctx context.Context, logger log.Logger) {
	if !srv.Cfg.UnifiedAlerting.IsEnabled() {
		return
	}
	if rowsAffected, err := srv.deleteExpiredImageService.DeleteExpired(ctx); err != nil {
		logger.Error("Failed to delete expired images", "error", err.Error())
	} else {
		logger.Debug("Deleted expired images", "rows affected", rowsAffected)
	}
}

func (srv *CleanUpService) deleteOldLoginAttempts(ctx context.Context, logger log.Logger) {
	err := srv.ServerLockService.LockAndExecute(ctx, "delete old login attempts",
		time.Minute*10, func(context.Context) {
			srv.deleteOldLoginAttemptsWithoutLock(ctx, logger)
		})
	if err != nil {
		logger.Error("failed to lock and execute cleanup of old login attempts", "error", err)
	}
}

func (srv *CleanUpService) deleteOldLoginAttemptsWithoutLock(ctx context.Context, logger log.Logger) {
	if srv.Cfg.DisableBruteForceLoginProtection {
		return
	}

	cmd := models.DeleteOldLoginAttemptsCommand{
		OlderThan: time.Now().Add(time.Minute * -10),
	}
	if err := srv.loginAttemptService.DeleteOldLoginAttempts(ctx, &cmd); err != nil {
		logger.Error("Problem deleting expired login attempts", "error", err.Error())
	} else {
		logger.Debug("Deleted expired login attempts", "rows affected", cmd.DeletedRows)
	}
}

func (srv *CleanUpService) expireOldUserInvites(ctx context.Context, logger log.Logger) {
	maxInviteLifetime := srv.Cfg.UserInviteMaxLifetime

	cmd := models.ExpireTempUsersCommand{
		OlderThan: time.Now().Add(-maxInviteLifetime),
	}

	if err := srv.tempUserService.ExpireOldUserInvites(ctx, &cmd); err != nil {
		logger.Error("Problem expiring user invites", "error", err.Error())
	} else {
		logger.Debug("Expired user invites", "rows affected", cmd.NumExpired)
	}
}

func (srv *CleanUpService) deleteStaleShortURLs(ctx context.Context, logger log.Logger) {
	cmd := models.DeleteShortUrlCommand{
		OlderThan: time.Now().Add(-time.Hour * 24 * 7),
	}
	if err := srv.ShortURLService.DeleteStaleShortURLs(ctx, &cmd); err != nil {
		logger.Error("Problem deleting stale short urls", "error", err.Error())
	} else {
		logger.Debug("Deleted short urls", "rows affected", cmd.NumDeleted)
	}
}

func (srv *CleanUpService) deleteStaleQueryHistory(ctx context.Context, logger log.Logger) {
	// Delete query history from 14+ days ago with exception of starred queries
	maxQueryHistoryLifetime := time.Hour * 24 * 14
	olderThan := time.Now().Add(-maxQueryHistoryLifetime).Unix()
	rowsCount, err := srv.QueryHistoryService.DeleteStaleQueriesInQueryHistory(ctx, olderThan)
	if err != nil {
		logger.Error("Problem deleting stale query history", "error", err.Error())
	} else {
		logger.Debug("Deleted stale query history", "rows affected", rowsCount)
	}

	// Enforce 200k limit for query_history table
	queryHistoryLimit := 200000
	rowsCount, err = srv.QueryHistoryService.EnforceRowLimitInQueryHistory(ctx, queryHistoryLimit, false)
	if err != nil {
		logger.Error("Problem with enforcing row limit for query_history", "error", err.Error())
	} else {
		logger.Debug("Enforced row limit for query_history", "rows affected", rowsCount)
	}

	// Enforce 150k limit for query_history_star table
	queryHistoryStarLimit := 150000
	rowsCount, err = srv.QueryHistoryService.EnforceRowLimitInQueryHistory(ctx, queryHistoryStarLimit, true)
	if err != nil {
		logger.Error("Problem with enforcing row limit for query_history_star", "error", err.Error())
	} else {
		logger.Debug("Enforced row limit for query_history_star", "rows affected", rowsCount)
	}
}
