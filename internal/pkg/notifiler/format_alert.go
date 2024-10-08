package notifiler

import (
	"fmt"
	"strings"
	"time"

	"github.com/lidofinance/finding-forwarder/generated/databus"
)

func FormatAlert(alert *databus.FindingDtoJson, source string) string {
	var (
		body   string
		footer string
	)

	if alert.Description != "" {
		body = alert.Description
		footer += "\n\n"
	}

	quorumTime := time.Now()
	if alert.BlockNumber != nil {
		eventToQuorumSecs := int(quorumTime.Unix()) - *alert.BlockTimestamp
		footer += fmt.Sprintf("Happened at block [%d](https://etherscan.io/block/%d/), ~ %d seconds ago", *alert.BlockNumber, *alert.BlockNumber, eventToQuorumSecs)
	}
	footer += fmt.Sprintf("\nTeam %s | %s | %s | quorum at %s by %s", alert.Team, alert.BotName, alert.AlertId, quorumTime.Format("15:04:05.000 MST"), source)

	if alert.TxHash != nil {
		footer += fmt.Sprintf("\nTx hash: [%s](https://etherscan.io/tx/%s/)", shortenHex(*alert.TxHash), *alert.TxHash)
	}

	return fmt.Sprintf("%s%s", body, footer)
}

func shortenHex(input string) string {
	if len(input) <= 5 {
		return input
	}
	return fmt.Sprintf("x%s...%s", input[2:5], input[len(input)-3:])
}

func TruncateMessageWithAlertID(message string, stringLimit int, warnMessage string) string {
	if len(message) <= stringLimit {
		return message
	}

	alertIndex := strings.LastIndex(message, "Alert Id:")
	if alertIndex == -1 {
		return fmt.Sprintf("%s\n%s", warnMessage, message[:stringLimit-len(warnMessage)-1])
	}

	alertText := message[alertIndex:]

	const formatSpecialCharsLength = 9
	maxTextLength := stringLimit - len(warnMessage) - len(alertText) - formatSpecialCharsLength

	if maxTextLength > 0 && alertIndex > maxTextLength {
		return fmt.Sprintf("%s\n...\n\n*%s*\n%s", message[:maxTextLength], warnMessage, alertText)
	}

	return fmt.Sprintf("%s\n%s", warnMessage, alertText)
}
