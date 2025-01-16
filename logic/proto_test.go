package logic

import "testing"

func TestSqlToProto(t *testing.T) {
	sql := ` 
		CREATE TABLE note_at_user (
			id bigint(20) DEFAULT NULL,
			note_id bigint(20) DEFAULT NULL,
			user_id bigint(20) DEFAULT NULL,
			created_at datetime DEFAULT CURRENT_TIMESTAMP,
			updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at datetime DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
	SqlToProto(sql, false)
}
