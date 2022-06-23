package utilscmd

import (
	"fmt"
	"github.com/nodauf/web3Toolbox/utilsCmd/weiConverter"
	"math/big"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var weiStr, gweiStr, etherStr string

// weiConverterCmd represents the weiConverter command
var weiConverterCmd = &cobra.Command{
	Use:   "weiConverter",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		gwei, err := weiConverter.ParseBigFloat(gweiStr)
		if err != nil {
			fmt.Println("Fail to parse gwei:", err)
			return
		}
		ether, err := weiConverter.ParseBigFloat(etherStr)
		if err != nil {
			fmt.Println("Fail to parse ether:", err)
			return
		}
		wei, ok := (new(big.Int)).SetString(weiStr, 10)
		if !ok {
			fmt.Println("Fail to parse wei:")
			return
		}
		fmt.Println(gwei)
		fmt.Println(ether)
		// If no argument spawn the gui
		if wei.Cmp(big.NewInt(0)) == 1 || gwei.Cmp(big.NewFloat(0)) == 1 || ether.Cmp(big.NewFloat(0)) == 1 {
			if wei.Cmp(big.NewInt(0)) == 1 {
				etherBigFloat := weiConverter.WeiToEther(wei)
				fmt.Printf("ether\t %s\n", etherBigFloat.Text('f', -1))
				fmt.Printf("gwei:\t%s\n", weiConverter.WeiToGwei(wei).Text('f', -1))
			}
			if gwei.Cmp(big.NewFloat(0)) == 1 {
				weiBigInt := weiConverter.GweiTowei(gwei)
				fmt.Printf("wei: \t%s\n", weiBigInt.String())
				fmt.Printf("ether:\t%s\n", weiConverter.WeiToEther(weiBigInt).Text('f', -1))
			}
			if ether.Cmp(big.NewFloat(0)) == 1 {
				weiBigInt := weiConverter.EtherToWei(ether)
				fmt.Printf("wei:\t%s\n", weiBigInt.String())
				fmt.Printf("gwei:\t%s\n", weiConverter.WeiToGwei(weiBigInt).Text('f', -1))
			}
			return
		}
		app := tview.NewApplication()
		weiInput := tview.NewInputField().
			SetLabel("Wei: ").
			SetPlaceholder("1000000000000000000").
			SetFieldWidth(50).SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return unicode.IsDigit(lastChar)
		})
		gweiInput := tview.NewInputField().
			SetLabel("Gwei: ").
			SetPlaceholder("1000000000").
			SetFieldWidth(50).SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return unicode.IsDigit(lastChar)
		})
		etherInput := tview.NewInputField().
			SetLabel("Eth: ").
			SetPlaceholder("1").
			SetFieldWidth(50).SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return unicode.IsDigit(lastChar)
		})

		/***** Auto update part *****/

		weiInput.SetChangedFunc(func(weiString string) {
			// Only update the other field if it's this field that has the focus
			if weiInput.HasFocus() {
				wei, _ := (new(big.Int)).SetString(weiString, 10)
				eth := weiConverter.WeiToEther(wei)
				// Update ether field
				etherInput.SetText(eth.Text('f', -1))
				// Update gwei field
				gweiInput.SetText(weiConverter.WeiToGwei(wei).Text('f', -1))
			}
		})

		gweiInput.SetChangedFunc(func(gweiString string) {
			// Only update the other field if it's this field that has the focus
			if gweiInput.HasFocus() {
				gwei, _ := weiConverter.ParseBigFloat(gweiString)
				wei := weiConverter.GweiTowei(gwei)
				// Update wei field
				weiInput.SetText(wei.String())
				eth := weiConverter.WeiToEther(wei)
				// Update ether field
				etherInput.SetText(eth.Text('f', -1))
			}
		})

		etherInput.SetChangedFunc(func(etherString string) {
			// Only update the other field if it's this field that has the focus
			if etherInput.HasFocus() {
				ether, _ := weiConverter.ParseBigFloat(etherString)
				wei := weiConverter.EtherToWei(ether)
				// Update wei field
				weiInput.SetText(wei.String())
				// Update gwei field
				gweiInput.SetText(weiConverter.WeiToGwei(wei).Text('f', -1))
			}
		})
		form := tview.NewForm().AddFormItem(weiInput).AddFormItem(gweiInput).AddFormItem(etherInput).
			AddButton("Quit", func() {
				app.Stop()
			})
		form.SetBorder(true).SetTitle("Wei converter").SetTitleAlign(tview.AlignLeft)

		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEsc {
				app.Stop()
			}
			return event
		})
		if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	UtilsCmd.AddCommand(weiConverterCmd)
	weiConverterCmd.Flags().StringVarP(&weiStr, "wei", "w", "0", "wei")
	weiConverterCmd.Flags().StringVarP(&gweiStr, "gwei", "g", "0", "gwei")
	weiConverterCmd.Flags().StringVarP(&etherStr, "ether", "e", "0", "ether")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// weiConverterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// weiConverterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
