package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

/*
 * Field 表示二维细胞
 * 切片s 表示细胞，切片
 * w，h 宽，高。
 */
type Field struct {
	s    [][]bool
	w, h int
}

/**
 * 初始化细胞阵列
 */
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s: s, w: w, h: h}
}

/*
 * Set 设置 Field 值
 * 这里 w 和 h 对应的坐标应该是 y,x，所以 f.s[y][x]，而不是f.s[x][y]
 */

func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
}

/**
 * 判断细胞是否存活
 * 如果 x 或者 y 坐标在边界之外，则替换到边界之内，在边界内，不影响
 * 比如 -1 等效为 -1+w 再取模
 *
 */
func (f *Field) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

/*
 * 在下一步返回指定位置的状态
 * 1，计算活的相邻细胞总数
 * 生命游戏中，对于任意细胞，规则如下：
 * 每个细胞有两种状态 - 存活或死亡，每个细胞与以自身为中心的周围八格细胞产生互动
 * 当前细胞为存活状态时，当周围的存活细胞低于2个时（不包含2个），该细胞变成死亡状态。（模拟生命数量稀少）
 * 当前细胞为存活状态时，当周围有2个或3个存活细胞时，该细胞保持原样。
 * 当前细胞为存活状态时，当周围有超过3个存活细胞时，该细胞变成死亡状态。（模拟生命数量过多）
 * 当前细胞为死亡状态时，当周围有3个存活细胞时，该细胞变成存活状态。（模拟繁殖）
 */
func (f *Field) Next(x, y int) bool {
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}

	sbool := true

	switch alive {
		case 3: sbool = true
		case 2: sbool = f.s[y][x]
		default: sbool = false
	}

	return sbool
	//return alive == 3 || alive == 2 && f.Alive(x, y)
}


/*
 * Life 某个时间状态下的 Field 的集合
 * a，为第一状态，b为第二状态，如此循环
 */
type Life struct {
	a, b *Field
	w, h int
}


/*
 * NewLife 初始化随机数的 life
 * 避免程序重启后，产生的随机数和上次一样，用时间戳重置种子
 * rand.Intn 在[0,n)中取随机数
 * 返回一个随机的 Field 和 初始化的 Field
 */
func NewLife(w, h int) *Life {
	a := NewField(w, h)
	for i := 0; i < (w * h / 4); i++ {
		rand.Seed(time.Now().UnixNano())
		a.Set(rand.Intn(w), rand.Intn(h), true)
	}
	return &Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
	}
}

/*
 * 运行一步，重新计算并更新所有细胞状态。
 * a是之前状态，b 是下一个状态，a 走下一步，然后交换 a 和 b
 */
func (l *Life) Step() {
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}

	l.a, l.b = l.b, l.a
}

/**
 * 将细胞阵列按照字符串返回
 */
func (l *Life) String() string {
	var buf bytes.Buffer
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			b := byte('-')
			if l.a.Alive(x, y) {
				b = '*'
			}
			buf.WriteByte(b)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

/**
 * 清屏
 */
func clearScreem() {

	clear := make(map[string]func())

	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd :=exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func main() {

	l := NewLife(100, 25)

	for {
		l.Step()
		fmt.Print("\x0c", l) // Clear screen and print field.
		time.Sleep(time.Second/30)
		clearScreem()
	}
}
