package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	width  = 101
	height = 26
)

type Universe [][]bool

//初始化二维网格
func NewUniverse() Universe {
	u := make(Universe, height)

	for i := range u {
		u[i] = make([]bool, width)
	}
	return u
}

//随机激活网格内25%的细胞
func (u Universe) Seed() {
	for i := 0; i < (width * height / 4); i++ {
		rand.Seed(time.Now().UnixNano())
		u.Set(rand.Intn(height), rand.Intn(width), true)
	}
}

//设置细胞状态
func (u Universe) Set(h, w int, b bool) {
	u[h][w] = b
}

//判断细胞是否存活
func (u Universe) Alive(h, w int) bool {
	h = (h + height) % height
	w = (w + width) % width
	return u[h][w]
}

//统计临近细胞存活数量
func (u Universe) Neighbors(h, w int) int {
	n := 0
	for hh := -1; hh <= 1; hh++ {
		for ww := -1; ww <= 1; ww++ {
			if !(hh == 0 && ww == 0) && u.Alive(h+hh, w+ww) {
				n++
			}
		}
	}
	return n
}

//返回执行细胞在下一代中的状态
func (u Universe) Next(h, w int) bool {
	n := u.Neighbors(h, w)

	//resBool := true
	//
	//switch n {
	//case 3:
	//	resBool = true
	//case 2:
	//	resBool = u.Alive(h, w)
	//default:
	//	resBool = false
	//}
	//
	//return resBool
	return n == 3 || n == 2 && u.Alive(h, w)
}

//将世界a的状态更新至下一代，并存储在b
func Step(a, b Universe) {
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			b.Set(h, w, a.Next(h, w))
		}
	}
}

//以字符串的形式返回整个世界
func (u Universe) String() string {
	var b byte

	buf := make([]byte, 0, (width+1)*height)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			b = ' '
			if u[h][w] {
				b = '*'
			}
			buf = append(buf, b)
		}
		buf = append(buf, '\n')
	}

	return string(buf)
}

//清空屏幕，显示网络
func (u Universe) Show() {
	fmt.Print("\033[H", u.String())
}

func main() {

	a, b := NewUniverse(), NewUniverse()
	a.Seed()
	fmt.Print("\033[H")

	for {
		Step(a, b)
		a.Show()
		time.Sleep(time.Second / 30)
		a, b = b, a
	}
}

