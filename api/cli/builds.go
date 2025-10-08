package cli

type CLIBuild struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Os           string `json:"os"`
	Architecture string `json:"architecture"`
	Type         string `json:"type"`
	File         string `json:"-"` // the name of the file in the fs
}

var CliBuilds = []CLIBuild{
	// architecture: arm64
	CLIBuild{
		Id:           "codebox-cli-setup-arm64",
		Name:         "Codebox CLI Windows Installer (arm64)",
		Os:           "windows",
		Architecture: "arm64",
		Type:         "package",
		File:         "codebox-cli-setup-arm64.exe",
	},
	CLIBuild{
		Id:           "codebox-cli-windows-arm64",
		Name:         "Codebox CLI Windows Binary (arm64)",
		Os:           "windows",
		Architecture: "arm64",
		Type:         "binaries",
		File:         "codebox-cli-windows-arm64.exe",
	},
	CLIBuild{
		Id:           "codebox-cli-arm64-deb",
		Name:         "Codebox CLI Debian Package (arm64)",
		Os:           "linux",
		Architecture: "arm64",
		Type:         "package",
		File:         "codebox-cli-arm64.deb",
	},
	CLIBuild{
		Id:           "codebox-cli-linux-arm64",
		Name:         "Codebox CLI Linux Binary (arm64)",
		Os:           "linux",
		Architecture: "arm64",
		Type:         "binaries",
		File:         "codebox-cli-linux-arm64",
	},
	// architecture: amd64
	CLIBuild{
		Id:           "codebox-cli-setup-amd64",
		Name:         "Codebox CLI Windows Installer (amd64)",
		Os:           "windows",
		Architecture: "amd64",
		Type:         "package",
		File:         "codebox-cli-setup-amd64.exe",
	},
	CLIBuild{
		Id:           "codebox-cli-windows-amd64",
		Name:         "Codebox CLI Windows Binary (amd64)",
		Os:           "windows",
		Architecture: "amd64",
		Type:         "binaries",
		File:         "codebox-cli-windows-amd64.exe",
	},
	CLIBuild{
		Id:           "codebox-cli-amd64-deb",
		Name:         "Codebox CLI Debian Package (amd64)",
		Os:           "linux",
		Architecture: "amd64",
		Type:         "package",
		File:         "codebox-cli-amd64.deb",
	},
	CLIBuild{
		Id:           "codebox-cli-linux-amd64",
		Name:         "Codebox CLI Linux Binary (amd64)",
		Os:           "linux",
		Architecture: "amd64",
		Type:         "binaries",
		File:         "codebox-cli-linux-amd64",
	},
	// architecture: arm
	CLIBuild{
		Id:           "codebox-cli-setup-arm",
		Name:         "Codebox CLI Windows Installer (arm)",
		Os:           "windows",
		Architecture: "arm",
		Type:         "package",
		File:         "codebox-cli-setup-arm.exe",
	},
	CLIBuild{
		Id:           "codebox-cli-windows-arm",
		Name:         "Codebox CLI Windows Binary (arm)",
		Os:           "windows",
		Architecture: "arm",
		Type:         "binaries",
		File:         "codebox-cli-windows-arm.exe",
	},
	CLIBuild{
		Id:           "codebox-cli-arm-deb",
		Name:         "Codebox CLI Debian Package (arm)",
		Os:           "linux",
		Architecture: "arm",
		Type:         "package",
		File:         "codebox-cli-arm.deb",
	},
	CLIBuild{
		Id:           "codebox-cli-linux-arm",
		Name:         "Codebox CLI Linux Binary (arm)",
		Os:           "linux",
		Architecture: "arm",
		Type:         "binaries",
		File:         "codebox-cli-linux-arm",
	},
	// architecture: 386
	CLIBuild{
		Id:           "codebox-cli-setup-386",
		Name:         "Codebox CLI Windows Installer (386)",
		Os:           "windows",
		Architecture: "386",
		Type:         "package",
		File:         "codebox-cli-setup-386.exe",
	},
	CLIBuild{
		Id:           "codebox-cli-windows-386",
		Name:         "Codebox CLI Windows Binary (386)",
		Os:           "windows",
		Architecture: "386",
		Type:         "binaries",
		File:         "codebox-cli-windows-386.exe",
	},
	CLIBuild{
		Id:           "codebox-cli-386-deb",
		Name:         "Codebox CLI Debian Package (386)",
		Os:           "linux",
		Architecture: "386",
		Type:         "package",
		File:         "codebox-cli-386.deb",
	},
	CLIBuild{
		Id:           "codebox-cli-linux-386",
		Name:         "Codebox CLI Linux Binary (386)",
		Os:           "linux",
		Architecture: "386",
		Type:         "binaries",
		File:         "codebox-cli-linux-386",
	},
}
