# Garu - Garu Assists Repository Users

Garu is an AUR helper designed to assist users in managing packages from the Arch User Repository (AUR). It simplifies the process of building and installing AUR packages, making it easier for Arch Linux users to access and install software not available in the official repositories.

## Installation

To install Garu, you can clone the AUR package using git. Here are the steps:

1. Open a terminal.

2. Clone the Garu AUR package repository:
```
git clone https://aur.archlinux.org/garu.git
```
3. Run makepkg
```
makepkg -si
```

The `-s` flag installs the required dependencies automatically, and the `-i` flag installs the package after building.

## Usage

Garu provides a set of commands to manage AUR packages. Some common commands include:

- `garu install <package_name>`: Install an AUR package.
- `garu update`: Synchronize and upgrade all AUR packages.
- `garu remove <package name>`: Remove an AUR package and its dependencies.

## Contributing

If you encounter any issues, have suggestions, or would like to contribute to Garu, feel free to open an issue or submit a pull request on the [Garu GitHub repository](https://github.com/Dannan21/garu).

## License

Garu is licensed under the GPL3 License. See the [LICENSE](https://github.com/your_username/garu/blob/main/LICENSE) file for more details.

