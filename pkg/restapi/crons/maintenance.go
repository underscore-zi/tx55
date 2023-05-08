package crons

import (
	"gorm.io/gorm"
)

func ClearOldSessions(db *gorm.DB) {
	db.Exec(`DELETE FROM sessions WHERE created_at < DATE_SUB(NOW(), INTERVAL 14 DAY) 
                 AND id NOT IN (
                     SELECT id FROM (
                         SELECT id, user_id FROM sessions GROUP BY user_id HAVING MAX(created_at)
                     ) AS newest_sessions
                 )`)
}
