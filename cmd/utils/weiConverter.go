package utilscmd

import (
	"fmt"
	"math/big"
	"strings"
	"unicode"

	"github.com/ethereum/go-ethereum/params"
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
		gwei, err := ParseBigFloat(gweiStr)
		if err != nil {
			fmt.Println("Fail to parse gwei:", err)
			return
		}
		ether, err := ParseBigFloat(etherStr)
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
				etherBigFloat := weiToEther(wei)
				fmt.Printf("ether\t %s\n", etherBigFloat.Text('f', -1))
				fmt.Printf("gwei:\t%s\n", weiToGwei(wei).Text('f', -1))
			}
			if gwei.Cmp(big.NewFloat(0)) == 1 {
				weiBigInt := gweiTowei(gwei)
				fmt.Printf("wei: \t%s\n", weiBigInt.String())
				fmt.Printf("ether:\t%s\n", weiToEther(weiBigInt).Text('f', -1))
			}
			if ether.Cmp(big.NewFloat(0)) == 1 {
				weiBigInt := etherToWei(ether)
				fmt.Printf("wei:\t%s\n", weiBigInt.String())
				fmt.Printf("gwei:\t%s\n", weiToGwei(weiBigInt).Text('f', -1))
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
				eth := weiToEther(wei)
				// Update ether field
				etherInput.SetText(eth.Text('f', -1))
				// Update gwei field
				gweiInput.SetText(weiToGwei(wei).Text('f', -1))
			}
		})

		gweiInput.SetChangedFunc(func(gweiString string) {
			// Only update the other field if it's this field that has the focus
			if gweiInput.HasFocus() {
				gwei, _ := ParseBigFloat(gweiString)
				wei := gweiTowei(gwei)
				// Update wei field
				weiInput.SetText(wei.String())
				eth := weiToEther(wei)
				// Update ether field
				etherInput.SetText(eth.Text('f', -1))
			}
		})

		etherInput.SetChangedFunc(func(etherString string) {
			// Only update the other field if it's this field that has the focus
			if etherInput.HasFocus() {
				ether, _ := ParseBigFloat(etherString)
				wei := etherToWei(ether)
				// Update wei field
				weiInput.SetText(wei.String())
				// Update gwei field
				gweiInput.SetText(weiToGwei(wei).Text('f', -1))
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

func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

func weiToGwei(wei *big.Int) *big.Float {
	//return new(big.Int).Div(wei, big.NewInt(params.GWei))
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.GWei))
	//return wei
}

func gweiTowei(gwei *big.Float) *big.Int {
	//return new(big.Int).Mul(gwei, big.NewInt(params.GWei))
	truncInt, _ := gwei.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.GWei))
	fracStr := strings.Split(fmt.Sprintf("%.18f", gwei), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func etherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

// ParseBigFloat parse string value to big.Float
func ParseBigFloat(value string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	return f, err
}
