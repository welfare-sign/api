package service

import (
	"context"

	"go.uber.org/zap"

	"welfare-sign/internal/apicode"
	"welfare-sign/internal/pkg/log"
	"welfare-sign/internal/pkg/wsgin"
)

// FailureIssueRecord 失效到达指定时间内，用户未核销完的福利
func (s *Service) FailureIssueRecord(ctx context.Context) (wsgin.APICode, error) {
	issueRecords, err := s.dao.GetNeedClearIssueRecords(ctx)
	if err != nil {
		log.Error(ctx, "FailureIssueRecord.GetNeedClearIssueRecords() error", zap.Error(err))
		return apicode.ErrTaskFailureIssueRecord, err
	}
	if len(issueRecords) == 0 {
		return wsgin.APICodeSuccess, nil
	}

	for i := 0; i < len(issueRecords); i++ {
		if err := s.dao.FailureIssueRecord(ctx, issueRecords[i]); err != nil {
			log.Error(ctx, "FailureIssueRecord error", zap.Error(err))
		}
	}
	return wsgin.APICodeSuccess, nil
}
