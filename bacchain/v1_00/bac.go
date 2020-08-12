package v1_00

var (
	//[must_modify] 在主网启动或者升级的时候需要修改 StartHeight 和StartBacAlreadyProduce
	StartParamBacAlreadyProduce =  "29016000000000000"

	StartParamInitHeight = int64(2800800)

	StartParamBeginGenBac = int64(28800)


)

/**
  	true 可以转账
  	false 不可以转账
 */
func CheckSendEnable( genesisSendEnable bool,height int64)  bool {
	return genesisSendEnable
}