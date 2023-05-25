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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update called")
		var tempPath = os.TempDir()

		command := exec.Command("pacman", "-Qm")
		output, err := command.Output()
		if err != nil {
			fmt.Println("Error:", err)
		}

		packages := strings.Split(strings.TrimSpace(string(output)), "\n")
		for i, pkg := range packages {
			//packages[i] = strings.Split(pkg, " ")[0]
			fmt.Println(pkg, i)
			URL := "https://aur.archlinux.org/cgit/aur.git/plain/PKGBUILD?h=" + pkg
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
				//Teste git

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
