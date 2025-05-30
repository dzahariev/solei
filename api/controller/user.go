package controller

import (
	"github.com/dzahariev/solei/api/model"
	"github.com/gofrs/uuid/v5"
)

// DBLoadUser loads an user by given ID
func (server *Server) DBLoadUser(userID string) (*model.User, error) {
	uid, err := uuid.FromString(userID)
	if err != nil {
		return nil, err
	}

	user := &model.User{}
	err = user.FindByID(server.DB, user, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DBSaveUser is caled to save an user
func (server *Server) DBSaveUser(user *model.User) error {
	err := user.Save(server.DB, user)

	if err != nil {
		return err
	}
	return nil
}
