package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	chars         = "GFW5XC3U9ZM6NTB7YR2HS8DVEJ4KQPAL" //随机字符串
	decimal int64 = 32                                 //进制数
	lenth   int64 = 8                                  //邀请码长度
	salt    int64 = 1234561                            //随机数据（盐值）
	prime1  int64 = 3                                  //prime1 与 chars 的长度 L 互质，可保证 (id * prime1) % L 在 [0,L] 上均匀分布
	prime2  int64 = 15                                 //prime2 与 lenth 互质，可保证 (index * prime2) % lenth 在 [0，lenth] 上均匀分布
)

func Encode(id int64) string {
	// 补位
	id = id*prime1 + salt
	// 将 id 转换成32进制的值
	var b [8]int64
	b[0] = id
	var i int64
	for i = 0; i < lenth-1; i++ {
		b[i+1] = b[i] / decimal
		// 按位扩散
		b[i] = (b[i] + int64(i)*b[0]) % decimal
	}
	b[7] = (b[0] + b[1] + b[2] + b[3] + b[4] + b[5] + b[6]) * prime1 % decimal
	// 混淆
	var c [8]int64
	for i = 0; i < lenth; i++ {
		c[i] = b[int64(i)*prime2%lenth]
	}
	// 取出的索引转换为字符
	var code string
	for _, value := range c {
		code += string(chars[value])
	}
	return code
}

func Decode(code string) int64 {
	// 将字符还原成对应数字
	var a [8]int64
	var i int64
	for i = 0; i < lenth; i++ {
		c := string(code[i])
		index := strings.Index(chars, c)
		if index == -1 {
			// 异常字符串
			return 0
		}
		a[i*prime2%lenth] = int64(index)
	}
	var b [8]int64
	for i = lenth - 2; i >= 0; i-- {
		b[i] = (a[i] - a[0]*i + decimal*i) % decimal
	}

	var result int64
	var num int64
	for i = lenth - 2; i >= 0; i-- {
		result += b[i]
		num = 1
		if i > 0 {
			num = decimal
		}
		result = result * num
	}

	return (result - salt) / prime1
}

//获取6位随机数
func GetRandNum() string {
	var code string
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for index := 0; index < 6; index++ {
		num := r.Intn(6)
		code += strconv.Itoa(num)
	}
	return code
}

func main() {
	randNum, _ := strconv.Atoi(GetRandNum())
	fmt.Println(randNum)
	code := Encode(int64(randNum))
	fmt.Println(code)

	id := Decode(code)
	fmt.Println(id)
}
