package hunters

type Boat struct {
	Type       string
	ID         string
	Kommandant string
	Forward    TorpedoSection
	Aft        TorpedoSection
	Damage
}

type TorpedoSection struct {
	Tubes           []Tube
	Capacity        int
	SteamReloads    int
	ElectricReloads int
}

type Tube struct {
	Number  int
	Torpedo Torpedo
}

type Torpedo string

const (
	Electric  Torpedo = "Electric"
	Steam     Torpedo = "Steam"
	EmptyTube Torpedo = "Empty"
)
