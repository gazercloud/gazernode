package gazer_dictionary

func ChannelUrl(channelId string) string {
	if len(channelId) > 0 {
		return "https://gazer.cloud/channel/" + channelId
	}
	return ""
}
