package utils

import (
	"log"
	"os"
)

// printa uma mensagem de aviso
func PrintWarning(text string) {
	log.Println("\x1B[33maviso\033[0m", text)
}

// printa uma mensagem de erro
func PrintError(text string) {
	log.Println("\x1B[31merro\033[0m", text)
}

// printa uma mensagem de sucesso
func PrintSuccess(text string) {
	log.Println("\x1B[32msucesso\033[0m", text)
}

// printa uma mensagem informativa
func PrintInfo(text string) {
	log.Println("\x1B[34minfo\033[0m", text)
}

// printa uma mensagem na cor verde
func PrintGreen(text string) {
	log.Println("\x1B[32m" + text + "\033[0m")
}

// printa uma mensagem na cor vermelha
func PrintRed(text string) {
	log.Println("\x1B[31m" + text + "\033[0m")
}

// printa uma mensagem em negrito
func PrintBold(text string) {
	log.Println("\x1B[1m" + text + "\033[0m")
}

// printa uma mensagem e deixa o processo
func Fatal(err error) {
	log.Println("\x1B[31mfatal\033[0m", err)
	
	os.Exit(1)
}