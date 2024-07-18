package pkg

import (
	"fmt"
	"log"
	"runtime"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ErrorWrapper(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	_, file, line, _ := runtime.Caller(1)
	log.Printf("ERROR: [%s:%d] %v \n", file, line, err)
	return err
}

func StrToObjectID(id string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	return objectId
}

func StrToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println("Invalid id")
	}

	return num
}

func StrToInt(str string) int {
	int, err := strconv.Atoi(str)
	// Check for errors during conversion
	if err != nil {
		return 0

	}
	return int
}

func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

func CountTrues(boolSlice []bool) int {
	count := 0
	for _, val := range boolSlice {
		if val {
			count++
		}
	}
	return count
}

func CreateBoolList(size int, sublistSize int) [][]bool {
	list := make([][]bool, size)
	for i := range list {
		list[i] = make([]bool, sublistSize)
	}
	return list
}
