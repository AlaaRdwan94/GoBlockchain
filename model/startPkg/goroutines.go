package startPkg

import (
	"../heartbeat"

	"../BackUp"
	"../account"
	"../broadcastHandle"
	"../systemupdate"
	proof "../proofofstake"
)

func StartGoRoutine() {
	go account.ClearDeadUser()

	go proof.Mining()

	go heartbeat.SendHeartBeat()

	go systemupdate.Update()

	go broadcastHandle.OpenSocket(Conf.Server.PrivIP + ":" + Conf.Server.SoketPort)
	go BackUp.CreatBackup()
}
