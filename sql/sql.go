package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "SELECT\n\t\t\tdashboard.id,\n\t\t\tdashboard.uid,\n\t\t\tdashboard.title,\n\t\t\tdashboard.slug,\n\t\t\tdashboard_tag.term,\n\t\t\tdashboard.is_folder,\n\t\t\tdashboard.folder_id,\n\t\t\tfolder.uid AS folder_uid,\n\t\t\tfolder.slug AS folder_slug,\n\t\t\tfolder.title AS folder_title \n\t\tFROM dashboard\n\t\tLEFT OUTER JOIN dashboard AS folder ON folder.id = dashboard.folder_id\n\t\tLEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id\n WHERE ((dashboard.uid IN (SELECT substr(scope, 16) FROM permission WHERE scope LIKE 'dashboards:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?) AND NOT dashboard.is_folder) OR ((dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  WHERE f1.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  WHERE f2.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  WHERE f3.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  WHERE f4.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  WHERE f5.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  WHERE f6.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  INNER JOIN folder f7 ON f6.parent_uid = f7.uid AND f6.org_id = f7.org_id  WHERE f7.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  INNER JOIN folder f7 ON f6.parent_uid = f7.uid AND f6.org_id = f7.org_id  INNER JOIN folder f8 ON f7.parent_uid = f8.uid AND f7.org_id = f8.org_id  WHERE f8.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  INNER JOIN folder f7 ON f6.parent_uid = f7.uid AND f6.org_id = f7.org_id  INNER JOIN folder f8 ON f7.parent_uid = f8.uid AND f7.org_id = f8.org_id  INNER JOIN folder f9 ON f8.parent_uid = f9.uid AND f8.org_id = f9.org_id  WHERE f9.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  INNER JOIN folder f7 ON f6.parent_uid = f7.uid AND f6.org_id = f7.org_id  INNER JOIN folder f8 ON f7.parent_uid = f8.uid AND f7.org_id = f8.org_id  INNER JOIN folder f9 ON f8.parent_uid = f9.uid AND f8.org_id = f9.org_id  INNER JOIN folder f10 ON f9.parent_uid = f10.uid AND f9.org_id = f10.org_id  WHERE f10.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT tr.role_id FROM team_role as tr\n\t\t\tWHERE tr.team_id IN(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)\n\t\t\tAND tr.org_id = ?\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?))) AND NOT dashboard.is_folder)) AND dashboard.org_id=? AND dashboard.title LIKE ? AND dashboard.is_folder = 0 ORDER BY dashboard.title ASC LIMIT 5000 OFFSET 0"
	args := []interface{}{
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
		int64(1),
		0,
		int64(1),
		int64(2),
		int64(3),
		int64(4),
		int64(5),
		int64(6),
		int64(7),
		int64(8),
		int64(9),
		int64(10),
		int64(11),
		int64(12),
		int64(13),
		int64(14),
		int64(15),
		int64(16),
		int64(17),
		int64(18),
		int64(19),
		int64(20),
		int64(21),
		int64(22),
		int64(23),
		int64(24),
		int64(25),
		int64(26),
		int64(27),
		int64(28),
		int64(29),
		int64(30),
		int64(31),
		int64(32),
		int64(33),
		int64(34),
		int64(35),
		int64(36),
		int64(37),
		int64(38),
		int64(39),
		int64(40),
		int64(41),
		int64(42),
		int64(43),
		int64(44),
		int64(45),
		int64(46),
		int64(47),
		int64(48),
		int64(49),
		int64(50),
		int64(1),
		"\"\"",
		int64(1),
		0,
		"\"dashboards:read\"",
		int64(1),
	}
	/*
		s := "SELECT\n\t\t\tdashboard.id,\n\t\t\tdashboard.uid,\n\t\t\tdashboard.title,\n\t\t\tdashboard.slug,\n\t\t\tdashboard_tag.term,\n\t\t\tdashboard.is_folder,\n\t\t\tdashboard.folder_id,\n\t\t\tfolder.uid AS folder_uid,\n\t\t\tfolder.slug AS folder_slug,\n\t\t\tfolder.title AS folder_title \n\t\tFROM dashboard\n\t\tLEFT OUTER JOIN dashboard AS folder ON folder.id = dashboard.folder_id\n\t\tLEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id\n WHERE (NOT dashboard.is_folder) AND dashboard.org_id=? AND dashboard_tag.term IN (?) AND dashboard.is_folder = false GROUP BY dashboard.title, dashboard.id HAVING COUNT(dashboard.id) >= ? ORDER BY dashboard.title ASC NULLS FIRST LIMIT 1000 OFFSET 0"
		args := []interface{}{
			int64(1)),
			"\"prod\""),
			1),
		}
	*/
	/*
		s := "\nSELECT title FROM dashboard WHERE ((dashboard.uid IN (SELECT substr(scope, 16) FROM permission WHERE scope LIKE 'dashboards:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?) AND NOT dashboard.is_folder) OR ((dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  WHERE f1.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  WHERE f2.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  WHERE f3.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  WHERE f4.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  WHERE f5.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  WHERE f6.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  INNER JOIN folder f7 ON f6.parent_uid = f7.uid AND f6.org_id = f7.org_id  WHERE f7.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  INNER JOIN folder f7 ON f6.parent_uid = f7.uid AND f6.org_id = f7.org_id  INNER JOIN folder f8 ON f7.parent_uid = f8.uid AND f7.org_id = f8.org_id  WHERE f8.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  INNER JOIN folder f7 ON f6.parent_uid = f7.uid AND f6.org_id = f7.org_id  INNER JOIN folder f8 ON f7.parent_uid = f8.uid AND f7.org_id = f8.org_id  INNER JOIN folder f9 ON f8.parent_uid = f9.uid AND f8.org_id = f9.org_id  WHERE f9.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) OR (dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  INNER JOIN folder f7 ON f6.parent_uid = f7.uid AND f6.org_id = f7.org_id  INNER JOIN folder f8 ON f7.parent_uid = f8.uid AND f7.org_id = f8.org_id  INNER JOIN folder f9 ON f8.parent_uid = f9.uid AND f8.org_id = f9.org_id  INNER JOIN folder f10 ON f9.parent_uid = f10.uid AND f9.org_id = f10.org_id  WHERE f10.uid IN (\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?)) AND NOT dashboard.is_folder))"
		args := []interface{}{
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
			"\"Viewer\""),
			int64(1)),
			0),
			"\"dashboards:read\""),
		}
	*/
	/*
		s = "\nSELECT DISTINCT\n\tle.name, le.id, le.org_id, le.folder_id, le.uid, le.kind, le.type, le.description, le.model, le.created, le.created_by, le.updated, le.updated_by, le.version\n\t, u1.login AS created_by_name\n\t, u1.email AS created_by_email\n\t, u2.login AS updated_by_name\n\t, u2.email AS updated_by_email\n\t, (SELECT COUNT(connection_id) FROM library_element_connection WHERE element_id = le.id AND kind=1) AS connected_dashboards, 'General' as folder_name , '' as folder_uid \nFROM library_element AS le\nLEFT JOIN `user` AS u1 ON le.created_by = u1.id\nLEFT JOIN `user` AS u2 ON le.updated_by = u2.id\n WHERE le.org_id=? AND le.uid=? AND le.folder_id=? UNION \nSELECT DISTINCT\n\tle.name, le.id, le.org_id, le.folder_id, le.uid, le.kind, le.type, le.description, le.model, le.created, le.created_by, le.updated, le.updated_by, le.version\n\t, u1.login AS created_by_name\n\t, u1.email AS created_by_email\n\t, u2.login AS updated_by_name\n\t, u2.email AS updated_by_email\n\t, (SELECT COUNT(connection_id) FROM library_element_connection WHERE element_id = le.id AND kind=1) AS connected_dashboards, dashboard.title as folder_name , dashboard.uid as folder_uid \nFROM library_element AS le\nLEFT JOIN `user` AS u1 ON le.created_by = u1.id\nLEFT JOIN `user` AS u2 ON le.updated_by = u2.id\n INNER JOIN dashboard AS dashboard on le.folder_id = dashboard.id AND le.folder_id <> 0 WHERE le.org_id=? AND le.uid=? AND (dashboard.uid IN (SELECT substr(scope, 16) FROM permission WHERE scope LIKE 'dashboards:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?) AND NOT dashboard.is_folder) OR (folder.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?) AND NOT dashboard.is_folder) OR (dashboard.uid IN (SELECT substr(scope, 13) FROM permission WHERE scope LIKE 'folders:uid:%' AND (\n\t\trole_id IN (\n\t\t\tSELECT ur.role_id\n\t\t\tFROM user_role AS ur\n\t\t\tWHERE ur.user_id = ?\n\t\t\tAND (ur.org_id = ? OR ur.org_id = ?)\n\t\t)\n\t\tOR\n\t\trole_id IN (\n\t\t\tSELECT br.role_id FROM builtin_role AS br\n\t\t\tWHERE br.role IN (?)\n\t\t\tAND (br.org_id = ? OR br.org_id = ?)\n\t\t)\n\t\t) AND action = ?) AND dashboard.is_folder) OR dashboard.id=0"
		args = []interface{}{
			int64(1),
			"\"b9891b26-f795-44bb-a283-b624979e7054\"",
			0,
			int64(1),
			"\"b9891b26-f795-44bb-a283-b624979e7054\"",
			int64(1),
			int64(1),
			0,
			"\"Viewer\"",
			int64(1),
			0,
			"\"dashboards:read\"",
			int64(1),
			int64(1),
			0,
			"\"Viewer\"",
			int64(1),
			0,
			"\"dashboards:read\"",
			int64(1),
			int64(1),
			0,
			"\"Viewer\"",
			int64(1),
			0,
			"\"folders:read\"",
		}
	*/
	for _, arg := range args {
		s = strings.Replace(s, "?", fmt.Sprintf("%v", arg), 1)
	}
	fmt.Println(s)
}
