package dao

import (
	"LaoQGChat/api/models/chat"
	"database/sql"
	"encoding/json"
	"errors"
)

type ChatDao interface {
	CheckSessionById(sessionId string) (bool, error)
	CheckUserSessionById(userName string, sessionId string) (bool, error)
	GetSessionContentsById(sessionId string) ([]chat.Content, error)
	InsertSessionContents(sessionId string, contents []chat.Content) error
	DeleteSessionById(sessionId string) error
}

type chatDao struct {
	checkSessionByIdStmt       *sql.Stmt
	checkUserSessionByIdStmt   *sql.Stmt
	getSessionContentsByIdStmt *sql.Stmt
	insertSessionContentsStmt  *sql.Stmt
	deleteSessionByIdStmt      *sql.Stmt
}

func NewChatDao(db *sql.DB) (ChatDao, error) {
	var (
		checkSessionByIdStmt       *sql.Stmt
		checkUserSessionByIdStmt   *sql.Stmt
		getSessionContentsByIdStmt *sql.Stmt
		insertSessionContentsStmt  *sql.Stmt
		deleteSessionByIdStmt      *sql.Stmt
		err                        error
	)
	checkSessionByIdStmt, err = db.Prepare(
		`SELECT session_id FROM chat_session
        WHERE session_id = $1 AND delete_flag = false
        FOR UPDATE NOWAIT`)
	if err != nil {
		return nil, err
	}

	checkUserSessionByIdStmt, err = db.Prepare(
		`SELECT session_id FROM chat_session
        WHERE user_name = $1 AND session_id = $2 AND delete_flag = false
        FOR UPDATE NOWAIT`)
	if err != nil {
		return nil, err
	}

	getSessionContentsByIdStmt, err = db.Prepare(
		`SELECT content FROM chat_content
        WHERE session_id = $1 AND delete_flag = false
        ORDER BY serial_number
        FOR UPDATE NOWAIT`)
	if err != nil {
		return nil, err
	}

	insertSessionContentsStmt, err = db.Prepare(
		`INSERT INTO chat_content (session_id, content) VALUES ($1, $2)`)
	if err != nil {
		return nil, err
	}

	deleteSessionByIdStmt, err = db.Prepare(
		`UPDATE chat_session SET delete_flag = true
        WHERE session_id = $1`)
	if err != nil {
		return nil, err
	}

	return chatDao{
		checkSessionByIdStmt:       checkSessionByIdStmt,
		checkUserSessionByIdStmt:   checkUserSessionByIdStmt,
		getSessionContentsByIdStmt: getSessionContentsByIdStmt,
		insertSessionContentsStmt:  insertSessionContentsStmt,
		deleteSessionByIdStmt:      deleteSessionByIdStmt,
	}, nil
}

func (dao chatDao) CheckSessionById(sessionId string) (bool, error) {
	row := dao.checkSessionByIdStmt.QueryRow(sessionId)
	if err := row.Scan(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (dao chatDao) CheckUserSessionById(userName string, sessionId string) (bool, error) {
	row := dao.checkUserSessionByIdStmt.QueryRow(userName, sessionId)
	if err := row.Scan(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (dao chatDao) GetSessionContentsById(sessionId string) ([]chat.Content, error) {
	rows, err := dao.getSessionContentsByIdStmt.Query(sessionId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	contents := make([]chat.Content, 0)
	for rows.Next() {
		var contentBytes []byte
		err = rows.Scan(&contentBytes)
		if err != nil {
			return nil, err
		}
		content := chat.Content{}
		err = json.Unmarshal(contentBytes, &content)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	return contents, nil
}

func (dao chatDao) InsertSessionContents(sessionId string, contents []chat.Content) error {
	for _, content := range contents {
		contentBytes, err := json.Marshal(content)
		if err != nil {
			return err
		}
		_, err = dao.insertSessionContentsStmt.Exec(sessionId, contentBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dao chatDao) DeleteSessionById(sessionId string) error {
	_, err := dao.deleteSessionByIdStmt.Exec(sessionId)
	if err != nil {
		return err
	}
	return nil
}
