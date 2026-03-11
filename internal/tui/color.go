package tui

import "charm.land/lipgloss/v2"

var (
	Black           = lipgloss.Color("#000000")
	DarkRed         = lipgloss.Color("#FF0000")
	Red             = lipgloss.Color("#FF5353")
	Purple          = lipgloss.Color("135")
	Orange          = lipgloss.Color("214")
	BurntOrange     = lipgloss.Color("214")
	Yellow          = lipgloss.Color("#DBBD70")
	Green           = lipgloss.Color("34")
	Turquoise       = lipgloss.Color("86")
	DarkGreen       = lipgloss.Color("#325451")
	LightGreen      = lipgloss.Color("47")
	GreenBlue       = lipgloss.Color("#00A095")
	DeepBlue        = lipgloss.Color("39")
	LightBlue       = lipgloss.Color("81")
	LightishBlue    = lipgloss.Color("75")
	Blue            = lipgloss.Color("63")
	Violet          = lipgloss.Color("13")
	Grey            = lipgloss.Color("#737373")
	LightGrey       = lipgloss.Color("245")
	LighterGrey     = lipgloss.Color("250")
	EvenLighterGrey = lipgloss.Color("253")
	DarkGrey        = lipgloss.Color("#606362")
	White           = lipgloss.Color("#ffffff")
	OffWhite        = lipgloss.Color("#a8a7a5")
	HotPink         = lipgloss.Color("200")
)

var (
	DebugLogLevel = Blue
	InfoLogLevel  = Turquoise  // In v2, just use the color directly or handle profile in renderer
	ErrorLogLevel = Red
	WarnLogLevel  = Yellow

	LogRecordAttributeKey = LightGrey

	HelpKey = lipgloss.Color("ff")

	HelpDesc = lipgloss.Color("248")

	InactivePreviewBorder = lipgloss.Color("244")

	CurrentBackground            = Grey
	CurrentForeground            = White
	SelectedBackground           = lipgloss.Color("110")
	SelectedForeground           = Black
	CurrentAndSelectedBackground = lipgloss.Color("117")
	CurrentAndSelectedForeground = Black

	TitleColor = lipgloss.Color("")

	GroupReportBackgroundColor = EvenLighterGrey
	TaskSummaryBackgroundColor = EvenLighterGrey

	ScrollPercentageBackground = DarkGrey
)
