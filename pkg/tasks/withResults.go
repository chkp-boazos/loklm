package tasks

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

func WithResults(onSuccess string, onError string) TaskDecorator {
	return func(task Task) Task {
		return func() error {
			green := color.New(color.FgGreen).SprintFunc()
			red := color.New(color.FgRed).SprintFunc()
			secondLine := color.New(color.FgCyan, color.BgBlack).SprintFunc()

			icon := green("✔")
			description := ""
			emojiName := ":tada:"
			text := onSuccess

			err := task()
			if err != nil {
				icon = red("x")
				description = fmt.Sprintf("\n  %s", err)
				emojiName = ":skull:"
				text = onError
			}
			fmt.Println(emoji.Sprintf("%s %s %s %s", icon, text, emojiName, secondLine(description)))
			return err
		}
	}
}
