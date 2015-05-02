package hunters

type Boat struct {
    Type string
    ID string
    Kommandant string
    Forward TorpedoSection
    Aft TorpedoSection
}

type TorpedoSection struct {
    Tubes []Torpedo
    Capacity int
    SteamReloads int
    ElectricReloads int
}

type Torpedo string

const (
    Electric Torpedo = "Electric"
    Steam Torpedo = "Steam"
)
