package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"time"
)

type SessionsRepository struct {
	db db.DB
}

func NewSessionsRepository(db db.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) ReadSessions() ([]model.Session, error) {
	records, err := u.db.Load("sessions")
	if err != nil {
		return nil, err
	}

	var listSessions []model.Session
	err = json.Unmarshal([]byte(records), &listSessions)
	if err != nil {
		return nil, err
	}

	return listSessions, nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}

	newSessionList := make([]model.Session, 0, len(listSessions))
	for i, ses := range listSessions {
		if ses.Token == tokenTarget {
			newSessionList = append(listSessions[:i], listSessions[i+1:]...)
			break
		}
	}

	jsonData, err := json.Marshal(newSessionList)
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	var dataSessions []model.Session

	jsonData, err := u.db.Load("sessions")
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &dataSessions)
	if err != nil {
		return err
	}

	dataSessions = append(dataSessions, session)

	data, err := json.Marshal(dataSessions)
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", data)
	return err
}

func (u *SessionsRepository) CheckExpireToken(token string) (model.Session, error) {
	sessions, err := u.db.Load("sessions")
	if err != nil {
		return model.Session{}, err
	}

	var listSessions []model.Session

	err = json.Unmarshal(sessions, &listSessions)
	if err != nil {
		return model.Session{}, err
	}

	for _, ses := range listSessions {
		if ses.Token == token {
			if u.TokenExpired(ses) {
				err = u.DeleteSessions(token)
				if err != nil {
					return model.Session{}, err
				}
				return model.Session{}, fmt.Errorf("Token is Expired!")
			}
			return ses, nil
		}
	}

	return model.Session{}, nil
}

func (u *SessionsRepository) ResetSessions() error {
	err := u.db.Reset("sessions", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) TokenExist(list []model.Session, req string) (model.Session, error) {
	for _, ses := range list {
		if ses.Token == req {
			return ses, nil
		}
	}
	return model.Session{}, fmt.Errorf("Token Not Found!")
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
