package command

func GeneralUsage() string {

	return "\"help [topic]\" to get a detailed explantion of the command.\nExample: \"help fly\"\nList of topics:\nfly\nset role\nview roles\nset airport\nview airports\ninvite\nrequest"
}

func FlyUsage() string {
	return "Usage: I (don't) want to fly\nExample: \"I don't want to fly\"\nor\n\"I want to fly\"\nIf you wanted to register for an airport:\n\"Help set airport\""
}

func RoleUsage() string {
	return "Usage: I want to be a [role]\nExample: \"I want to be an MO\""
}

func AirportUsage() string {
	return "Usage: I (don't) want to fly at [IATA Code]\nExample: I want to fly at KAWO\nor\nI don't want to fly at KAWO"
}

func ShowUsage() string {
	// CAREFUL, THIS MESSAGE IS EXACTLY 160 CHARACTERS (SMS Limit). Careful when appending something to this string
	return "Usage: I want to view [my/all] [roles/airports]\nExamples: \"I want to view my roles\" shows YOUR roles and \"I want to view all airports\" shows ALL airport options"
}

func ShowRoleUsage() string {
	return "\"I want to view my roles\" shows you YOUR roles\n\"I want to view all roles\" shows you ALL role options\n\"I want to view all detailed roles\" shows you ALL roles and what the acronym means"
}

func ShowAirportUsage() string {
	return "\"I want to view my airports\" shows you YOUR airports\n\"I want to view all airports\" shows you ALL airport options"
}

func InviteUsage() string {
	return "Usage: I want to invite [number]\nExample: \"I want to invite 3606851331\"\nor\n\"I want to invite 13606851331\""
}

func RequestUsage() string {
	return "Usage: I need a [role] at [IATA]\nExample: \"I need an MO at KBLI\""
}
