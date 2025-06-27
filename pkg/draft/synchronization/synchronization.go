package synchronization

import (
	"slices"
	"time"
)

type LtSrvDTO struct {
	ID        int
	Name      string
	UpdatedAt time.Time // unix milliseconds int64
}

type SrvDTO struct {
	ID        int
	Name      string
	Data      []byte
	CreatedAt time.Time // unix milliseconds int64
	UpdatedAt time.Time // unix milliseconds int64
	Deleted   bool
}

type LtClDTO struct {
	ID        int
	Name      string
	UpdatedAt time.Time // unix milliseconds int64
	SyncID    int
}

type ClDTO struct {
	ID        int
	Name      string
	Data      []byte
	CreatedAt time.Time // unix milliseconds int64
	UpdatedAt time.Time // unix milliseconds int64
	Deleted   bool
	SyncID    int
}

type Server interface {
	GetLtData() []LtSrvDTO
	GetAllData() []SrvDTO
	GetData([]int) []SrvDTO   // slice of ID's
	InsertData([]ClDTO) []int // return ID slice for SyncID
}

type Client interface {
	GetLtData() []LtClDTO
	GetAllData() []ClDTO
	GetData([]int) []ClDTO // slice of ID's
	InsertData([]SrvDTO)   // don't forget about sync_id it SrvDTO.ID
	InsertSyncID([][2]int) // pass key value pairs [ID, syncID] for set sync_id to rows
}

func SyncWithUniqueName(srv Server, cl Client) {
	ltSrvData := srv.GetLtData()
	if len(ltSrvData) == 0 {
		clData := cl.GetAllData()
		sendToServer(srv, cl, clData)
		return
	}

	ltClData := cl.GetLtData()
	if len(ltClData) == 0 {
		srvData := srv.GetAllData()
		saveFromServer(cl, srvData)
		return
	}

	var sToC []int
	var cToS []int
	syncedDataMap, newDataMap := makeClDataMaps(ltClData)
	ltSrvData, sToC, cToS = handleSyncedData(syncedDataMap, ltSrvData, sToC, cToS)
	sToC, cToS = handleNewData(ltSrvData, newDataMap, sToC, cToS)

	if len(sToC) != 0 {
		srvData := srv.GetData(sToC)
		saveFromServer(cl, srvData)
	}

	if len(cToS) != 0 {
		clData := cl.GetData(cToS)
		sendToServer(srv, cl, clData)
	}
}

func makeIDSyncIDPairs(clData []ClDTO, srvID []int) [][2]int {
	var s [][2]int
	for i, clDTO := range clData {
		if clDTO.SyncID == 0 {
			s = append(s, [2]int{clData[i].ID, srvID[i]})
		}
	}
	return s
}

func makeClDataMaps(ltClData []LtClDTO) (map[int]LtClDTO, map[string]LtClDTO) {
	syncedData := make(map[int]LtClDTO)
	newData := make(map[string]LtClDTO)
	for _, v := range ltClData {
		if v.SyncID != 0 {
			syncedData[v.SyncID] = v
			continue
		}
		newData[v.Name] = v
	}
	return syncedData, newData
}

// Retruns LtSrvData, sToC, cTos
func handleSyncedData(
	syncedMap map[int]LtClDTO, ltSrvData []LtSrvDTO, sToC, cToS []int,
) ([]LtSrvDTO, []int, []int) {
	if len(syncedMap) != 0 {
		for i := 0; i < len(ltSrvData); {
			dtoS := ltSrvData[i]
			dtoC, ok := syncedMap[dtoS.ID]
			if !ok {
				i++
				continue
			}
			switch dtoS.UpdatedAt.Compare(dtoC.UpdatedAt) {
			case -1:
				cToS = append(cToS, dtoC.ID)
			case 1:
				sToC = append(sToC, dtoS.ID)
			case 0:
			}
			ltSrvData = ltSrvData[i+1:]
		}
	}
	return ltSrvData, sToC, cToS
}

// Returns sToC, cToS
func handleNewData(
	ltSrvData []LtSrvDTO, newDataMap map[string]LtClDTO, sToC, cToS []int,
) ([]int, []int) {
	for _, dtoS := range ltSrvData {
		if dtoC, ok := newDataMap[dtoS.Name]; ok {
			if dtoS.UpdatedAt.After(dtoC.UpdatedAt) {
				sToC = append(sToC, dtoS.ID)
			} else {
				cToS = append(cToS, dtoC.ID)
			}
			delete(newDataMap, dtoS.Name)
			continue
		}
		sToC = append(sToC, dtoS.ID)
	}

	for _, dtoC := range newDataMap {
		cToS = append(cToS, dtoC.ID)
	}
	slices.Sort(cToS)

	return sToC, cToS
}

func sendToServer(srv Server, cl Client, clData []ClDTO) {
	if len(clData) == 0 {
		return
	}
	syncIDs := srv.InsertData(clData)
	idSyncID := makeIDSyncIDPairs(clData, syncIDs)
	cl.InsertSyncID(idSyncID)
}

func saveFromServer(cl Client, srvData []SrvDTO) {
	if len(srvData) == 0 {
		return
	}
	cl.InsertData(srvData)
}
