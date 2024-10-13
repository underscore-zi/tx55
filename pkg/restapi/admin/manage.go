package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelAdmin, "GET", "/admin/users/list", GetUsersAndRoles)
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/users/update", ManageUsers)
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/roles/update", ManageRoles)
}

type ArgsManageUsers struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   uint   `json:"role_id"`
	Delete   bool   `json:"delete"`
}

// ManageUsers godoc
// @Summary      Manage Administrative User
// @Description  Change an Administrative User's username, password or role ID. If the UserID is 0, a new user will be created.
// @Tags         AdminLogin
// @Accept       json
// @Produce      json
// @Param        body     body  ArgsManageUsers  true  "New user information"
// @Success      200  {object}  restapi.ResponseJSON{data=User}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/users/update [post]
// @Security ApiKeyAuth
func ManageUsers(c *gin.Context) {
	if !CheckPrivilege(c, PrivManageUsers) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}
	l := c.MustGet("logger").(*logrus.Logger)

	var user User
	var args ArgsManageUsers
	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, err.Error())
		return
	}

	adminDB := c.MustGet("adminDB").(*gorm.DB)
	// Don't allow handling out privileged roles without the AllPrivileges flag
	var role Role
	if tx := adminDB.First(&role, args.RoleID); tx.Error != nil {
		l.WithError(tx.Error).Error("failed to find role")
		restapi.Error(c, 500, "database error")
		return
	}
	if role.AllPrivileges && !CheckPrivilege(c, PrivAll) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}

	if args.UserID == 0 {
		// New User
		if args.Username == "" || args.Password == "" {
			restapi.Error(c, 400, "username and password are required")
			return
		}

		user.RoleID = args.RoleID
		user.Username = args.Username
		user.Password = args.Password
		if tx := adminDB.Create(&user); tx.Error != nil {
			l.WithError(tx.Error).Error("failed to create admin user")
			restapi.Error(c, 400, tx.Error.Error())
			return
		}
		restapi.Success(c, user)
	} else {
		user.ID = args.UserID
		adminDB.Model(&user).Joins("Role").First(&user)
		if user.HasPrivilege(PrivAll) {
			restapi.Error(c, 400, "cannot modify users with all privileges")
			return
		}

		if args.Delete {
			if tx := adminDB.Delete(&user); tx.Error != nil {
				restapi.Error(c, 400, tx.Error.Error())
				return
			}
		} else {
			var updates User;
			if args.Username != "" {
				updates.Username = args.Username
			}
			if args.Password != "" {
				updates.Password = args.Password
			}
			if args.RoleID > 0 {
				updates.RoleID = args.RoleID
			}
			// save the updated user
			if tx := adminDB.Model(&user).Updates(&updates); tx.Error != nil {
				restapi.Error(c, 400, tx.Error.Error())
				return
			}
		}
		restapi.Success(c, user)
	}
}

type UserRolesJSON struct {
	Users []User `json:"users"`
	Roles []Role `json:"roles"`
}

// GetUsersAndRoles godoc
// @Summary      Get Users/Roles
// @Description  Get all administrative Users/Roles
// @Tags         AdminLogin
// @Produce      json
// @Success      200  {object}  restapi.ResponseJSON{data=User}
// @Failure      403  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/users/list [post]
// @Security ApiKeyAuth
func GetUsersAndRoles(c *gin.Context) {
	if !CheckPrivilege(c, PrivManageUsers) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}
	var users []User
	var roles []Role
	l := c.MustGet("logger").(*logrus.Logger)
	adminDB := c.MustGet("adminDB").(*gorm.DB)

	if tx := adminDB.Find(&users); tx.Error != nil {
		l.WithError(tx.Error).Error("failed to get users")
		restapi.Error(c, 400, "database error")
		return
	}

	if tx := adminDB.Find(&roles); tx.Error != nil {
		l.WithError(tx.Error).Error("failed to get roles")
		restapi.Error(c, 400, "database error")
		return
	}

	restapi.Success(c, UserRolesJSON{Users: users, Roles: roles})
}

type ArgsManageRoles struct {
	Role   Role `json:"role"`
	Delete bool `json:"delete"`
}

// ManageRoles godoc
// @Summary      Manage Administrative Roles
// @Description  Update the privielges granted by a particular role
// @Tags         AdminLogin
// @Accept       json
// @Produce      json
// @Param        body     body  ArgsManageRoles  true  "New Role information"
// @Success      200  {object}  restapi.ResponseJSON{data=User}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      403  {object}  restapi.ResponseJSON{data=string}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/roles/update [post]
// @Security ApiKeyAuth
func ManageRoles(c *gin.Context) {
	if !CheckPrivilege(c, PrivManageUsers) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}

	adminDB := c.MustGet("adminDB").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var args ArgsManageRoles
	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, err.Error())
		return
	}

	changeRequiresAllPrivs := false
	if args.Role.AllPrivileges {
		changeRequiresAllPrivs = true
	} else if args.Role.ID != 0 {
		var original Role
		if tx := adminDB.First(&original, args.Role.ID); tx.Error != nil {
			l.WithError(tx.Error).Error("failed to get role")
			restapi.Error(c, 400, "database error")
			return
		}
		changeRequiresAllPrivs = original.AllPrivileges
	}

	if changeRequiresAllPrivs && !CheckPrivilege(c, PrivAll) {
		restapi.Error(c, 403, "You cannot modify roles with all-privileges")
		return
	}

	if args.Delete {
		var users []User
		if args.Role.ID > 0 {
			if tx := adminDB.Model(&User{}).Find(&users, "role_id = ?", args.Role.ID); tx.Error != nil {
				l.WithError(tx.Error).Error("failed to get users of role")
				restapi.Error(c, 500, "database error")
				return
			}
			if len(users) > 0 {
				restapi.Error(c, 400, "cannot delete role with users")
				return
			}
			if tx := adminDB.Model(&args.Role).Delete(&args.Role); tx.Error != nil {
				l.WithError(tx.Error).Error("failed to delete role")
				restapi.Error(c, 500, "database error")
				return
			}
		} else {
			restapi.Error(c, 400, "role ID required")
			return
		}
	} else {
		if tx := adminDB.Save(&args.Role); tx.Error != nil {
			l.WithError(tx.Error).Error("failed to save role")
			restapi.Error(c, 500, "database error")
			return
		}
	}

	restapi.Success(c, nil)
	return

	// Make sure that roles with All Privs can't be modified unless this user has all privs
}
