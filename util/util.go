// gen by iyfiysi at 2020-08-02 22:41:29.9166513 +0800 CST m=+13.827069201
package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/spf13/viper"
	"math/rand"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var gSnowflake *snowflake.Node

const CRLF = "\n"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Now() uint {
	t := time.Now().Unix()
	return uint(t)
}

func NowMs() uint64 {
	t := time.Now().UnixNano()
	return uint64(t / int64(time.Millisecond))
}

func NowHourTs() uint {
	now := time.Now()
	hourTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	return uint(hourTime.Unix())
}

func DayTs() uint {
	t := time.Now() //.UTC()
	d := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, //hour
		0, //min
		0, //sec
		0, //ms
		//time.UTC)
		time.Local)

	return uint(d.Unix())
}

/*
hourTimeStr=09:10
hourTs= 1571706600
*/
func ParseHourTimeStr2Ts(hourTimeStr string) (hourTs uint, err error) {
	tmp := strings.Split(hourTimeStr, ":")
	if len(tmp) != 2 {
		err = fmt.Errorf("hourTimeStr=%s format err", hourTimeStr)
		return
	}
	hourPart := Str2Num(tmp[0])
	MinutePart := Str2Num(tmp[1])
	now := time.Now()
	hourTime := time.Date(now.Year(), now.Month(), now.Day(), int(hourPart), int(MinutePart), 0, 0, now.Location())
	hourTs = uint(hourTime.Unix())
	return
}

/*
hourTimeStr=2020-03-20
hourTs= 1571706600
*/
func ParseDateStr2Ts(dateStr string) (ts uint, err error) {
	ts = 0
	tmp := strings.Split(dateStr, "-")
	if len(tmp) != 3 {
		err = fmt.Errorf("ParseDateStr2Ts=%s format err", dateStr)
		return
	}
	yearPart := Str2Num(tmp[0])
	monthPart := Str2Num(tmp[1])
	dayPart := Str2Num(tmp[2])
	now := time.Now()
	tsTime := time.Date(int(yearPart), time.Month(monthPart), int(dayPart), 0, 0, 0, 0, now.Location())
	ts = uint(tsTime.Unix())
	return
}

//return: [0,n)
func RandN(n int) int {
	return rand.Intn(n)
}

func GetUpperCase() string {
	return "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
}

func GetLowwerCase() string {
	return "abcdefghijklmnopqrstuvwxyz"
}

func GetStrDigit() string {
	return "0123456789"
}

func RandInt() int {
	return rand.Int()
}

func RandIntn(start, end int) int {
	if end == start {
		return start
	}
	if end < start {
		start, end = end, start
	}
	range_ := end - start + 1
	r := rand.Intn(range_)
	return r + start
}

func RandStr(l int) string {
	str := GetUpperCase() + GetLowwerCase() + GetStrDigit()
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < l; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func CopyStuct(dst interface{}, src interface{}) (err error) {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		err = errors.New("dst isn't a pointer to struct")
		return
	}
	dstElem := dstValue.Elem()
	if dstElem.Kind() != reflect.Struct {
		err = errors.New("pointer doesn't point to struct")
		return
	}

	srcValue := reflect.ValueOf(src)
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Struct {
		err = errors.New("src isn't struct")
		return
	}

	for i := 0; i < srcType.NumField(); i++ {
		sf := srcType.Field(i)
		sv := srcValue.FieldByName(sf.Name)
		// make sure the value which in dst is valid and can set
		if dv := dstElem.FieldByName(sf.Name); dv.IsValid() && dv.CanSet() {
			dv.Set(sv)
		}
	}
	return
}

//string to float64
//"123.435"=>123.435
//"123.435hahahha"=>123.435
//"123.435hahahh567.123"=>123.435
//"hahah"=>err
func Str2Float(str string) (num float64, err error) {
	var valid = regexp.MustCompile("(^[0-9][0-9.]*)")
	result := valid.FindAllStringSubmatch(str, -1)

	if len(result) <= 0 {
		err = fmt.Errorf("str=" + str + " not valid")
		return
	}
	if len(result[0]) <= 1 {
		err = fmt.Errorf("str=" + str + " not valid cuz substr")
		return
	}

	str = result[0][1]
	num, err = strconv.ParseFloat(str, 0)
	return
}

//string to int
func Str2Num(str string) (num int64) {
	num = 0
	tmp, _ := strconv.ParseInt(str, 10, 0)
	num = int64(tmp)
	return
}

func Num2Str(num interface{}) (str string) {
	str = ""
	switch num.(type) {
	case int:
		str = strconv.FormatInt(int64(num.(int)), 10)
	case int8:
		str = strconv.FormatInt(int64(num.(int8)), 10)
	case int16:
		str = strconv.FormatInt(int64(num.(int16)), 10)
	case int32:
		str = strconv.FormatInt(int64(num.(int32)), 10)
	case int64:
		str = strconv.FormatInt(int64(num.(int64)), 10)
	case uint:
		str = strconv.FormatInt(int64(num.(uint)), 10)
	case uint8:
		str = strconv.FormatInt(int64(num.(uint8)), 10)
	case uint16:
		str = strconv.FormatInt(int64(num.(uint16)), 10)
	case uint32:
		str = strconv.FormatInt(int64(num.(uint32)), 10)
	case uint64:
		str = strconv.FormatInt(int64(num.(uint64)), 10)
	case float32:
		str = strconv.FormatFloat(float64(num.(float32)), 'f', -1, 32)
	case float64:
		str = strconv.FormatFloat(float64(num.(float64)), 'f', -1, 64)
	default:
		panic("type err")
	}
	return
}

func Struct2Str(v interface{}) (str string) {
	str = ""
	b, err := json.Marshal(v)
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return
	}
	str = out.String()
	return
}

func JsonStr2Map(jsonStr string) (m map[string]interface{}, err error) {
	m = make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return
	}
	return
}

func UUID() (id int64) {
	if gSnowflake == nil {
		nodeId := viper.GetInt64("snowflakeNodeID")
		gSnowflake, _ = snowflake.NewNode(nodeId)
	}
	return gSnowflake.Generate().Int64()
}

func Elasped(start time.Time) int {
	elaspedMs := int(time.Since(start).Nanoseconds() / 1000 / 1000)
	return elaspedMs
}

func GetKeysByTag(structVal interface{}, tagName string) (keys []string) {
	keys = make([]string, 0)
	t := reflect.TypeOf(structVal).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagKey := field.Tag.Get(tagName)
		keys = append(keys, tagKey)
	}
	return
}

type helperStringIntSortStructItem struct {
	key string
	val interface{}
}

type helperStringIntSortStructSlice []helperStringIntSortStructItem

func (h helperStringIntSortStructSlice) Len() int {
	return len(h)
}

func (h helperStringIntSortStructSlice) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h helperStringIntSortStructSlice) Less(i, j int) bool { //重写Less()方法,从大到小排序
	switch h[i].val.(type) {
	case int:
		return h[j].val.(int) < h[i].val.(int)
	case int32:
		return h[j].val.(int32) < h[i].val.(int32)
	case int64:
		return h[j].val.(int64) < h[i].val.(int64)
	case uint64:
		return h[j].val.(uint64) < h[i].val.(uint64)
	case float32:
		return h[j].val.(float32) < h[i].val.(float32)
	case float64:
		return h[j].val.(float64) < h[i].val.(float64)
	}
	return true
}

func SortStringIntMapDesc(data map[string]interface{}) (result []string) {
	result = make([]string, 0)
	tmp := helperStringIntSortStructSlice{}
	for key, value := range data {
		tmp = append(tmp, helperStringIntSortStructItem{key, value})
	}
	sort.Stable(tmp)
	for _, v := range tmp {
		result = append(result, v.key)
	}
	return
}

func Shuffle(data []string) {
	for i := range data {
		j := RandN(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

func GetBetweenStr(str, start, end string) string {
	s := strings.Index(str, start)
	if s < 0 {
		return ""
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e < 0 {
		return ""
	}
	return str[s : s+e]
}

//将str的头尾的标点符号去掉
//包含中文标点符号，英文标点符号（空格不是标点符号）
func TrimPunct(str string) (trimStr string) {
	puncts := []*unicode.RangeTable{unicode.Punct}
	runes := []rune(str)

	//去头
	headIdx := 0
	for i, r := range runes {
		if !unicode.IsOneOf(puncts, r) {
			headIdx = i
			break
		}
	}

	//去尾
	tailIdx := len(runes)
	for i := len(runes) - 1; i >= 0; i-- {
		if !unicode.IsOneOf(puncts, runes[i]) {
			tailIdx = i + 1
			break
		}
	}
	//s = s[low : high : max]
	//range=(low,high]
	trimStr = string(runes[headIdx:tailIdx])
	return
}

//分割字符串：使用sepList来分割
//举例：
// str=hello world##hello world2#.
// sepList=[]string{"##","#."}
//ret:
//	[]string{"hello world","hello world2"}
//若是str不存在sepList，则返回为空list（[]string{}）
func SplitBySeps(str string, seps []string) (splitList []string) {
	splitList = make([]string, 0)

	splitStrs := []string{str}
	for _, sep := range seps {
		tmps := make([]string, 0)
		for _, s := range splitStrs {
			if s == "" {
				continue
			}
			tmp := strings.Split(s, sep)
			for _, t := range tmp {
				if t == "" {
					continue
				}
				tmps = append(tmps, t)
			}
		}
		splitStrs = tmps
	}
	splitList = splitStrs
	return
}

//和SplitBySeps功能一致，不过，其会将sep附加在分割后的句子的前面
//举例："say:hello worldsay:ggg"
func SplitBeforeBySeps(str string, seps []string) (splitList []string) {
	splitList = make([]string, 0)

	splitStrs := []string{str}
	for _, sep := range seps {
		tmps := make([]string, 0)
		for _, s := range splitStrs {
			if s == "" {
				continue
			}

			tmp := strings.Split(s, sep)
			if len(tmp) == 1 {
				tmps = append(tmps, tmp[0])
				continue
			}
			first := true
			for _, t := range tmp {
				if t == "" {
					first = false
					continue
				}
				tt := t

				//第一个
				if first {
					tmps = append(tmps, tt)
					first = false
					continue
				}

				tmps = append(tmps, sep+tt)
			}
		}
		splitStrs = tmps
	}
	splitList = splitStrs
	return
}

//将Unicode编码转换为string
//fmt.Println(U2S("\u54c8\u54c8"))==>"哈哈"
func U2S(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, "\\u", "", -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func KSort(srcMap map[string]string) (sortStr string) {
	ks := make([]string, 0)
	for k, _ := range srcMap {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for idx, k := range ks {
		ks[idx] += "=" + srcMap[k]
	}

	sortStr = strings.Join(ks, "&")
	return
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func UniqueStringList(src []string) (dst []string) {
	dst = make([]string, 0)

	tmp := make(map[string]bool)
	for _, item := range src {
		if _, ok := tmp[item]; !ok {
			tmp[item] = true
			dst = append(dst, item)
		}
	}
	return
}

func DeferWhenCoreDump() {
	if err := recover(); err != nil {
		//打印栈信息
		stack := debug.Stack()
		MainLogger.Error(string(stack))
		panic(err)
	}
}
