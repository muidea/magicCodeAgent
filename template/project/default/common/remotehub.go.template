package common

import cd "github.com/muidea/magicCommon/def"

const (
	SubscribeEvent   = "/remoteHub/event/subscribe/"
	UnsubscribeEvent = "/remoteHub/event/unsubscribe/"
	QueryEvent       = "/remoteHub/event/query/"
	NotifyEvent      = "/remoteHub/event/notify/"
)

type SubscribeParam struct {
	Event    []string `json:"event"`
	CallBack string   `json:"callBack"`
}

type SubscribeResult struct {
	cd.Result
}

type UnsubscribeParam SubscribeParam
type UnsubscribeResult SubscribeResult

type QuerySubscribeResult struct {
	cd.Result
	Event []string `json:"event"`
}

type Event struct {
	Event  string `json:"event"`
	Action int    `json:"action"`
	DataID int    `json:"dataID"`
}

const RemoteHubModule = "/kernel/remotehub"
