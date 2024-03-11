package db

import "github.com/busyLambda/bbkk/internal/models"

func (d *DbManager) InsertSession(s models.Session) error {
	return d.Conn.Create(s).Error
}

func (d *DbManager) DeleteSession(id uint) error {
	s, err := d.GetSessionById(id)
	if err != nil {
		return err
	}

	return d.Conn.Delete(&s).Error
}

func (d *DbManager) GetSessionById(id uint) (s models.Session, err error) {
	err = d.Conn.Where("id = ?", id).First(&s).Error
	return
}
