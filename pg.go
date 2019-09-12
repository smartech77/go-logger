package logger

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	pgx "github.com/lib/pq"
)

// TODO: add more codes.
func ParsePGCode(code string, customMessage string) (outError *InformationConstruct) {
	switch code {
	case "42703": // Column not found error
		outError = BadRequest(nil, customMessage)
		outError.Hint = "Unknown column"
	case "42601": // bad syntax error
		outError = BadRequest(nil, customMessage)
		outError.Hint = "Your syntax might be off, review all your column and table references."
	case "22P02": // bad syntax error for UUID
		outError = BadRequest(nil, customMessage)
		outError.Hint = "You have an inalid UUID in your database transaction"
	case "42P01": // table not found
		outError = BadRequest(nil, customMessage)
		outError.Hint = "You are trying to save data to a table that does not exist, double check your table names"
	default:
		outError = BadRequest(nil, code+":"+customMessage)
	}
	return
}
func ParsePGError(er error) (outError *InformationConstruct) {
	ispgerror := false
	switch er.(type) {
	case *pgx.Error:
		ispgerror = true
	default:
		return GenericError(er)
	}
	if !ispgerror {
		return GenericError(er)
	}
	err := er.(*pgx.Error)
	switch err.Code {
	case "42703": // Column not found error
		outError = BadRequest(err, err.Routine)
		outError.Hint = err.Hint
		outError.Message = "This column does not appear to exist: " + strings.Split(err.Message, " ")[2]
	case "42601": // bad syntax error
		outError = BadRequest(err, err.Routine)
		outError.Hint = "Your syntax might be off, review all your column and table references."
		outError.Message = err.Message
	case "22P02": // bad syntax error for UUID
		outError = BadRequest(err, err.Routine)
		outError.Hint = "You have an inalid PRIMARY ID in your database transaction"
		outError.Message = err.Message
	case "42P01": // table not found
		outError = BadRequest(err, err.Routine)
		outError.Hint = "You are trying to save data to a table that does not exist, double check your table names"
		outError.Message = "This table does not appear to exist: " + strings.Split(err.Message, " ")[2]
	case "42701":
		outError = BadRequest(err, err.Routine)
		outError.Message = err.Message
	default:
		// this is  away to catch errors that are not supported.
		// so that they can be added.
		PrintObject(err)
	}
	return
}
func PrintObject(Object interface{}) {
	fields := reflect.TypeOf(Object).Elem()
	values := reflect.ValueOf(Object).Elem()
	num := fields.NumField()
	parseFields(num, fields, values)
}

func parseFields(num int, fields reflect.Type, values reflect.Value) {
	log.Println("!!!!!!!!!! UN-HANDLED POSTGRES ERROR !!!!!!!!!!")
	for i := 0; i < num; i++ {
		value := values.Field(i)
		field := fields.Field(i)

		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valueInt := strconv.FormatInt(value.Int(), 64)
			if valueInt != "" {
				fmt.Println(field.Name, valueInt)
			}
		case reflect.String:
			if value.String() != "" {
				fmt.Println(field.Name, value.String())
			}
		}
	}
	log.Println("!!!!!!!!!! UN-HANDLED POSTGRES ERROR !!!!!!!!!!")
}
