package manager

func getAndroidVersionInfo(apiLevel string) (string, string) {
	versionMap := map[string][2]string{
		"35": {"15", "Vanilla Ice Cream"},
		"34": {"14", "Upside Down Cake"},
		"33": {"13", "Tiramisu"},
		"32": {"12L", "S-V2"},
		"31": {"12", "Snow Cone"},
		"30": {"11", "Red Velvet Cake"},
		"29": {"10", "Quince Tart"},
		"28": {"9.0", "Pie"},
		"27": {"8.1", "Oreo"},
		"26": {"8.0", "Oreo"},
		"25": {"7.1", "Nougat"},
		"24": {"7.0", "Nougat"},
		"23": {"6.0", "Marshmallow"},
		"22": {"5.1", "Lollipop"},
		"21": {"5.0", "Lollipop"},
		"20": {"4.4W", "KitKat"},
		"19": {"4.4", "KitKat"},
		"18": {"4.3", "Jelly Bean"},
		"17": {"4.2", "Jelly Bean"},
		"16": {"4.1", "Jelly Bean"},
	}

	if info, ok := versionMap[apiLevel]; ok {
		return "Android " + info[0], info[1]
	}
	return "Android " + apiLevel, ""
}
