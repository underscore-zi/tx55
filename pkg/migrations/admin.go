package migrations

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/restapi/admin"
)

func initAdminDB(db *gorm.DB) (err error) {
	Logger.WithField("type", AdminDBMigrationType).Info("Initialization complete")
	if err = firstUserInit(db); err != nil {
		Logger.WithError(err).WithField("type", AdminDBMigrationType).Error("Error initializing database")
		return
	}
	Logger.WithField("type", AdminDBMigrationType).Info("Initialization complete")
	return
}

func firstUserInit(db *gorm.DB) (err error) {
	var adminRole admin.Role

	// Make sure we have a super-user role
	Logger.Info("Checking for super-user role")
	if err = db.First(&adminRole, "all_privileges = 1").Error; err == gorm.ErrRecordNotFound {
		Logger.Debug("Creating new super-user role")
		adminRole.Name = "SuperUser"
		adminRole.AllPrivileges = true
		if err = db.Create(&adminRole).Error; err != nil {
			return
		}
	} else if err != nil {
		return
	} else {
		Logger.WithField("name", adminRole.Name).Info("Found super-user role")
	}

	// Make sure we have atleast one user with the super-user role
	var adminUser admin.User
	Logger.Info("Checking for super-user user")
	if err = db.First(&adminUser, "role_id = ?", adminRole.ID).Error; err == gorm.ErrRecordNotFound {
		adminUser.Username = "admin"
		adminUser.Password = StringWithCharset(16, charsetAlphaNumSpecial)
		adminUser.RoleID = adminRole.ID

		Logger.WithFields(logrus.Fields{
			"username": adminUser.Username,
			"password": adminUser.Password,
		}).Error("Creating new admin user")
		if err = db.Create(&adminUser).Error; err != nil {
			return
		}
	} else if err != nil {
		return
	} else {
		Logger.WithField("username", adminUser.Username).Info("Found super-user user")
	}
	return
}
