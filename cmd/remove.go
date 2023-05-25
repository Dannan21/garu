/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes an AUR package",
	Long: `Removes an AUR package from your system. The syntax is:
gaur remove [PACKAGE NAME]`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Removing package:" + "...")
		pkgName := "pkg"

		if len(args) >= 1 && args[0] != "" {
			//Capturando o nome do pacote
			pkgName = args[0]
		} else {
			//Erro(observar padrão, vai se repetir)
			fmt.Println("Error: no install arguments")
			return
		}

		rmvCmd := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", pkgName)

		pkgStdout, err := rmvCmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}

		go func() {
			if _, err := io.Copy(os.Stdout, pkgStdout); err != nil {
				fmt.Println("Error in stdou:", err)
			}
		}()

		stderr, err := rmvCmd.StderrPipe()
		if err != nil {
			fmt.Println("Error creating stderr pipe:", err)
			return
		}

		go func() {
			if _, err := io.Copy(os.Stderr, stderr); err != nil {
				fmt.Println("Error copying stderr:", err)
			}
		}()

		if err := rmvCmd.Start(); err != nil {
			fmt.Println("Error removing package:", err)
			return
		}

		if err := rmvCmd.Wait(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Package removed")
		}

	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
