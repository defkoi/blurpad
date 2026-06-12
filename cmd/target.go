/*
Copyright © 2026 defkoi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"blurpad/internal"
	"blurpad/lib"
	"fmt"
	"image"
	"strings"

	"github.com/spf13/cobra"
)

const targetLong = //
`Supported targets:
  instagram`

// targetCmd represents the target command
var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "Pad for target ratio",
	Long:  targetLong,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Undefined target.")
			return
		}

		args[0] = strings.ToLower(args[0])
		target, ok := internal.Targets[args[0]]
		if !ok {
			fmt.Println("Unknown target.")
			return
		}

		lib.OpenDoSave(
			inputFile, outputFile,
			func(src image.Image) (image.Image, error) {
				pad := internal.PaddingFromRatio(src, target)
				return internal.Process(src, pad), nil
			},
		)
	},
}

func init() {
	rootCmd.AddCommand(targetCmd)
}
