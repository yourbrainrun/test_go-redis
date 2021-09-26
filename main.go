package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"io/ioutil"
	rand2 "math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"test_go-redis/drivers"
	"test_go-redis/helpers"
	"time"
	"unsafe"
)

func init() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
	}
}

func main() {
	//zSetAdd()
	//return

	//list()
	//AddMember()
	//return
	//randClean()
	//return

	scan()
	return

	for i := 0; i < 5; i++ {
		s, err := rand1()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(s)
	}
	return

	s2, err := rand1()
	if err != nil {
		return
	}
	fmt.Println(s2)
	//backData()
	return

	s, err := rand()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(s)
	return

	for i := 10; i > 0; i-- {
		//go test()
	}
	time.Sleep(1 * time.Second)

	url, err := getBootApiUrl()
	if err != nil {
		return
	}
	fmt.Println(fmt.Sprintf("%s://%s", "http", url))
}

func backData() {
	key := "sort:list:push"
	member := "rtx:124.0.9.999"
	ctx1 := context.Background()
	rdb := drivers.GetRedisClient()

	_, err := rdb.ZIncrBy(ctx1, key, 1, member).Result()
	fmt.Println(err)
}

func AddMember() {
	ctx := context.Background()
	rdb := drivers.GetRedisClient()
	key := "sort:list:push"

	member := "rtx:11:11:11:11"

	result, err := rdb.ZIncrBy(ctx, key, 1, member).Result()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
}

func zSetAdd() {
	ctx := context.Background()
	rdb := drivers.GetRedisClient()
	key := "sort:list:push"
	args := make([]*redis.Z, 0)
	args = append(args, &redis.Z{
		Score:  10,
		Member: "rtx:124.0.9.24",
	})

	for i := 0; i < 1000; i++ {
		args = append(args, &redis.Z{
			Score:  20,
			Member: "rtx:124.0.9." + fmt.Sprintf("%d", i),
		})
	}

	rdb.ZAdd(ctx, key, args...)
}

func scan() {
	ctx := context.Background()
	rdb := drivers.GetRedisClient()

	var cursor uint64 = 0
	var count int64 = 1

	for {
		result, retCursor, err := rdb.Scan(ctx, cursor, "rtx:*", count).Result()
		cursor = retCursor
		if err != nil {
			return
		}
		for index, value := range result {
			tempSlice := strings.Split(value, ":")
			fmt.Println(index, "=", value, tempSlice[1])
			countT, err := rdb.Exists(ctx, value).Result()
			if err != nil {
				return
			}
			fmt.Println(countT)

			countT, err = rdb.Exists(ctx, "noon").Result()
			if err != nil {
				return
			}
			fmt.Println(countT)
		}

		if cursor <= 0 {
			break
		}
	}
}

func randClean() (string, error) {
	//ZINCRBY key increment member

	ctx := context.Background()
	rdb := drivers.GetRedisClient()
	key := "sort:list:push"

	keys := make([]string, 0)
	scores := make([]int64, 0)

	var result []string
	var err error
	var step int64
	var cursor uint64 = 0
	var floatPointer int64 = 0
	step = 500
	for {
		result, cursor, err = rdb.ZScan(ctx, key, cursor, "", step).Result()
		if err != nil {
			return "", err
		}

		for index, value := range result {
			if index%2 == 0 {
				scriptOne(key, value)
				keys = append(keys, value)
			} else {
				temp, err := strconv.ParseInt(value, 10, 64)
				if temp <= 0 {
					scores = append(scores, 0)
					continue
				}

				if err != nil {
					return "", err
				}
				floatPointer = floatPointer + temp
				scores = append(scores, floatPointer)
			}
		}
		if cursor <= 0 || len(result) == 0 {
			break
		}
	}
	if floatPointer <= 0 {
		return "", errors.New("没有符合条件 score")
	}
	return "", nil
}

func scriptOne(key string, member string) {
	var luaScript = redis.NewScript(`
local member = KEYS[1]
local key = KEYS[2]
local flag
redis.log(redis.LOG_NOTICE,"member=",member)
flag = redis.call('EXISTS',member)
if flag < 1 then
    redis.call('ZREM',key,member)
end
`)
	rdb := drivers.GetRedisClient()

	luaScript.Run(context.Background(), rdb, []string{member, key}, 2)
}

func script(key string, member string) {
	f, err := os.Open("./scripts/clean_keys.lua")
	if err != nil {
		return
	}
	myBytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err.Error())
	}
	scriptString := UnsafeBytesToString(myBytes)
	var luaScript = redis.NewScript(scriptString)
	rdb := drivers.GetRedisClient()

	luaScript.Run(context.Background(), rdb, []string{member, key}, 2)
}

func UnsafeBytesToString(bytes []byte) string {
	hdr := &reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&bytes[0])),
		Len:  len(bytes),
	}
	return *(*string)(unsafe.Pointer(hdr))
}

func rand() (string, error) {
	ctx := context.Background()
	rdb := drivers.GetRedisClient()
	key := "sort:list:push"

	keys := make([]string, 0)
	scores := make([]int64, 0)

	var err error
	var floatPointer int64 = 0
	result, err := rdb.ZRange(ctx, key, 0, -1).Result()
	fmt.Println(result)
	if err != nil {
		fmt.Println("cursor error", err.Error())
		return "", err
	}

	for index, value := range result {
		if index%2 == 0 {
			keys = append(keys, value)
		} else {
			temp, err := strconv.ParseInt(value, 10, 64)
			if temp <= 0 {
				scores = append(scores, 0)
				continue
			}

			if err != nil {
				return "", err
			}
			floatPointer = floatPointer + temp
			scores = append(scores, floatPointer)
		}
	}

	if floatPointer <= 0 {
		return "", errors.New("没有符合条件 score")
	}

	var s2 = rand2.NewSource(time.Now().UnixNano())
	r2 := rand2.New(s2)
	randValue := r2.Int63n(floatPointer)

	var member string
	for index, value := range scores {
		if value > randValue {
			member = keys[index]
			break
		}
	}

	rdb.ZIncrBy(ctx, key, -1, member)
	return member, err
}

func rand1() (string, error) {
	//ZINCRBY key increment member

	ctx := context.Background()
	rdb := drivers.GetRedisClient()
	key := "sort:list:push"

	keys := make([]string, 0)
	scores := make([]int64, 0)

	var err error
	var step int64
	var cursor uint64 = 0
	var floatPointer int64 = 0
	step = 500
	for {
		result, retCursor, err := rdb.ZScan(ctx, key, cursor, "", step).Result()
		cursor = retCursor
		if err != nil {
			return "", err
		}

		for index, value := range result {
			if index%2 == 0 {
				keys = append(keys, value)
			} else {
				temp, err := strconv.ParseInt(value, 10, 64)
				if temp <= 0 {
					scores = append(scores, 0)
					continue
				}

				if err != nil {
					return "", err
				}
				floatPointer = floatPointer + temp
				scores = append(scores, floatPointer)
			}
		}
		if cursor <= 0 || len(result) == 0 {
			break
		}
	}
	if floatPointer <= 0 {
		return "", errors.New("没有符合条件 score")
	}

	s2 := rand2.NewSource(time.Now().UnixNano())
	r2 := rand2.New(s2)
	randValue := r2.Int63n(floatPointer)

	var member string
	for index, value := range scores {
		if value > randValue {
			member = keys[index]
			break
		}
	}

	rdb.ZIncrBy(ctx, key, -1, member)
	return member, err
}

func list() {
	ctx := context.Background()
	rdb := drivers.GetRedisClient()
	key := "list:push"
	rdb.LPush(ctx, key, "rtx:124.0.9.23")
	rdb.LPush(ctx, key, "rtx:124.0.9.24")
	rdb.Set(ctx, "rtx:124.0.9.23", "123.233.11.22", -1)
	rdb.Set(ctx, "rtx:124.0.9.24", "rtx:124.0.9.24", -1)
}

func getBootApiUrl() (string, error) {
	key := "list:push"
	ctx := context.Background()
	rdb := drivers.GetRedisClient()
	for {
		pullerKey, err := rdb.LPop(ctx, key).Result()
		if err != nil {
			return "", err
		}

		pullerIp, err := rdb.Get(ctx, pullerKey).Result()
		if err == nil {
			return pullerIp, nil
		}
	}
}

func test() {
	key := "eqweqweqweqweqwe1234"
	lock := &helpers.RedisLock{ReleaseDuration: 10 * time.Second}
	if lock.Lock(key) == false {
		fmt.Println("already")
		return
	}
	fmt.Println("locked")
	time.Sleep(10 * time.Second)
	defer lock.Unlock(key)
}
