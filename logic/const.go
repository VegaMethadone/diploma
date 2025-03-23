package logic

import (
	"labyrinth/config"
	"labyrinth/database/postgres"
)

var ps = postgres.NewPostgres(config.Conf)

var InjectionKeywords = []string{
	"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER", "TRUNCATE",
	"EXEC", "EXECUTE", "UNION", "JOIN", "GRANT", "REVOKE", "SHOW", "DATABASE",
	"TABLE", "OR", "AND", "WHERE", "FROM", "INTO", "VALUES", "SET", "HAVING",
	"GROUP BY", "ORDER BY", "LIMIT", "OFFSET", "UNION ALL", "UNION SELECT",
	"WAITFOR", "DELAY", "xp_cmdshell", "LOAD_FILE", "INTO OUTFILE", "INTO DUMPFILE",
	"CONCAT", "--", "#", "/*", "*/", ";", "' OR '1'='1", "'", "\"", "\\",
}
