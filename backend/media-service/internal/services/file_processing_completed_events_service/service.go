package fileprocessingcompletedeventsservice

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/event"
	producersinterfaces "github.com/ZaiiiRan/messenger/backend/media-service/internal/producers/interfaces"
	producersmodels "github.com/ZaiiiRan/messenger/backend/media-service/internal/producers/models"
	uow "github.com/ZaiiiRan/messenger/backend/media-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type FileProcessingCompletedEventsService interface {
	CreateFileProcessingCompletedEvents(ctx context.Context, payloads []*producersmodels.FileProcessingCompletedEvent, uow *uow.UnitOfWork) error
	SendFileProcessingCompletedEvents(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int) error
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *fileProcessingCompletedEventsDataProvider
	producer     producersinterfaces.FileProcessingCompletedEventsProducer
}

func New(
	pgClient *postgres.PostgresClient,
	producer producersinterfaces.FileProcessingCompletedEventsProducer,
	log *zap.SugaredLogger,
) FileProcessingCompletedEventsService {
	return &service{
		log:          log,
		dataProvider: newFileProcessingCompletedEventsDataProvider(pgClient),
		producer:     producer,
	}
}

func (s *service) CreateFileProcessingCompletedEvents(
	ctx context.Context,
	payloads []*producersmodels.FileProcessingCompletedEvent,
	uow *uow.UnitOfWork,
) error {
	l := s.log.With("op", "create_file_processing_completed_events", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(payloads) == 0 {
		return nil
	}

	events := make([]*event.Event, 0, len(payloads))
	for _, payload := range payloads {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			l.Errorw("file_processing_completed_events.create_file_processing_completed_events_failed.marshal_error", "err", err, "payload", payload)
			return ErrCreateFileProcessingCompletedEvents
		}
		events = append(events, event.New("", jsonPayload))
	}

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("file_processing_completed_events.create_file_processing_completed_events_failed.begin_transaction_error", "err", err)
			return ErrCreateFileProcessingCompletedEvents
		}
	}

	if err := s.dataProvider.createEvents(ctx, events, uow); err != nil {
		l.Errorw("file_processing_completed_events.create_file_processing_completed_events_failed.create_error", "err", err)
		return ErrCreateFileProcessingCompletedEvents
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("file_processing_completed_events.create_file_processing_completed_events_failed.commit_error", "err", err)
			return ErrCreateFileProcessingCompletedEvents
		}
	}

	l.Infow("file_processing_completed_events.create_file_processing_completed_events.success", "count", len(events))
	return nil
}

func (s *service) SendFileProcessingCompletedEvents(
	ctx context.Context,
	workerID string,
	retryIntervalMS uint,
	batchSize int,
) error {
	l := s.log.With("op", "send_file_processing_completed_events", "worker_id", workerID)

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("file_processing_completed_events.send_file_processing_completed_events_failed.begin_transaction_error", "err", err)
		return ErrSendFileProcessingCompletedEvents
	}

	now := time.Now()
	retryAfter := now.Add(-1 * time.Duration(retryIntervalMS) * time.Millisecond)
	outboxEvents, err := s.dataProvider.getEventsLocked(ctx, batchSize, retryAfter, uow)
	if err != nil {
		l.Errorw("file_processing_completed_events.send_file_processing_completed_events_failed.get_events_error", "err", err)
		return ErrSendFileProcessingCompletedEvents
	}
	if outboxEvents == nil {
		return nil
	}

	outboxEventsFailed := make([]*event.Event, 0, len(outboxEvents))
	outboxEventsSuccess := make([]*event.Event, 0, len(outboxEvents))

	for _, evt := range outboxEvents {
		var payload producersmodels.FileProcessingCompletedEvent
		if err := json.Unmarshal(evt.GetPayload(), &payload); err != nil {
			l.Errorw(
				"file_processing_completed_events.send_file_processing_completed_events_failed.unmarshal_payload_error",
				"err", err,
				"event", evt.GetID(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			if err := s.markFailed(evt, now, l); err != nil {
				continue
			}
			outboxEventsFailed = append(outboxEventsFailed, evt)
			continue
		}
		payload.Id = evt.GetID()

		if err := s.producer.ProduceFileProcessingCompletedEvent(ctx, &payload); err != nil {
			l.Errorw(
				"file_processing_completed_events.send_file_processing_completed_events_failed.produce_error",
				"err", err,
				"event", evt.GetID(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			if err := s.markFailed(evt, now, l); err != nil {
				continue
			}
			outboxEventsFailed = append(outboxEventsFailed, evt)
			continue
		}

		outboxEventsSuccess = append(outboxEventsSuccess, evt)
	}

	if len(outboxEventsSuccess) > 0 {
		if err := s.dataProvider.deleteEvents(ctx, outboxEventsSuccess, uow); err != nil {
			l.Errorw("file_processing_completed_events.send_file_processing_completed_events_failed.delete_error", "err", err)
			return ErrSendFileProcessingCompletedEvents
		}
	}
	if len(outboxEventsFailed) > 0 {
		if err := s.dataProvider.updateEvents(ctx, outboxEventsFailed, uow); err != nil {
			l.Errorw("file_processing_completed_events.send_file_processing_completed_events_failed.update_error", "err", err)
			return ErrSendFileProcessingCompletedEvents
		}
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("file_processing_completed_events.send_file_processing_completed_events_failed.commit_error", "err", err)
		return ErrSendFileProcessingCompletedEvents
	}

	if len(outboxEventsSuccess) > 0 {
		l.Infow("file_processing_completed_events.send_file_processing_completed_events.success", "sent", len(outboxEventsSuccess))
	}
	if len(outboxEventsFailed) > 0 {
		l.Warnw("file_processing_completed_events.send_file_processing_completed_events.partial", "not_sent", len(outboxEventsFailed))
	}

	return nil
}

func (s *service) markFailed(evt *event.Event, now time.Time, log *zap.SugaredLogger) error {
	if err := evt.IncrementAttempts(); err != nil {
		log.Errorw(
			"file_processing_completed_events.mark_event_as_failed_error",
			"err", err,
			"event", evt.GetID(),
			"attempts", evt.GetAttempts(),
			"status", evt.GetStatus(),
		)
		return err
	}
	evt.SetUpdatedAt(&now)
	if err := evt.SetStatus(event.EventStatusFailed); err != nil {
		log.Errorw(
			"file_processing_completed_events.mark_event_as_failed_error",
			"err", err,
			"event", evt.GetID(),
			"attempts", evt.GetAttempts(),
			"status", evt.GetStatus(),
		)
		return err
	}
	return nil
}
