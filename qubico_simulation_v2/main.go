package main

import (
	"log"
	"strconv"
)

const SLEEP_W = 300
const LOW_W = 2000
const DEFAULT_W = 3000
const LOW_TH = 75
const DEFAULT_TH = 100
const UP_RAMP_WPS = 50
const SITE_CONTAINERS = 10
const ASIC_PER_CONTAINER = 336
const VENTILATION_W = 36000 //18 vents each 2000W
const WARM_UP_TIME = 10     //wait time for asics to start ramping after they are turned on
const COOL_DOWN_TIME = 20   //wait time for the asic to be defined in sleep state after a down ramp (all asics have to be cooled down before turning off the container)

var target_W float64
var siteCurrent_W float64
var currentTime float64
var Containers []Container

// auxiliary variables to calculate target states
var MinContainerW float64
var AllLowContainerW float64
var MaxContainerW float64
var MaxSiteConsumption float64

var validAsicTargetStates = map[string]bool{
	"SLEEP":   true,
	"LOW":     true,
	"DEFAULT": true,
}

type Asic struct {
	AsicId       string
	TargetState  string //SLEEP, LOW, DEFAULT
	CurrentState string //OFF, WARM_UP, SLEEP, RAMP_UP, LOW, DEFAULT, COOL_DOWN (PREVER EL PERFIL CON ANTICIPACIÓN)
	Current_W    float64
	Current_TH   float64
	ContainerId  int
	WaitTime     float64
}

type Container struct {
	ContainerId    int
	TargetState    string //OFF, ON
	CurrentState   string //OFF, ON
	ContainerAsics []Asic
}

func main() {
	InitializeSiteState()
	UpdateSite(0)
	UpdateTargetStates(0.5)
	for t := 0; t < 120; t++ {
		UpdateSite(1)
		currentTime = float64(t)
		//print status every 10 seconds
		if t%5 == 0 {
			PrintCurrentStatus()
		}
	}
}

func InitializeSiteState() {
	currentTime = 0
	MinContainerW = float64(VENTILATION_W+SLEEP_W*ASIC_PER_CONTAINER) / 1000000
	AllLowContainerW = float64(VENTILATION_W+LOW_W*ASIC_PER_CONTAINER) / 1000000
	MaxContainerW = float64(VENTILATION_W+DEFAULT_W*ASIC_PER_CONTAINER) / 1000000
	MaxSiteConsumption = float64(MaxContainerW * SITE_CONTAINERS)

	target_W = 0
	for i := 0; i < SITE_CONTAINERS; i++ {
		new_container := Container{
			ContainerId:  i,
			TargetState:  "OFF",
			CurrentState: "OFF",
		}
		Containers = append(Containers, new_container)
		for j := 0; j < ASIC_PER_CONTAINER; j++ {
			asic_id := strconv.Itoa(i) + "_" + strconv.Itoa(j)
			new_asic := Asic{
				AsicId:       asic_id,
				TargetState:  "SLEEP",
				CurrentState: "OFF",
				Current_W:    0,
				Current_TH:   0,
				ContainerId:  i,
			}
			Containers[i].ContainerAsics = append(Containers[i].ContainerAsics, new_asic)
		}
	}
}

func GetCurrentTotal_MW() float64 {
	var totalW float64
	totalW = 0
	for i := range Containers {
		if Containers[i].CurrentState == "ON" {
			totalW = totalW + VENTILATION_W
			for j := 0; j < ASIC_PER_CONTAINER; j++ {
				totalW = totalW + Containers[i].ContainerAsics[j].Current_W
			}
		}
	}
	return totalW / 1000000
}

func ChangeAsicTargetState(asic *Asic, targetState string) {
	if ValidateAsicTargetState(targetState) {
		asic.TargetState = targetState
	}
}

func ValidateAsicTargetState(input string) bool {
	if !validAsicTargetStates[input] {
		log.Println("Invalide Asic State Requested: " + input)
		return false
	}
	return true
}

func UpdateAsic(deltaTime float64, asic *Asic, container *Container) {
	//ASIC states: OFF, WARM_UP, SLEEP, RAMP_UP, LOW, DEFAULT, COOL_DOWM

	//if container is OFF, ASIC also goes to off, and we reset ASIC
	if container.CurrentState == "OFF" {
		asic.CurrentState = "OFF"
		asic.WaitTime = 0
		asic.Current_W = 0
		asic.Current_TH = 0
		return
	}

	//When OFF
	if asic.CurrentState == "OFF" {
		//if the container is on, go to WARM_UP and start timer
		if container.CurrentState == "ON" {
			asic.CurrentState = "WARM_UP"
			asic.Current_W = SLEEP_W
			asic.Current_TH = 0
			asic.WaitTime = WARM_UP_TIME
		}
		return
	}

	//When in WARM_UP, deduct time from timer and go to SLEEP if time is up
	if asic.CurrentState == "WARM_UP" {
		asic.WaitTime = asic.WaitTime - deltaTime
		if asic.WaitTime <= 0 {
			asic.CurrentState = "SLEEP"
			asic.WaitTime = 0
		}
		return
	}

	//When Sleep, if target state is LOW or DEFAULT, start ramping... otherwise keep SLEEP
	if asic.CurrentState == "SLEEP" {
		if asic.TargetState == "LOW" || asic.TargetState == "DEFAULT" {
			asic.CurrentState = "RAMP_UP"
		}
		return
	}

	//If cooling down, we update timer and go to sleep if enough time has passed
	if asic.CurrentState == "COOL_DOWN" {
		asic.WaitTime = asic.WaitTime - deltaTime
		if asic.WaitTime <= 0 {
			asic.CurrentState = "SLEEP"
			asic.WaitTime = 0
		}
		return
	}

	//If in LOW, we update towards target
	if asic.CurrentState == "LOW" {
		if asic.TargetState == "DEFAULT" {
			asic.CurrentState = "RAMP_UP"
		}
		if asic.TargetState == "SLEEP" {
			asic.CurrentState = "COOL_DOWN"
			asic.Current_W = SLEEP_W
			asic.Current_TH = 0
			asic.WaitTime = COOL_DOWN_TIME
		}
		return
	}
	//If in DEFAULT, we update towards target
	if asic.CurrentState == "DEFAULT" {
		if asic.TargetState == "LOW" {
			asic.CurrentState = "LOW"
			asic.Current_W = LOW_W
			asic.Current_TH = LOW_TH
		}
		if asic.TargetState == "SLEEP" {
			asic.CurrentState = "COOL_DOWN"
			asic.Current_W = SLEEP_W
			asic.Current_TH = 0
			asic.WaitTime = COOL_DOWN_TIME
		}
		return
	}

	//When Ramping, we update consumption and hashrate and stop if target is reached
	if asic.CurrentState == "RAMP_UP" {
		//if target is SLEEP, go to COOL_DOWN
		if asic.TargetState == "SLEEP" {
			asic.CurrentState = "COOL_DOWN"
			asic.Current_W = SLEEP_W
			asic.Current_TH = 0
			return
		}

		//if we are below target consumption we keep increasing consumption and hashrate
		//we define ramp TH speed based on current consumption trench
		var th_ramp float64
		th_ramp = float64(LOW_TH) / float64(LOW_W)
		if asic.Current_W > LOW_TH {
			th_ramp = float64(DEFAULT_TH-LOW_TH) / float64(DEFAULT_W-LOW_W)
		}
		//update hashrate and W
		asic.Current_W = asic.Current_W + deltaTime*UP_RAMP_WPS
		asic.Current_TH = asic.Current_TH + deltaTime*UP_RAMP_WPS*th_ramp
		//stop if reached target
		if asic.TargetState == "LOW" && asic.Current_W >= LOW_W {
			asic.CurrentState = "LOW"
			asic.Current_W = LOW_W
			asic.Current_TH = LOW_TH
		}
		if asic.TargetState == "DEFAULT" && asic.Current_W >= DEFAULT_W {
			asic.CurrentState = "DEFAULT"
			asic.Current_W = DEFAULT_W
			asic.Current_TH = DEFAULT_TH
		}
		return
	}
}

func UpdateContainer(deltaTime float64, container *Container) {
	//if the container is off and it's target state is on, we turn the container on
	if container.CurrentState == "OFF" && container.TargetState == "ON" {
		container.CurrentState = "ON"
		return
	}

	allAsicsInSleep := true
	//we update all Asics within the container
	for i := range container.ContainerAsics {
		UpdateAsic(deltaTime, &container.ContainerAsics[i], container)
		//after that update if there is any asic that is not in sleep state we cannot turn the container off
		if container.ContainerAsics[i].CurrentState != "SLEEP" {
			allAsicsInSleep = false
		}
	}
	//if the target state of the container is off, and all asics are in sleep state, we can turn the container off
	if allAsicsInSleep && container.TargetState == "OFF" {
		container.CurrentState = "OFF"
	}
}

func UpdateSite(deltaTime float64) {
	//update current state
	for i := range Containers {
		UpdateContainer(deltaTime, &Containers[i])
	}
}

func ChangeContainerTargetState(container *Container, targetState string) {
	container.TargetState = targetState
}

func PrintCointainerState(container Container) {
	log.Println("Container", container.ContainerId, "state is", container.CurrentState, "Target state is", container.TargetState)
	for i := range container.ContainerAsics {
		PrintAsicState(container.ContainerAsics[i])
	}
}

func PrintAsicState(asic Asic) {
	log.Println(asic.AsicId, "Current state:", asic.CurrentState, "Target State:", asic.TargetState, "Current_W:", asic.Current_W, "Current_TH", asic.Current_TH)
}

func UpdateTargetStates(target float64) {
	target_W = target
	//calculate full containers
	var nFullContainers int
	nFullContainers = min(int(target/MaxContainerW), SITE_CONTAINERS)
	//calculate consumption for the marginal container
	var marginalW = min(target, MaxContainerW*SITE_CONTAINERS) - float64(nFullContainers)*MaxContainerW

	//calculate distribution of asic target states for the marginal container
	var marginalContainerOn bool
	var nMarginalLow int
	var nMarginalDefault int
	var auxMarginalForDefault float64
	var auxDefaultLowDelta float64
	var auxMarginalForLow float64
	var auxLowSleepDelta float64

	nMarginalLow = 0
	nMarginalDefault = 0

	marginalContainerOn = true
	if marginalW < MinContainerW {
		marginalContainerOn = false
	}

	if marginalContainerOn {
		if marginalW <= AllLowContainerW {
			auxMarginalForLow = AllLowContainerW - marginalW
			auxLowSleepDelta = float64(LOW_W-SLEEP_W) / 1000000
			nMarginalLow = int(auxMarginalForLow / auxLowSleepDelta)
		}
		if marginalW > AllLowContainerW {
			auxMarginalForDefault = MaxContainerW - marginalW
			auxDefaultLowDelta = float64(DEFAULT_W-LOW_W) / 1000000
			nMarginalDefault = int(auxMarginalForDefault / auxDefaultLowDelta)
			nMarginalLow = ASIC_PER_CONTAINER - nMarginalDefault
		}
	}
	//update targets for full containers
	for i := 0; i < nFullContainers; i++ {
		Containers[i].TargetState = "ON"
		for j := 0; j < ASIC_PER_CONTAINER; j++ {
			Containers[i].ContainerAsics[j].TargetState = "DEFAULT"
		}
	}
	//update marginal container state
	if marginalContainerOn {
		Containers[nFullContainers].TargetState = "ON"
	}
	//update the marginal container's asics' targets
	if nMarginalDefault > 0 {
		for j := 0; j < ASIC_PER_CONTAINER; j++ {
			if j < nMarginalDefault {
				Containers[nFullContainers].ContainerAsics[j].TargetState = "DEFAULT"
			} else {
				Containers[nFullContainers].ContainerAsics[j].TargetState = "LOW"
			}
		}
	} else {
		for j := 0; j < ASIC_PER_CONTAINER; j++ {
			if j < nMarginalLow {
				Containers[nFullContainers].ContainerAsics[j].TargetState = "LOW"
			} else {
				Containers[nFullContainers].ContainerAsics[j].TargetState = "SLEEP"
			}
		}
	}
}

func PrintCurrentStatus() {
	currentConsumption := GetCurrentTotal_MW()
	log.Println("Elapsed Time", currentTime)
	log.Println("Target Consumption:", target_W)
	log.Println("Current Consumption:", currentConsumption)
	for i := 0; i < SITE_CONTAINERS; i++ {
		nOff := 0
		nSleep := 0
		nLow := 0
		nDefault := 0
		nRamp := 0
		nCoolDown := 0
		nWarmUp := 0
		for j := 0; j < ASIC_PER_CONTAINER; j++ {
			if Containers[i].ContainerAsics[j].CurrentState == "OFF" {
				nOff++
			}
			if Containers[i].ContainerAsics[j].CurrentState == "SLEEP" {
				nSleep++
			}
			if Containers[i].ContainerAsics[j].CurrentState == "LOW" {
				nLow++
			}
			if Containers[i].ContainerAsics[j].CurrentState == "DEFAULT" {
				nDefault++
			}
			if Containers[i].ContainerAsics[j].CurrentState == "RAMP" {
				nRamp++
			}
			if Containers[i].ContainerAsics[j].CurrentState == "COOL_DOWN" {
				nCoolDown++
			}
			if Containers[i].ContainerAsics[j].CurrentState == "WARM_UP" {
				nWarmUp++
			}
		}
		log.Println("Container:", Containers[i].CurrentState, "OFF:", nOff, "SLEEP", nSleep, "LOW", nLow, "DEFAULT", nDefault, "RAMP", nRamp, "COOL_DOWN", nCoolDown, "WARM_UP", nWarmUp)
	}
}
