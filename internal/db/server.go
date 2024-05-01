package db

import "github.com/busyLambda/bbkk/internal/models"

func (d *DbManager) InsertServer(s *models.Server) error {
	return d.Conn.Create(s).Error
}

func (d *DbManager) DeleteServer(id int) error {
	s, err := d.GetServerByID(id)
	if err != nil {
		return err
	}

	return d.Conn.Delete(&s).Error
}

func (d *DbManager) GetServerByID(id int) (s models.Server, err error) {
	err = d.Conn.Where("id = ?", id).First(&s).Error
	return
}

func (d *DbManager) GetServerByName(sn string) (s models.Server, err error) {
	err = d.Conn.Where("name = ?", sn).First(&s).Error
	return
}

func (d *DbManager) GetAllServers() (servers []models.Server, err error) {
	err = d.Conn.Find(&servers).Error
	return
}
