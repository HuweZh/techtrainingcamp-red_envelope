package commons

import (
	"errors" // 生成错误
	"fmt"
	"sync" // 使用互斥锁
	"time" // 获取时间
)

/*
 * 这里求最大值使用了位运算，-1 的二进制表示为 1 的补码
 */

const (
	nodeBits  uint8 = 10                    // 节点 ID 的位数
	stepBits  uint8 = 12                    // 序列号的位数
	nodeMax   int64 = -1 ^ (-1 << nodeBits) // 节点 ID 的最大值，用于检测溢出
	stepMax   int64 = -1 ^ (-1 << stepBits) // 序列号的最大值，用于检测溢出
	timeShift uint8 = nodeBits + stepBits   // 时间戳向左的偏移量
	nodeShift uint8 = stepBits              // 节点 ID 向左的偏移量
)

//设置初始时间的时间戳 (毫秒表示) int64(time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano() / (1000 * 1000))
var Epoch int64 = 1514764800000

//这里我们申明一个 int64 的 ID 类型 （这样可以为此类型定义方法，比直接使用 int64 变量更灵活）
type ID int64

//创建channel
var snowflakeIDChann chan ID
var node *Node
var err error

//Node 结构用来存储一个节点 (机器) 上的基础数据
type Node struct {
	mu        sync.Mutex // 添加互斥锁，保证并发安全
	timestamp int64      // 时间戳部分
	node      int64      // 节点 ID 部分
	step      int64      // 序列号 ID 部分
}

//获取 Node 类型实例的函数，用于获得当前节点的 Node 实例
func newNode(node int64) (*Node, error) {
	// 如果超出节点的最大范围，产生一个 error
	if node < 0 || node > nodeMax {
		return nil, errors.New("Node number must be between 0 and 1023")
	}
	// 生成并返回节点实例的指针
	return &Node{
		timestamp: 0,
		node:      node,
		step:      0,
	}, nil
}

//最后一步，生成 ID 的方法
func (n *Node) generate() ID {

	n.mu.Lock()         // 保证并发安全, 加锁
	defer n.mu.Unlock() // 方法运行完毕后解锁

	// 获取当前时间的时间戳 (毫秒数显示)
	now := time.Now().UnixNano() / 1e6

	if n.timestamp == now {
		// step 步进 1
		n.step++

		// 当前 step 用完
		if n.step > stepMax {
			// 等待本毫秒结束
			for now <= n.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}

	} else {
		// 本毫秒内 step 用完
		n.step = 0
	}

	n.timestamp = now
	// 移位运算，生产最终 ID
	result := ID((now-Epoch)<<timeShift | (n.node << nodeShift) | (n.step))

	return result
}

//初始化snowflakeIDChann，并开启goroutine循环生成id
func init() {
	//新建有缓存的channel，缓存大小为1000个snowflakeID
	snowflakeIDChann = make(chan ID, 10000)
	// 生成节点实例
	node, err = newNode(1)
	if err != nil {
		fmt.Printf("create node [Err:%s]\n", err.Error())
		return
	}
	//开启一个协程循环生成snowflakeID
	go func() {
		for {
			id := node.generate()
			snowflakeIDChann <- id
		}
	}()
}

//获取一个雪花id
func GetID() ID {
	return <-snowflakeIDChann
}
