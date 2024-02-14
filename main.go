package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/creack/pty"
)

const (
	ScreenName = "nethack_select"
)

func main() {
	closeAllScreen()

	c := exec.Command("/usr/bin/screen", "-S", ScreenName)
	pf, err := pty.Start(c)
	if err != nil {
		panic(err)
	}
	defer closeAllScreen()
	defer pf.Close()

	// must wait
	time.Sleep(time.Second)

	screenList := getScreenList()
	if len(screenList) != 1 {
		fmt.Printf("find screen error: %v\n", screenList)
		return
	}
	name := screenList[0]

	for {
		pf.WriteString("nethack\n")
		pf.WriteString("          i")

		time.Sleep(100 * time.Millisecond)
		inventory, err := getInventory(name)
		for err != nil {
			fmt.Print(".")
			time.Sleep(100 * time.Millisecond)
			inventory, err = getInventory(name)
		}
		fmt.Println()
		fmt.Println(inventory)

		itemList := []string{
			"spellbook of identify",
			"spellbook of extra healing",
			"magic marker",
		}
		satisfied := true
		for _, item := range itemList {
			if !strings.Contains(inventory, item) {
				fmt.Printf("item not satisfied: %v\n", item)
				//fmt.Printf("screen -x %v\n", name)
				satisfied = false
				break
			}
		}
		if satisfied {
			pf.WriteString("   #save\ny")
			break
		} else {
			pf.WriteString("   #quit\nyq")
		}

		pf.WriteString("\n\n\n\n\n")
		// must wait
		time.Sleep(1 * time.Second)
	}
}

func getScreenList() []string {
	var list []string
	cmd := exec.Command("/usr/bin/screen", "-ls")
	output, _ := cmd.Output()
	reg := regexp.MustCompile(`[^\s]*` + ScreenName + `[^\s]*`)
	for _, line := range strings.Split(string(output), "\n") {
		matsch := reg.FindString(line)
		if matsch != "" {
			list = append(list, matsch)
		}
	}
	return list
}

func closeAllScreen() {
	screenList := getScreenList()
	if len(screenList) > 0 {
		exec.Command("/usr/bin/screen", "-wipe").Run()
		for _, v := range screenList {
			exec.Command("/usr/bin/screen", "-S", v, "-X", "quit").Run()
		}
	}
}

func getInventory(name string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("screen -S %v -X hardcopy /tmp/hardcopy.txt", name))
	cmd.Run()
	cmd = exec.Command("/bin/sh", "-c", `grep -o -E "[a-z] - .*" /tmp/hardcopy.txt`)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func printHardcopy(name string) {
	cmd := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("screen -S %v -X hardcopy /tmp/hardcopy.txt", name))
	cmd.Run()
	cmd = exec.Command("/bin/sh", "-c", `grep -vE "^\s*$" /tmp/hardcopy.txt`)
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
}
