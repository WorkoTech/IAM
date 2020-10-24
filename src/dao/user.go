package dao

import (
    "worko.tech/iam/src/models"

    "github.com/jinzhu/gorm"
)

type UserDao struct {
    db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
    return &UserDao {
        db: db,
    }
}

func (dao *UserDao) GetByEmail(email string) (*models.User, error) {
    var user models.User

    if err := dao.db.Where("email = ?", email).First(&user).Error; err != nil {
        if gorm.IsRecordNotFoundError(err) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

func (dao *UserDao) Create(user *models.User) (err error) {
    return dao.db.Create(user).Error
}

func (dao *UserDao) GetById(id int64) (*models.User, error) {
    var user models.User

    if err := dao.db.Where("id = ?", id).First(&user).Error; err != nil {
        if gorm.IsRecordNotFoundError(err) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

func (dao *UserDao) GetByIds(ids []int64) (*[]models.User, error) {
    var users []models.User

    if err := dao.db.Where("id IN (?)", ids).Find(&users).Error; err != nil {
        if gorm.IsRecordNotFoundError(err) {
            return nil, nil
        }
        return nil, err
    }

    return &users, nil
}
