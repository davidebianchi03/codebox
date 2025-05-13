import ConsoleIcon from "../../assets/icons/console.svg";
import DockerIcon from "../../assets/icons/docker.svg";
import JSONIcon from "../../assets/icons/json.svg";
import MarkdownIcon from "../../assets/icons/markdown.svg";
import PowerShellIcon from "../../assets/icons/powershell.svg";
import TerraformIcon from "../../assets/icons/terraform.svg";
import YamlIcon from "../../assets/icons/yaml.svg";
import FileIcon from "../../assets/icons/file.svg";


interface IconMap {
    extensions: string[]
    icon: string
}

const IconsMap: IconMap[] = [
    {
        extensions: ["sh"],
        icon: ConsoleIcon
    },
    {
        extensions: ["Dockerfile", "dockerfile"],
        icon: DockerIcon
    },
    {
        extensions: ["json"],
        icon: JSONIcon
    },
    {
        extensions: ["md"],
        icon: MarkdownIcon
    },
    {
        extensions: ["ps1"],
        icon: PowerShellIcon
    },
    {
        extensions: ["tf"],
        icon: TerraformIcon
    },
    {
        extensions: ["yaml", "yml"],
        icon: YamlIcon
    },
]

export function GetIconForFile(filename: string): string {
    var icon = FileIcon;
    IconsMap.forEach(ft => {
        ft.extensions.forEach(ex => {
            if (filename.endsWith(ex)) {
                icon = ft.icon;
            }
        })
    });
    return icon;
}
