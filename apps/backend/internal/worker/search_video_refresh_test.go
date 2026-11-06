package worker

import (
	"errors"
	"testing"
	"time"

	"tunetrend-backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRefreshUsecase struct {
	domain.SearchUsecase
	n   int
	err error
}

func (s *stubRefreshUsecase) RefreshSearchVideoStats() (int, error) { return s.n, s.err }

type stubLogRepo struct {
	domain.WorkerLogRepository
	logs     []domain.WorkerLog
	lockKeys []int64
}

func (r *stubLogRepo) CreateLog(l domain.WorkerLog) error {
	r.logs = append(r.logs, l)
	return nil
}

func (r *stubLogRepo) WithLock(key int64, fn func() error) (bool, error) {
	r.lockKeys = append(r.lockKeys, key)
	return true, fn()
}

func TestRunSearchVideoRefresh_LogsSuccessAndFailure(t *testing.T) {
	logs := &stubLogRepo{}
	runSearchVideoRefresh(&stubRefreshUsecase{n: 3}, logs)
	runSearchVideoRefresh(&stubRefreshUsecase{err: errors.New("yt down")}, logs)

	require.Len(t, logs.logs, 2)
	assert.Equal(t, []int64{searchVideoRefreshLockKey, searchVideoRefreshLockKey}, logs.lockKeys)
	assert.Equal(t, domain.SearchVideoRefreshJobName, logs.logs[0].JobName)
	assert.Equal(t, domain.WorkerLogStatusSuccess, logs.logs[0].Status)
	assert.Contains(t, logs.logs[0].Message, "3")
	assert.Equal(t, int(searchVideoRefreshInterval/time.Minute), logs.logs[0].IntervalMinutes)
	assert.Equal(t, domain.WorkerLogStatusFailed, logs.logs[1].Status)
	assert.Equal(t, "yt down", logs.logs[1].Message)
	assert.Equal(t, 6*time.Hour, searchVideoRefreshInterval)
}
