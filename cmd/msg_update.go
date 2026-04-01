package cmd

import (
	"fmt"
	"os"

	"github.com/riba2534/feishu-cli/internal/client"
	"github.com/riba2534/feishu-cli/internal/config"
	"github.com/spf13/cobra"
)

var msgUpdateCmd = &cobra.Command{
	Use:   "update <message_id>",
	Short: "更新卡片消息内容",
	Long: `更新已发送的卡片消息（interactive 类型）。

参数:
  message_id           消息 ID（必填）
  --content, -c        新的卡片 JSON 内容
  --content-file       从文件读取卡片 JSON

注意:
  - 只能更新 interactive（卡片）类型的消息
  - 卡片 JSON 必须包含 "config": {"update_multi": true}
  - 只能更新 14 天内发送的消息
  - 频率限制: 5 QPS / message

示例:
  # 更新卡片消息
  feishu-cli msg update om_xxx -c '{"schema":"2.0","config":{"update_multi":true},...}'

  # 从文件读取卡片 JSON
  feishu-cli msg update om_xxx --content-file card.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Validate(); err != nil {
			return err
		}

		token := resolveOptionalUserToken(cmd)

		messageID := args[0]
		content, _ := cmd.Flags().GetString("content")
		contentFile, _ := cmd.Flags().GetString("content-file")

		// 互斥校验
		if content == "" && contentFile == "" {
			return fmt.Errorf("请指定 --content 或 --content-file")
		}
		if content != "" && contentFile != "" {
			return fmt.Errorf("--content 和 --content-file 不能同时指定")
		}

		// 从文件读取
		if contentFile != "" {
			data, err := os.ReadFile(contentFile)
			if err != nil {
				return fmt.Errorf("读取文件失败: %w", err)
			}
			content = string(data)
		}

		if err := client.UpdateMessage(messageID, content, token); err != nil {
			return err
		}

		output, _ := cmd.Flags().GetString("output")
		if output == "json" {
			if err := printJSON(map[string]any{
				"success":    true,
				"message_id": messageID,
			}); err != nil {
				return err
			}
		} else {
			fmt.Printf("消息更新成功！\n")
			fmt.Printf("  消息 ID: %s\n", messageID)
		}

		return nil
	},
}

func init() {
	msgCmd.AddCommand(msgUpdateCmd)
	msgUpdateCmd.Flags().StringP("content", "c", "", "卡片 JSON 内容")
	msgUpdateCmd.Flags().String("content-file", "", "卡片 JSON 文件路径")
	msgUpdateCmd.Flags().StringP("output", "o", "", "输出格式（json）")
	msgUpdateCmd.Flags().String("user-access-token", "", "User Access Token（用户授权令牌）")
}
