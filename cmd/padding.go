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

	"github.com/spf13/cobra"
)

// paddingCmd represents the padding command
var paddingCmd = &cobra.Command{
	Use:   "padding",
	Short: "Pad for CSS padding",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Undefined padding.")
			return
		}

		pad, err := internal.ParsePadding(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := lib.OpenDoSave(
			inputFile, outputFile,
			func(src image.Image) (image.Image, error) {
				return internal.Process(src, pad), nil
			}); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(paddingCmd)
}
