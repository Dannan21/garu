/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates all AUR packages",
	Long: `gets all AUR packages to the latest version. The syntax is:
gaur update`,
	Run: func(cmd *cobra.Command, args []string) {
		var tempPath = os.TempDir()

		command := exec.Command("pacman", "-Qm")
		output, err := command.Output()
		if err != nil {
			fmt.Println("Error:", err)
		}

		packages := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, pkg := range packages {
			//packages[i] = strings.Split(pkg, " ")[0]
			//fmt.Println(pkg, i)
			var pkgName = ""
			parts := strings.Fields(pkg)
			if len(parts) > 0 {
				pkgName = parts[0]
				fmt.Println(pkgName) // Output: Hello
			}
			URL := "https://aur.archlinux.org/cgit/aur.git/plain/PKGBUILD?h=" + pkgName
			response, err := http.Get(URL)
			if err != nil {
				fmt.Println(err)
			}
			if response.StatusCode == 200 {
				fmt.Println("==> Updates Found!")
				out, err := os.Create(tempPath + "/PKGBUILD")
				if err != nil {
					fmt.Println(err)
				}
				defer out.Close()

				_, err = io.Copy(out, response.Body)
				if err != nil {
					fmt.Println(err)
				}
				_ = exec.Command("sh", "-c", "cd "+tempPath+" && makepkg --printsrcinfo > "+tempPath+"/.SRCINFO")
				//outSRCINFO, err := getSRCINFO.CombinedOutput()

				if err != nil {
					fmt.Println(err)
				}

				makepkgCmd := exec.Command("sh", "-c", "cd "+tempPath+" && makepkg -i --noconfirm")

				pkgStdout, err := makepkgCmd.StdoutPipe()
				if err != nil {
					fmt.Println(err)
				}

				go func() {
					if _, err := io.Copy(os.Stdout, pkgStdout); err != nil {
						fmt.Println("Erro no stdout:", err)
					}
				}()

				stderr, err := makepkgCmd.StderrPipe()
				if err != nil {
					fmt.Println("Error creating stderr pipe:", err)
					return
				}

				go func() {
					if _, err := io.Copy(os.Stderr, stderr); err != nil {
						fmt.Println("Error copying stderr:", err)
					}
				}()

				if err := makepkgCmd.Start(); err != nil {
					fmt.Println("Erro no makepkg:", err)
					return
				}

				if err := makepkgCmd.Wait(); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("==> Update Installed")
				}

				_ = exec.Command("sudo", "rm", "-rf", tempPath+"/PKGBUILD", tempPath+"/.SRCINFO")
				if err != nil {
					fmt.Println(err)
				}

			} else {
				fmt.Println("==> Updates Not Found")
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
