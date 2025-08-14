package db

import (
	"database/sql"
	"fmt"
	"github.com/donnyhardyanto/dxlib/database2/database_type"
	"github.com/donnyhardyanto/dxlib/database2/db/raw"
	"github.com/donnyhardyanto/dxlib/database2/sqlchecker"
	"github.com/donnyhardyanto/dxlib/database2/utils/sql_expression"
	"github.com/donnyhardyanto/dxlib/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
)

// SQLPartInsertFieldNamesFieldValues generates the field names and values parts for an SQL INSERT statement
func SQLPartInsertFieldNamesFieldValues(insertKeyValues utils.JSON, driverName string) (fieldNames string, fieldValues string) {
	for k, v := range insertKeyValues {
		switch driverName {
		case "oracle":
			k = strings.ToUpper(k)
		}
		if fieldNames != "" {
			fieldNames = fieldNames + ","
		}
		fieldNames = fieldNames + k
		if fieldValues != "" {
			fieldValues = fieldValues + ","
		}
		switch v.(type) {
		case sql_expression.SQLExpression:
			fieldValues = fieldValues + v.(sql_expression.SQLExpression).String()
		default:
			fieldValues = fieldValues + ":" + k
		}
	}
	return fieldNames, fieldValues
}

// Insert performs a database insert with support for returning values across different database types
// Parameters:
//   - db: Database connection
//   - tableName: Target table name
//   - setFieldValues: Map of column names to values
//   - returningFieldNames: List of field names to return after insert
//
// Returns:
//   - returningFieldValues: Map of field names to their values after insert
//   - err: Error if any occurred
func Insert(db *sqlx.DB, tableName string, setFieldValues utils.JSON, returningFieldNames []string) (result sql.Result, returningFieldValues utils.JSON, err error) {
	// Basic input validation
	if db == nil {
		return nil, nil, errors.New("database connection is nil")
	}
	if tableName == "" {
		return nil, nil, errors.New("table name cannot be empty")
	}

	// Get the database driver name
	driverName := strings.ToLower(db.DriverName())
	dbType := database_type.StringToDXDatabaseType(driverName)

	// Validate table name explicitly
	if err := sqlchecker.CheckIdentifier(dbType, tableName); err != nil {
		return nil, nil, errors.Wrap(err, "invalid table name")
	}

	// Validate field names in setFieldValues
	for fieldName := range setFieldValues {
		if err := sqlchecker.CheckIdentifier(dbType, fieldName); err != nil {
			return nil, nil, errors.Wrapf(err, "invalid field name: %s", fieldName)
		}
	}

	// Validate returning field names
	for _, fieldName := range returningFieldNames {
		if err := sqlchecker.CheckIdentifier(dbType, fieldName); err != nil {
			return nil, nil, errors.Wrapf(err, "invalid returning field name: %s", fieldName)
		}
	}

	// Prepare field names and values for the INSERT statement
	fieldNames, fieldValues := SQLPartInsertFieldNamesFieldValues(setFieldValues, driverName)

	// Base INSERT statement
	baseSQL := strings.Join([]string{
		"INSERT INTO",
		tableName,
		fmt.Sprintf("(%s)", fieldNames),
		"VALUES",
		fmt.Sprintf("(%s)", fieldValues),
	}, " ")

	// Initialize return values
	returningFieldValues = utils.JSON{}

	// If no returning keys requested, simply execute the insert
	if returningFieldNames == nil || len(returningFieldNames) == 0 {
		result, err := raw.Exec(db, baseSQL, setFieldValues)
		return result, returningFieldValues, err
	}

	// Handle database-specific RETURNING clauses
	switch driverName {
	case "postgres", "mariadb":
		// Both PostgreSQL and MariaDB 10.5.0+ support RETURNING clause with the same syntax
		sqlStatement := fmt.Sprintf("%s RETURNING %s", baseSQL, strings.Join(returningFieldNames, ", "))
		_, rows, err := raw.QueryRows(db, nil, sqlStatement, setFieldValues)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error executing insert with RETURNING clause")
		}

		if len(rows) > 0 {
			returningFieldValues = rows[0]
		}

	case "sqlserver", "mssql":
		// SQL Server supports OUTPUT clause
		// Build OUTPUT clause
		var outputFields []string
		for _, key := range returningFieldNames {
			outputFields = append(outputFields, fmt.Sprintf("INSERTED.%s", key))
		}

		// Insert with OUTPUT clause
		sqlStatement := strings.Join([]string{
			"INSERT INTO",
			tableName,
			fmt.Sprintf("(%s)", fieldNames),
			"OUTPUT",
			strings.Join(outputFields, ", "),
			"VALUES",
			fmt.Sprintf("(%s)", fieldValues),
		}, " ")

		_, rows, err := raw.QueryRows(db, nil, sqlStatement, setFieldValues)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error executing insert with OUTPUT clause")
		}

		if len(rows) > 0 {
			returningFieldValues = rows[0]
		}

	case "oracle":
		// Oracle uses RETURNING INTO syntax
		// Prepare named arguments for Oracle
		namedArgs := make([]interface{}, 0, len(setFieldValues))
		for name, value := range setFieldValues {
			// Skip SQL expressions
			if _, ok := value.(sql_expression.SQLExpression); !ok {
				namedArgs = append(namedArgs, sql.Named(strings.ToUpper(name), value))
			}
		}

		// Build RETURNING INTO clause
		var returningFields []string
		var returningIntoFields []string

		for _, key := range returningFieldNames {
			returningFields = append(returningFields, key)
			returningIntoFields = append(returningIntoFields, fmt.Sprintf(":%s_out", key))

			// Add output parameters
			var outParam interface{}
			namedArgs = append(namedArgs, sql.Named(key+"_out", sql.Out{Dest: &outParam}))
		}

		sqlStatement := fmt.Sprintf("%s RETURNING %s INTO %s",
			baseSQL,
			strings.Join(returningFields, ", "),
			strings.Join(returningIntoFields, ", "))

		// Execute directly for Oracle with output parameters
		result, err = db.Exec(sqlStatement, namedArgs...)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error executing oracle insert with RETURNING INTO")
		}

		// Extract output parameters
		for _, arg := range namedArgs {
			namedArg, ok := arg.(sql.NamedArg)
			if !ok {
				continue
			}

			if strings.HasSuffix(namedArg.Name, "_out") {
				outArg, ok := namedArg.Value.(sql.Out)
				if !ok {
					continue
				}

				originalKey := strings.TrimSuffix(namedArg.Name, "_out")
				if outArg.Dest != nil {
					dest := outArg.Dest.(*interface{})
					returningFieldValues[originalKey] = *dest
				}
			}
		}

	case "mysql":
		// MySQL doesn't support RETURNING, so we need to do a separate query
		result, err := raw.Exec(db, baseSQL, setFieldValues)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error executing mysql insert")
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return nil, nil, errors.Wrap(err, "error getting last insert ID")
		}
		// Get the last insert ID and check if it's valid
		if lastInsertId <= 0 {
			// Some tables might not have auto-increment IDs, so this isn't always an error
			// Just return empty map if that's what the user wants
			return result, returningFieldValues, nil
		}

		// Use first common ID field name
		idFieldNames := []string{"id", "ID", "Id"}
		idField := ""

		// Find which ID field name was requested
		for _, fieldName := range idFieldNames {
			for _, key := range returningFieldNames {
				if strings.EqualFold(key, fieldName) {
					returningFieldValues[key] = result.LastInsertId
					idField = key
					break
				}
			}
			if idField != "" {
				break
			}
		}

		// If we didn't find a matching ID field but got an ID, use the first ID field name
		if idField == "" && len(idFieldNames) > 0 {
			idField = idFieldNames[0]
			returningFieldValues[idField] = result.LastInsertId
		}

		// Build fields to select, excluding the id field we already have
		var selectFields []string
		for _, key := range returningFieldNames {
			if !strings.EqualFold(key, idField) {
				selectFields = append(selectFields, key)
			}
		}

		if len(selectFields) > 0 {
			// Query for additional fields
			selectSQL := strings.Join([]string{
				"SELECT",
				strings.Join(selectFields, ", "),
				"FROM",
				tableName,
				"WHERE",
				fmt.Sprintf("%s = :%s", idField, idField),
			}, " ")

			selectArgs := utils.JSON{
				idField: result.LastInsertId,
			}

			_, rows, err := raw.QueryRows(db, nil, selectSQL, selectArgs)
			if err == nil && len(rows) > 0 {
				// Merge additional values
				for k, v := range rows[0] {
					returningFieldValues[k] = v
				}
			}
		}

	default:
		// Unsupported database type
		return nil, nil, errors.Errorf("unsupported database driver: %s", driverName)
	}

	return result, returningFieldValues, nil
}

func TxInsert(tx *sqlx.Tx, tableName string, setFieldValues utils.JSON, returningFieldNames []string) (result sql.Result, returningFieldValues utils.JSON, err error) {
	// Basic input validation
	if tx == nil {
		return nil, nil, errors.New("database transaction connection is nil")
	}
	if tableName == "" {
		return nil, nil, errors.New("table name cannot be empty")
	}

	// Get the database driver name
	driverName := strings.ToLower(tx.DriverName())
	dbType := database_type.StringToDXDatabaseType(driverName)

	// Validate table name explicitly
	if err := sqlchecker.CheckIdentifier(dbType, tableName); err != nil {
		return nil, nil, errors.Wrap(err, "invalid table name")
	}

	// Validate field names in setFieldValues
	for fieldName := range setFieldValues {
		if err := sqlchecker.CheckIdentifier(dbType, fieldName); err != nil {
			return nil, nil, errors.Wrapf(err, "invalid field name: %s", fieldName)
		}
	}

	// Validate returning field names
	for _, fieldName := range returningFieldNames {
		if err := sqlchecker.CheckIdentifier(dbType, fieldName); err != nil {
			return nil, nil, errors.Wrapf(err, "invalid returning field name: %s", fieldName)
		}
	}

	// Prepare field names and values for the INSERT statement
	fieldNames, fieldValues := SQLPartInsertFieldNamesFieldValues(setFieldValues, driverName)

	// Base INSERT statement
	baseSQL := strings.Join([]string{
		"INSERT INTO",
		tableName,
		fmt.Sprintf("(%s)", fieldNames),
		"VALUES",
		fmt.Sprintf("(%s)", fieldValues),
	}, " ")

	// Initialize return values
	returningFieldValues = utils.JSON{}

	// If no returning keys requested, simply execute the insert
	if returningFieldNames == nil || len(returningFieldNames) == 0 {
		result, err := raw.TxExec(tx, baseSQL, setFieldValues)
		return result, returningFieldValues, err
	}

	// Handle database-specific RETURNING clauses
	switch driverName {
	case "postgres", "mariadb":
		// Both PostgreSQL and MariaDB 10.5.0+ support RETURNING clause with the same syntax
		sqlStatement := fmt.Sprintf("%s RETURNING %s", baseSQL, strings.Join(returningFieldNames, ", "))
		_, rows, err := raw.TxQueryRows(tx, nil, sqlStatement, setFieldValues)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error executing insert with RETURNING clause")
		}

		if len(rows) > 0 {
			returningFieldValues = rows[0]
		}

	case "sqlserver", "mssql":
		// SQL Server supports OUTPUT clause
		// Build OUTPUT clause
		var outputFields []string
		for _, key := range returningFieldNames {
			outputFields = append(outputFields, fmt.Sprintf("INSERTED.%s", key))
		}

		// Insert with OUTPUT clause
		sqlStatement := strings.Join([]string{
			"INSERT INTO",
			tableName,
			fmt.Sprintf("(%s)", fieldNames),
			"OUTPUT",
			strings.Join(outputFields, ", "),
			"VALUES",
			fmt.Sprintf("(%s)", fieldValues),
		}, " ")

		_, rows, err := raw.TxQueryRows(tx, nil, sqlStatement, setFieldValues)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error executing insert with OUTPUT clause")
		}

		if len(rows) > 0 {
			returningFieldValues = rows[0]
		}

	case "oracle":
		// Oracle uses RETURNING INTO syntax
		// Prepare named arguments for Oracle
		namedArgs := make([]interface{}, 0, len(setFieldValues))
		for name, value := range setFieldValues {
			// Skip SQL expressions
			if _, ok := value.(sql_expression.SQLExpression); !ok {
				namedArgs = append(namedArgs, sql.Named(strings.ToUpper(name), value))
			}
		}

		// Build RETURNING INTO clause
		var returningFields []string
		var returningIntoFields []string

		for _, key := range returningFieldNames {
			returningFields = append(returningFields, key)
			returningIntoFields = append(returningIntoFields, fmt.Sprintf(":%s_out", key))

			// Add output parameters
			var outParam interface{}
			namedArgs = append(namedArgs, sql.Named(key+"_out", sql.Out{Dest: &outParam}))
		}

		sqlStatement := fmt.Sprintf("%s RETURNING %s INTO %s",
			baseSQL,
			strings.Join(returningFields, ", "),
			strings.Join(returningIntoFields, ", "))

		// Execute directly for Oracle with output parameters
		_, err = tx.Exec(sqlStatement, namedArgs...)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error executing oracle insert with RETURNING INTO")
		}

		// Extract output parameters
		for _, arg := range namedArgs {
			namedArg, ok := arg.(sql.NamedArg)
			if !ok {
				continue
			}

			if strings.HasSuffix(namedArg.Name, "_out") {
				outArg, ok := namedArg.Value.(sql.Out)
				if !ok {
					continue
				}

				originalKey := strings.TrimSuffix(namedArg.Name, "_out")
				if outArg.Dest != nil {
					dest := outArg.Dest.(*interface{})
					returningFieldValues[originalKey] = *dest
				}
			}
		}

	case "mysql":
		// MySQL doesn't support RETURNING, so we need to do a separate query
		result, err := raw.TxExec(tx, baseSQL, setFieldValues)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error executing mysql insert")
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return nil, nil, errors.Wrap(err, "error getting last insert ID")
		}
		// Get the last insert ID and check if it's valid
		if lastInsertId <= 0 {
			// Some tables might not have auto-increment IDs, so this isn't always an error
			// Just return empty map if that's what the user wants
			return result, returningFieldValues, nil
		}

		// Use first common ID field name
		idFieldNames := []string{"id", "ID", "Id"}
		idField := ""

		// Find which ID field name was requested
		for _, fieldName := range idFieldNames {
			for _, key := range returningFieldNames {
				if strings.EqualFold(key, fieldName) {
					returningFieldValues[key] = result.LastInsertId
					idField = key
					break
				}
			}
			if idField != "" {
				break
			}
		}

		// If we didn't find a matching ID field but got an ID, use the first ID field name
		if idField == "" && len(idFieldNames) > 0 {
			idField = idFieldNames[0]
			returningFieldValues[idField] = result.LastInsertId
		}

		// Build fields to select, excluding the id field we already have
		var selectFields []string
		for _, key := range returningFieldNames {
			if !strings.EqualFold(key, idField) {
				selectFields = append(selectFields, key)
			}
		}

		if len(selectFields) > 0 {
			// Query for additional fields
			selectSQL := strings.Join([]string{
				"SELECT",
				strings.Join(selectFields, ", "),
				"FROM",
				tableName,
				"WHERE",
				fmt.Sprintf("%s = :%s", idField, idField),
			}, " ")

			selectArgs := utils.JSON{
				idField: result.LastInsertId,
			}

			_, rows, err := raw.TxQueryRows(tx, nil, selectSQL, selectArgs)
			if err == nil && len(rows) > 0 {
				// Merge additional values
				for k, v := range rows[0] {
					returningFieldValues[k] = v
				}
			}
		}

	default:
		// Unsupported database type
		return nil, nil, errors.Errorf("unsupported database driver: %s", driverName)
	}

	return result, returningFieldValues, nil
}
