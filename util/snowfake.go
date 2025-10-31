package util

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// SnowflakeGenerator 雪花ID生成器
type SnowflakeGenerator struct {
	mutex        sync.Mutex
	epoch        int64 // 起始时间戳 (2024-01-01 00:00:00 UTC)
	machineID    int64 // 机器ID (0-1023)
	datacenterID int64 // 数据中心ID (0-31)
	sequence     int64 // 序列号 (0-4095)
	lastTime     int64 // 上次生成ID的时间戳
}

// 雪花ID的位数分配
const (
	// 时间戳占用41位，可以使用69年
	timestampBits = 41
	// 数据中心ID占用5位，支持32个数据中心
	datacenterIDBits = 5
	// 机器ID占用5位，支持32台机器
	machineIDBits = 5
	// 序列号占用12位，支持每毫秒4096个ID
	sequenceBits = 12

	// 最大值
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits) // 31
	maxMachineID    = -1 ^ (-1 << machineIDBits)    // 31
	maxSequence     = -1 ^ (-1 << sequenceBits)     // 4095

	// 位移
	machineIDShift    = sequenceBits
	datacenterIDShift = sequenceBits + machineIDBits
	timestampShift    = sequenceBits + machineIDBits + datacenterIDBits
)

// NewSnowflakeGenerator 创建雪花ID生成器
func NewSnowflakeGenerator(machineID, datacenterID int64) (*SnowflakeGenerator, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("机器ID必须在0-31之间")
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("数据中心ID必须在0-31之间")
	}

	return &SnowflakeGenerator{
		epoch:        1640995200000, // 2024-01-01 00:00:00 UTC 的时间戳(毫秒)
		machineID:    machineID,
		datacenterID: datacenterID,
		sequence:     0,
		lastTime:     -1,
	}, nil
}

// GenerateID 生成雪花ID
func (s *SnowflakeGenerator) GenerateID() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixMilli()

	// 如果当前时间小于上次生成ID的时间，说明系统时钟回退
	if now < s.lastTime {
		return 0, errors.New("系统时钟回退，无法生成ID")
	}

	// 如果是同一毫秒内，序列号递增
	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		// 如果序列号溢出，等待下一毫秒
		if s.sequence == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 新的毫秒，序列号重置
		s.sequence = 0
	}

	s.lastTime = now

	// 生成ID
	id := ((now - s.epoch) << timestampShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.machineID << machineIDShift) |
		s.sequence

	return id, nil
}

// GenerateStringID 生成字符串格式的ID
func (s *SnowflakeGenerator) GenerateStringID() (string, error) {
	id, err := s.GenerateID()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

// 全局实例
var Snowflake *SnowflakeGenerator

// InitSnowflake 初始化雪花ID生成器
func InitSnowflake(machineID, datacenterID int64) error {
	generator, err := NewSnowflakeGenerator(machineID, datacenterID)
	if err != nil {
		return err
	}
	Snowflake = generator
	return nil
}

// GenerateUserID 生成用户ID的便捷方法
func GenerateUserID() (string, error) {
	if Snowflake == nil {
		return "", errors.New("雪花ID生成器未初始化")
	}

	id, err := Snowflake.GenerateID()
	if err != nil {
		return "", err
	}

	return string(rune(id)), nil
}
