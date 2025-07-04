package syncservice

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/niksmo/gophkeeper/internal/model"
	"github.com/niksmo/gophkeeper/pkg/logger"
)

type LocalRepo interface {
	GetComparable(context.Context) ([]model.LocalComparable, error)
	GetAll(context.Context) ([]model.LocalPayload, error)
	GetSliceByIDs(ctx context.Context, IDs []int64) ([]model.LocalPayload, error)
	UpdateSliceBySyncIDs(ctx context.Context, data []model.SyncPayload) error
	InsertSlice(ctx context.Context, data []model.LocalPayload) error
	InsertSliceSyncID(ctx context.Context, IDSyncIDPairs [][2]int64) error
}

type ServerClient interface {
	SetToken(string)
	GetComparable(context.Context) ([]model.SyncComparable, error)
	GetAll(context.Context) ([]model.SyncPayload, error)
	GetSliceByIDs(ctx context.Context, IDs []int64) ([]model.SyncPayload, error)
	UpdateSliceByIDs(ctx context.Context, data []model.SyncPayload) error
	InsertSlice(ctx context.Context, data []model.LocalPayload) ([]int64, error)
}

type lists struct {
	insert []int64
	update []int64
}

type Worker struct {
	logger logger.Logger
	local  LocalRepo
	server ServerClient
}

func NewWorker(l logger.Logger, clR LocalRepo, srvR ServerClient) *Worker {
	return &Worker{l, clR, srvR}
}

func (w *Worker) DoJob(ctx context.Context, token string) {
	const op = "Worker.DoJob"
	log := w.logger.WithOp(op)

	sync.OnceFunc(func() {
		log.Debug().Str("token", token).Msg("set token to server client")
		w.server.SetToken(token)
	})()

	srvComp, err := w.getServerComparable(ctx)
	if err != nil {
		return
	}

	if w.serverNoData(srvComp) {
		log.Debug().Msg("server no data")
		locData, err := w.getLocalAll(ctx)
		if err != nil {
			return
		}
		w.insertToServer(ctx, locData)
		return
	}

	locComp, err := w.getLocalComparable(ctx)
	if err != nil {
		return
	}

	if w.localNoData(locComp) {
		log.Debug().Msg("no local data")
		srvData, err := w.getServerAll(ctx)
		if err != nil {
			return
		}
		w.insertToLocal(ctx, srvData)
		return
	}

	log.Debug().Int(
		"locCompLen", len(locComp)).Int(
		"srvCompLen", len(srvComp)).Msg(
		"start compare between local and server")

	srvIDs, locIDs := w.compare(locComp, srvComp)

	log.Debug().Ints64(
		"insertFromServer", srvIDs.insert).Ints64(
		"updateFromServer", srvIDs.update).Ints64(
		"insertFromLocal", locIDs.insert).Ints64(
		"updateFromLocal", locIDs.update).Msg("compare result")

	go w.handleServerData(ctx, srvIDs)
	go w.handleLocalData(ctx, locIDs)
}

func (w *Worker) serverNoData(srvComp []model.SyncComparable) bool {
	return len(srvComp) == 0
}

func (w *Worker) localNoData(locComp []model.LocalComparable) bool {
	return len(locComp) == 0
}

func (w *Worker) insertToServer(
	ctx context.Context, locData []model.LocalPayload,
) {
	const op = "Worker.insertToServer"
	log := w.logger.WithOp(op)

	if len(locData) == 0 {
		log.Debug().Msg("no local data to send")
		return
	}

	log.Debug().Msg("start insert local data to the server")

	syncIDs, err := w.server.InsertSlice(ctx, locData)
	if err != nil {
		log.Error().Err(err).Msg("failed to send local data ot server")
		return
	}

	log.Debug().Ints64("syncIDs", syncIDs).Msg(
		"insert local data to the server successfully")

	if len(syncIDs) != len(locData) {
		log.Error().Err(err).Int(
			"locDataLen", len(locData)).Int(
			"syncIDsLen", len(syncIDs)).Msg(
			"unexpected syncIDs returned length")
		return
	}

	IDSyncIDPairs := w.makeIDSyncIDPairs(locData, syncIDs)

	log.Debug().Any(
		"idSyncIDPairs", IDSyncIDPairs).Msg(
		"start insert syncID to local")

	err = w.local.InsertSliceSyncID(ctx, IDSyncIDPairs)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert syncIDs to local data")
	}
	log.Debug().Msg("insert syncID successfully")
}

func (w *Worker) makeIDSyncIDPairs(
	locData []model.LocalPayload, syncIDs []int64,
) [][2]int64 {
	var s [][2]int64
	for i, syncID := range syncIDs {
		s = append(s, [2]int64{locData[i].ID, syncID})
	}
	return s
}

func (w *Worker) insertToLocal(
	ctx context.Context, srvData []model.SyncPayload,
) error {
	const op = "Worker.InsertToLocal"
	log := w.logger.WithOp(op)

	if len(srvData) == 0 {
		log.Debug().Msg("no server data to insert")
		return nil
	}

	insertData := w.convertSrvToLoc(srvData)

	log.Debug().Msg("start insert server data to local")

	err := w.local.InsertSlice(ctx, insertData)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert data to local")
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Msg("insert server data successfully")
	return nil
}

func (w *Worker) convertSrvToLoc(
	srvData []model.SyncPayload,
) []model.LocalPayload {
	s := make([]model.LocalPayload, 0, len(srvData))
	for _, o := range srvData {
		lp := model.LocalPayload{
			SyncPayload: o,
			SyncID:      o.ID,
		}
		lp.ID = -1
		s = append(s, lp)
	}
	return s
}

func (w *Worker) makeLocalComparableMaps(
	locComp []model.LocalComparable,
) (
	map[int64]model.LocalComparable,
	map[string]model.LocalComparable,
) {
	syncIDModelMap := make(map[int64]model.LocalComparable)
	nameModelMap := make(map[string]model.LocalComparable)

	for _, o := range locComp {
		if o.SyncID != 0 {
			syncIDModelMap[o.SyncID] = o
		} else {
			nameModelMap[o.Name] = o
		}
	}
	return syncIDModelMap, nameModelMap
}

func (w *Worker) compareForUpdate(
	srvComp []model.SyncComparable,
	syncLocalCompMap map[int64]model.LocalComparable,
) (
	fromSrv []int64,
	fromLoc []int64,
	notSyncYet []model.SyncComparable,
) {
	for _, srvObj := range srvComp {
		if locObj, ok := syncLocalCompMap[srvObj.ID]; ok {
			switch locObj.UpdatedAt.Compare(srvObj.UpdatedAt) {
			case -1:
				fromSrv = append(fromSrv, srvObj.ID)
			case 1:
				fromLoc = append(fromLoc, locObj.ID)
			}
			continue
		}
		notSyncYet = append(notSyncYet, srvObj)
	}
	slices.Sort(fromLoc)
	return
}

func (w *Worker) compareForInsert(
	notSyncYet []model.SyncComparable,
	newLocalCompMap map[string]model.LocalComparable,
) (fromSrv []int64, fromLoc []int64) {
	for _, srvObj := range notSyncYet {
		if locObj, ok := newLocalCompMap[srvObj.Name]; ok {
			switch locObj.UpdatedAt.Compare(srvObj.UpdatedAt) {
			case -1:
				fromSrv = append(fromSrv, srvObj.ID)
			case 1:
				fromLoc = append(fromLoc, locObj.ID)
			}
			delete(newLocalCompMap, srvObj.Name)
			continue
		}
		fromSrv = append(fromSrv, srvObj.ID)
	}

	for _, locObj := range newLocalCompMap {
		fromLoc = append(fromLoc, locObj.ID)
	}
	slices.Sort(fromLoc)
	return
}

func (w *Worker) updateLocal(
	ctx context.Context, srvData []model.SyncPayload,
) error {
	const op = "Worker.updateLocal"
	log := w.logger.WithOp(op)

	if len(srvData) == 0 {
		log.Debug().Msg("no data for update local")
		return nil
	}

	log.Debug().Msg("start update local data")
	if err := w.local.UpdateSliceBySyncIDs(ctx, srvData); err != nil {
		log.Error().Err(err).Msg("failed to update local data")
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Msg("update local successfully")
	return nil
}

func (w *Worker) updateServer(
	ctx context.Context, locData []model.LocalPayload,
) error {
	const op = "Worker.updateServer"
	log := w.logger.WithOp(op)

	if len(locData) == 0 {
		log.Debug().Msg("no data for server update")
		return nil
	}

	updateData := w.convertLocToSrv(locData)

	log.Debug().Msg("start update server data")
	err := w.server.UpdateSliceByIDs(ctx, updateData)
	if err != nil {
		log.Error().Err(err).Msg("failed to update server by IDs")
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Msg("update server successfully")
	return nil
}

func (w *Worker) convertLocToSrv(
	locData []model.LocalPayload,
) []model.SyncPayload {
	s := make([]model.SyncPayload, 0, len(locData))

	for _, o := range locData {
		o.SyncPayload.ID = o.SyncID
		s = append(s, o.SyncPayload)
	}
	return s
}

func (w *Worker) getLocalAll(
	ctx context.Context,
) ([]model.LocalPayload, error) {
	const op = "Worker.getLocalAll"
	log := w.logger.WithOp(op)

	log.Debug().Msg("start get all local data")

	locData, err := w.local.GetAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all local data")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int("locDataLen", len(locData)).Msg("get all local data")
	return locData, nil
}

func (w *Worker) getServerAll(
	ctx context.Context,
) ([]model.SyncPayload, error) {
	const op = "Worker.getServerAll"
	log := w.logger.WithOp(op)

	log.Debug().Msg("start get all from server")

	srvData, err := w.server.GetAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all server data")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug().Int(
		"srvDataLen", len(srvData)).Msg(
		"receive data from server")

	return srvData, nil
}

func (w *Worker) getLocalComparable(
	ctx context.Context,
) ([]model.LocalComparable, error) {
	const op = "Worker.getLocalComparable"
	log := w.logger.WithOp(op)

	log.Debug().Msg("start get local comparable")

	locComp, err := w.local.GetComparable(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get local comparable objects")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int("locCompLen", len(locComp)).Msg("get local comparable")
	return locComp, nil
}

func (w *Worker) getServerComparable(
	ctx context.Context,
) ([]model.SyncComparable, error) {
	const op = "Worker.getServerComparable"
	log := w.logger.WithOp(op)

	log.Debug().Msg("start get comparable from server")

	srvComp, err := w.server.GetComparable(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get server comparable objects")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int(
		"srvComLen", len(srvComp)).Msg(
		"receive comparable from server")
	return srvComp, nil
}

func (w *Worker) getLocalSlice(
	ctx context.Context, IDs []int64,
) ([]model.LocalPayload, error) {
	const op = "Worker.getLocalSlice"
	log := w.logger.WithOp(op)

	if len(IDs) == 0 {
		log.Debug().Msg("no IDs for get slice from local")
		return []model.LocalPayload{}, nil
	}

	log.Debug().Ints64("IDs", IDs).Msg("start get slice from local")

	locData, err := w.local.GetSliceByIDs(ctx, IDs)
	if err != nil {
		log.Error().Err(err).Msg("failed to get local data by IDs")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int(
		"locDataLen", len(locData)).Msg(
		"get slice from local successfully")
	return locData, nil
}

func (w *Worker) getServerSlice(
	ctx context.Context, IDs []int64,
) ([]model.SyncPayload, error) {
	const op = "Worker.getLocalSlice"
	log := w.logger.WithOp(op)

	if len(IDs) == 0 {
		log.Debug().Msg("no IDs for get slice from server")
		return []model.SyncPayload{}, nil
	}

	log.Debug().Ints64("IDs", IDs).Msg("start get slice from server")
	srvData, err := w.server.GetSliceByIDs(ctx, IDs)
	if err != nil {
		log.Error().Err(err).Msg("failed to get data from server by IDs")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug().Int(
		"srvDataLen", len(srvData)).Msg(
		"get slice from server successfully")

	return srvData, nil
}

func (w *Worker) handleServerData(
	ctx context.Context, srvIDs lists,
) {
	const op = "Worker.handleServerData"
	log := w.logger.WithOp(op)
	log.Debug().Msg("start op")

	srvData, err := w.getServerSlice(
		ctx, append(srvIDs.update, srvIDs.insert...))
	if err != nil {
		return
	}

	updFromSrvSize := len(srvIDs.update)

	updFromSrvData := srvData[:updFromSrvSize]
	insFromSrvData := srvData[updFromSrvSize:]

	err = w.updateLocal(ctx, updFromSrvData)
	if err != nil {
		return
	}

	w.insertToLocal(ctx, insFromSrvData)
	log.Debug().Msg("end op")
}

func (w *Worker) handleLocalData(
	ctx context.Context, locIDs lists,
) {
	const op = "Worker.handleLocalData"

	log := w.logger.WithOp(op)

	log.Debug().Msg("start op")

	locData, err := w.getLocalSlice(
		ctx, append(locIDs.update, locIDs.insert...))
	if err != nil {
		return
	}

	updFromLocSize := len(locIDs.update)

	updFromLocData := locData[:updFromLocSize]
	insFromLocData := locData[updFromLocSize:]

	err = w.updateServer(ctx, updFromLocData)
	if err != nil {
		return
	}

	w.insertToServer(ctx, insFromLocData)
	log.Debug().Msg("end op")
}

func (w *Worker) compare(
	locComp []model.LocalComparable, srvComp []model.SyncComparable,
) (fromSrvLists, fromLocLists lists) {
	var notSyncYet []model.SyncComparable
	syncLocCompMap, newLocCompMap := w.makeLocalComparableMaps(locComp)

	fromSrvLists.update, fromLocLists.update, notSyncYet =
		w.compareForUpdate(srvComp, syncLocCompMap)

	fromSrvLists.insert, fromLocLists.insert =
		w.compareForInsert(notSyncYet, newLocCompMap)

	return
}
