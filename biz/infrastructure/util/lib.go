package util

func GetMsg(resp interface{ GetMsg() string }, msg string) string {
	if resp == nil {
		return msg
	}
	return resp.GetMsg()
}
