package socialUser

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"database/sql"
)

// get friend status id from db by string
func getFriendStatusIDFromDB(status string) (int, error) {
	db := pgDB.GetDB()
	var id int
	err := db.QueryRow(`SELECT id FROM friend_statuses WHERE name = $1`, status).Scan(&id)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "get friend status id from db", status, err)
		return 0, appErr.InternalServerError("internal server error")
	}
	return id, nil
}

// update friend status in db
func updateFriendStatusInDB(actorID, targetID uint64, statusID int) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`
		UPDATE friends 
		SET status_id = $1 
		WHERE (friend_1_id = $2 AND friend_2_id = $3) OR (friend_1_id = $3 AND friend_2_id = $2)
	`, statusID, actorID, targetID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "update friend status", map[string]interface{}{
			"actorID": actorID, "targetID": targetID, "statusID": statusID,
		}, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// insert friend data to db
func insertFriendDataToDB(actorID, targetID uint64, statusID int) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`
		INSERT INTO friends (friend_1_id, friend_2_id, status_id)
		VALUES ($1, $2, $3)
	`, actorID, targetID, statusID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "insert friend data", map[string]interface{}{
			"actorID": actorID, "targetID": targetID, "statusID": statusID,
		}, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// remove friend data from db
func removeFriendDataFromDB(actorID, targetID uint64) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`
		DELETE FROM friends
		WHERE 
			(friend_1_id = $1 AND friend_2_id = $2)
			OR (friend_1_id = $2 AND friend_2_id = $1)
	`, actorID, targetID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "remove friend", map[string]interface{}{
			"actorID": actorID, "targetID": targetID,
		}, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// remove friend data by status id from db
func removeFriendDataByStatusIDFromDB(actorID, targetID uint64, statusID int) error {
	db := pgDB.GetDB()
	_, err := db.Exec(`
		DELETE FROM friends 
		WHERE 
			friend_1_id = $1 AND friend_2_id = $2
			AND status_id = $3
	`, actorID, targetID, statusID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "remove friend data by status id", map[string]interface{}{
			"actorID": actorID, "targetID": targetID, "statusID": statusID,
		}, err)
		return appErr.InternalServerError("inernal server error")
	}
	return nil
}

// get relations between two users from db
func getRelationsFromDB(actorID, targetID uint64) (*string, error) {
	db := pgDB.GetDB()
	query := `
		SELECT 
			CASE 
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted') THEN 'accepted'
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked') 
					AND f.friend_1_id = $1 THEN 'blocked'
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked') 
					AND f.friend_2_id = $1 THEN 'blocked by target'
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request') 
					AND f.friend_1_id = $1 THEN 'outgoing request'
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request') 
					AND f.friend_2_id = $1 THEN 'incoming request'
				ELSE NULL
			END AS friendship_status
		FROM friends f
		WHERE 
			(f.friend_1_id = $1 AND f.friend_2_id = $2)
			OR (f.friend_1_id = $2 AND f.friend_2_id = $1)
		ORDER BY 
			CASE 
				WHEN f.friend_1_id = $1 AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked') THEN 1
				WHEN f.friend_2_id = $1 AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked') THEN 2
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted') THEN 3
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request') THEN 4
				ELSE 5
			END
		LIMIT 1
	`

	var friendshipStatus string
	err := db.QueryRow(query, actorID, targetID).Scan(&friendshipStatus)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get relations between users", map[string]interface{}{
			"userID": actorID, "targetID": targetID,
		}, err)
		return nil, appErr.InternalServerError("internal server error")
	}

	return &friendshipStatus, nil
}
