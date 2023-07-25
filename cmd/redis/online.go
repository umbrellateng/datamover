/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/20 3:44 下午
 */
package redis

import (
	"context"
	"core.bank/datamover/utils"
	"fmt"
	"io"
	"os"
	"os/exec"

	"core.bank/datamover/log"

	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

var (
	from   string
	target string
	fromPassword string
	targetPassword string

	fromDB     int
	targetDB   int
)




// 定义不同的键类型
const (
	TypeString = "string"
	TypeList   = "list"
	TypeSet    = "set"
	TypeHash   = "hash"
	TypeZSet   = "zset"
)

var onlineCmd = &cobra.Command{
	Use:   "online",
	Short: "move redis data from source cluster target the target cluster",
	Args:  cobra.NoArgs,
	Run:   onlineCommandLibraryFunc,
}

func init() {
	onlineCmd.Flags().StringVarP(&from, "from", "f", "", "source redis cluster url")
	onlineCmd.Flags().StringVarP(&target, "target", "t", "", "target redis cluster url")
	onlineCmd.Flags().StringVar(&fromPassword, "from-password", "", "")
	onlineCmd.Flags().StringVar(&targetPassword, "target-password", "", "")
	onlineCmd.Flags().IntVar(&fromDB, "from-db", 0, "source redis db number ")
	onlineCmd.Flags().IntVar(&targetDB, "target-db", 0, "target redis db number ")

	_ = onlineCmd.MarkFlagRequired("from")
	_ = onlineCmd.MarkFlagRequired("target")
}

func onlineCommandLibraryFunc(cmd *cobra.Command, args[] string) {

	utils.HandleError(migrateAll())
}

// 定义一个函数，用于迁移所有的键值对
func migrateAll() error {
	var srcClient = redis.NewClient(&redis.Options{
		Addr:    from,
		Password: fromPassword,
		DB:       fromDB,
	})

	var dstClient = redis.NewClient(&redis.Options{
		Addr:   target,
		Password: targetPassword,
		DB:       targetDB,
	})
	var srcPipe = srcClient.Pipeline()
	var dstPipe = dstClient.Pipeline()

	defer srcClient.Close()
	defer dstClient.Close()

	var ctx = context.Background()

	// 使用 Scan 命令代替 Keys 命令，避免阻塞 redis 服务器
	var cursor uint64
	var keys []string
	for {
		var err error
		keys, cursor, err = srcClient.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			return err
		}
		// 遍历每个键名，调用迁移函数，并使用错误处理函数统一处理错误
		for _, key := range keys {
			// 获取键的类型
			keyType, err := srcClient.Type(ctx, key).Result()
			if err != nil {
				return err
			}
			utils.HandleError(migrateKey(ctx, srcClient, srcPipe, dstPipe, key, keyType))
		}
		// 如果游标为 0，表示扫描完成，退出循环
		if cursor == 0 {
			break
		}
	}
	// 迁移完成，返回 nil 错误
	return nil
}

// 定义一个函数，用于迁移一个键值对
func migrateKey(ctx context.Context, srcClient *redis.Client, srcPipe, dstPipe redis.Pipeliner, key, keyType string) error {
	var err error
	// 根据不同的类型，使用不同的命令获取值，并使用管道命令批量执行
	switch keyType {
	case TypeString:
		srcCmd := srcPipe.Get(ctx, key)
		_, err = srcPipe.Exec(ctx)
		if err != nil {
			return err
		}

		//value := srcClient.Get(ctx, key).Val()
		_ = dstPipe.Set(ctx, key, srcCmd.Val(), 0)
		log.Logger.Info("migrate key value: %s: %s", key, srcCmd.Val())

		_, err = dstPipe.Exec(ctx)
		if err != nil {
			return err
		}
	case TypeList:
		// 获取列表的长度
		length, err := srcClient.LLen(ctx, key).Result()
		if err != nil {
			return err
		}

		_, err = srcPipe.Exec(ctx)
		if err != nil {
			return err
		}
		// 获取列表的所有元素，并将元素插入到目标 redis 的列表中
		srcCmd := srcPipe.LRange(ctx, key, 0, length-1)
		_ = dstPipe.RPush(ctx, key, srcCmd.Val())

		_, err = dstPipe.Exec(ctx)
		if err != nil {
			return err
		}
	case TypeSet:
		// 获取集合的所有元素，并将元素添加到目标 redis 的集合中
		srcCmd := srcPipe.SMembers(ctx, key)
		_, err = srcPipe.Exec(ctx)
		if err != nil {
			return err
		}

		_ = dstPipe.SAdd(ctx, key, srcCmd.Val())
		_, err = dstPipe.Exec(ctx)
		if err != nil {
			return err
		}
	case TypeHash:
		// 获取哈希表的所有字段和值，并将字段和值设置到目标 redis 的哈希表中
		srcCmd := srcPipe.HGetAll(ctx, key)
		_, err = srcPipe.Exec(ctx)
		if err != nil {
			return err
		}

		_ = dstPipe.HSet(ctx, key, srcCmd.Val())

		_, err = dstPipe.Exec(ctx)
		if err != nil {
			return err
		}

	case TypeZSet:
		// 获取有序集合的所有元素和分数，并将元素和分数添加到目标 redis 的有序集合中
		srcCmd := srcPipe.ZRangeWithScores(ctx, key, 0, -1)
		_, err = srcPipe.Exec(ctx)
		if err != nil {
			return err
		}

		_ = dstPipe.ZAdd(ctx, key, Value2Pointer(srcCmd.Val())...)
		_, err = dstPipe.Exec(ctx)
		if err != nil {
			return err
		}
	default:
		// 不支持的类型，返回错误信息
		return fmt.Errorf("unsupported type: %s", keyType)
	}
	// 迁移成功，返回 nil 错误
	return nil
}

func Value2Pointer(datas []redis.Z) []*redis.Z {
	var ret []*redis.Z
	for _, data := range datas {
		ret = append(ret, &data)
	}
	return ret
}

func onlineCommandFunc(cmd *cobra.Command, args []string) {

	err := redisOnlineMove(from, target, fromDB)
	if err != nil {
		log.Logger.Error("redis data migration error: " + err.Error())
		return
	}

	fmt.Println()
	log.Logger.Info("redis data migration on success!")
}

func redisOnlineMove(sourceURL, targetURL string, db int) error {
	// 创建redis-cli命令对象
	cliCmd := exec.Command("redis-cli", "-u", sourceURL, "-n", fmt.Sprint(db), "keys", "*")
	// 执行redis-cli命令并获取所有键的列表
	keys, err := cliCmd.Output()
	if err != nil {
		return fmt.Errorf("redis cliCmd.Output error: " + err.Error())
	}

	log.Logger.Info("keys: " + string(keys))

	// 创建xargs命令对象
	xargsCmd := exec.Command("xargs", "-I", "{}", "redis-cli", "-u", sourceURL, "-n", fmt.Sprint(db),
		"migrate", targetURL, "{}", fmt.Sprint(db), "10000", "COPY")
	// 创建一个读写管道
	stdin, err := xargsCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("redis xargsCmd.StdinPipe error: " + err.Error())
	}
	stdout, err := xargsCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("redis xargsCmd.StdoutPipe error: " + err.Error())
	}
	// 将键的列表写入到管道中
	go func() {
		defer stdin.Close()
		_, _ = io.WriteString(stdin, string(keys))
	}()

	// 从管道中读取输出信息，并打印到标准输出中
	go func() {
		defer stdout.Close()
		_, _ = io.Copy(os.Stdout, stdout)
	}()
	// 执行xargs命令并等待完成
	err = xargsCmd.Run()
	if err != nil {
		return fmt.Errorf("xargsCmd.Run error: " + err.Error())
	}

	return nil
}
