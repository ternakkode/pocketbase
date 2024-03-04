package daos

import "github.com/pocketbase/pocketbase/models"

func (dao *Dao) SaveAddress(address *models.Address) error {
	return dao.Save(address)
}
