package dbwriter

import (
	"strings"
	"reflect"
	"github.com/json-iterator/go"
	"fmt"
	"time"
	orm "github.com/adtechpotok/go-orm"
)

func parseTagSetting(tags reflect.StructTag) map[string]string {
	setting := map[string]string{}
	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}

func makeSqlValue(element interface{}, serverId int) string {
	var result string
	v := reflect.ValueOf(element)
	s := reflect.Indirect(v)
	t := s.Type()
	num := t.NumField()
	for k := 0; k < num; k++ {
		var elementStringed string
		parsedTag := parseTagSetting(t.Field(k).Tag)
		name := s.Field(k).Type().Name()

		if val, ok := parsedTag["COLUMN"]; !ok || val == "COLUMN" {
			continue
		}

		switch name {
		case "bool":
			if s.Field(k).Bool() {
				elementStringed = "1"
			} else {
				elementStringed = "0"
			}

		case "int":
			elementStringed = fmt.Sprintf("%d", s.Field(k).Int())
		case "string":
			if _, ok := parsedTag["CLEARQUOTES"]; ok {
				elementStringed = clearQuotes(s.Field(k).String())
			} else {
				elementStringed = s.Field(k).String()
			}
			elementStringed = "'" + elementStringed + "'"
		case "Time":
			elementStringed = "'" + s.Field(k).Interface().(time.Time).Format("2006-01-02 15:04:05") + "'"
		case "float64":
			elementStringed = fmt.Sprintf("%f", s.Field(k).Float())
		default:
			if _, ok := parsedTag["FROMJSON"]; ok {
				json, _ := jsoniter.Marshal(s.Field(k).Interface())
				elementStringed = "'" + string(json) + "'"
			}
		}

		if elementStringed != "" {
			result += elementStringed + ","
		}
	}
	result += fmt.Sprintf("%d", serverId)

	return result

}

func makeSqlField(element interface{}) string {
	var result string
	v := reflect.ValueOf(element)
	s := reflect.Indirect(v)
	t := s.Type()
	num := t.NumField()
	for k := 0; k < num; k++ {
		parsedTag := parseTagSetting(t.Field(k).Tag)
		if val, ok := parsedTag["COLUMN"]; !ok || val == "COLUMN" {
			continue
		}

		result += parsedTag["COLUMN"] + ","
	}
	result += "server_id"
	return result
}

func makeInsertQuery(v orm.SchemaPotok, serverId int) string {
	res := fmt.Sprintf("INSERT INTO `%s`.`%s` (%s) VALUES (%s)", v.SchemaName(), v.TableName(), makeSqlField(v), makeSqlValue(v, serverId))
	if t, ok := v.(AfterSqlInterface); ok {
		res += " " + t.AfterSql()
	}
	return res
}

func clearQuotes(value string) string {
	replace := map[string]string{"\\": "\\\\", "'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}

type AfterSqlInterface interface {
	AfterSql() string
}
