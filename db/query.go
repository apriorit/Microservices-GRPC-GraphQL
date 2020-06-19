package db

import "fmt"

func ListRecords(collectionName string) string {
	const listRecordsQuery = `
	FOR record IN %s
		RETURN record`
	return fmt.Sprintf(listRecordsQuery, collectionName)
}
