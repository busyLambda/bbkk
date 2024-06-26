package db

import "github.com/busyLambda/bbkk/internal/models"

func (d *DbManager) InsertSession(s *models.Session) error {
	return d.Conn.Create(s).Error
}

func (d *DbManager) DeleteSession(id string) error {
	s, err := d.GetSessionById(id)
	if err != nil {
		return err
	}

	return d.Conn.Delete(&s).Error
}

func (d *DbManager) GetSessionById(id string) (s models.Session, err error) {
	err = d.Conn.Where("id = ?", id).First(&s).Error
	return
}

func (d *DbManager) GetSessionByUser(u uint) (s models.Session, err error) {
	err = d.Conn.Where("user_id = ?", u).First(&s).Error
	return
}
