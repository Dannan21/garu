/*
Copyright © 2023 Davi Seidel Brandão <daviseidel.brandao@gmail.com>
*/
package cmd

import (
	//Standart Library
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	//Lib para parsear PKGBUILD

	//Cobra para cli
	"github.com/spf13/cobra"
)

// installCmd represents the install command

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs an AUR package",
	Long: `Will install a package from the AUR in you system. The syntax is:
gaur install [PACKAGE NAME]`,
	Run: func(cmd *cobra.Command, args []string) {
		//Criando as váriaves que vão se usadas mais pra fente
		//Nome do pacote
		var pkgName = ""
		//Diretório /tmp
		var tempPath = os.TempDir()
		//Checando se há um argumento
		if len(args) >= 1 && args[0] != "" {
			//Capturando o nome do pacote
			pkgName = args[0]
		} else {
			//Erro(observar padrão, vai se repetir)
			fmt.Println("Error: no install arguments")
			return
		}
		//Printando o pacote
		fmt.Println("==> Package to be installed: " + pkgName)
		//Capturando o scipt
		URL := "https://aur.archlinux.org/cgit/aur.git/plain/PKGBUILD?h=" + pkgName
		//Fazendo o request
		response, err := http.Get(URL)
		if err != nil {
			fmt.Println(err)
		}

		//Checando o sucesso
		if response.StatusCode == 200 {
			fmt.Println("==> Package Found!")
			//Se sim, crie o PKGBUILD no diretório temporário
			out, err := os.Create(tempPath + "/PKGBUILD")
			if err != nil {
				fmt.Println(err)
			}
			defer out.Close()

			//Copiar para o PKGBUILD vazio
			_, err = io.Copy(out, response.Body)
			if err != nil {
				fmt.Println(err)
			}

			_ = exec.Command("sh", "-c", "cd "+tempPath+" && makepkg --printsrcinfo > "+tempPath+"/.SRCINFO")
			//outSRCINFO, err := getSRCINFO.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}

			//Ver Dependencias (que dor de cabeça, ainda bem q passou)
			/*
				pkgInfo, err := pkgbuild.ParseSRCINFO(tempPath + "/.SRCINFO")
				if err != nil {
					fmt.Println(err)
				}

					fmt.Println("==> Dependency list:")

					pkgBultDeps := pkgInfo.BuildDepends()

					for _, dep := range pkgBultDeps {
						//depString := dep.String()
						fmt.Println(dep.Name)
						if dep.MinVer != nil {
							fmt.Println(dep.MinVer.Version)
						}
					}
			*/
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
				fmt.Println("No errors!")
			}

		} else {
			fmt.Println("==> Package Not Found")
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("nodeps", "n", false, "Ignore dependency list")
}
