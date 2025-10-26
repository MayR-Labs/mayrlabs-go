package commands

import (
	"fmt"
	"time"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// AddLicenseCmd creates a LICENSE file
var AddLicenseCmd = &cobra.Command{
	Use:   "add-license",
	Short: "Create a LICENSE file based on selected type, year, and author",
	Long:  "Generate a LICENSE file for your project with popular open source licenses",
	RunE: func(cmd *cobra.Command, args []string) error {
		licenseType, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		}
		author, err := cmd.Flags().GetString("author")
		if err != nil {
			return err
		}
		year, err := cmd.Flags().GetString("year")
		if err != nil {
			return err
		}
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		// Check if LICENSE already exists
		if utils.FileExists("LICENSE") && !force {
			return fmt.Errorf("LICENSE file already exists. Use --force to overwrite")
		}

		// Prompt for missing values
		if licenseType == "" {
			fmt.Println("Available licenses: mit, apache2, gpl3, bsd3")
			licenseType, err = utils.PromptInput("Select license type: ")
			if err != nil {
				return err
			}
		}

		if author == "" {
			author, err = utils.PromptInput("Enter author name: ")
			if err != nil {
				return err
			}
		}

		if year == "" {
			year = fmt.Sprintf("%d", time.Now().Year())
		}

		// Generate license content
		content, err := generateLicense(licenseType, author, year)
		if err != nil {
			return err
		}

		// Write to file
		if err := utils.WriteFile("LICENSE", content); err != nil {
			return fmt.Errorf("failed to write LICENSE file: %w", err)
		}

		fmt.Printf("✅ LICENSE file created successfully! (%s, %s, %s)\n", licenseType, author, year)
		return nil
	},
}

func init() {
	AddLicenseCmd.Flags().StringP("type", "t", "", "License type (mit, apache2, gpl3, bsd3)")
	AddLicenseCmd.Flags().StringP("author", "a", "", "Author name")
	AddLicenseCmd.Flags().StringP("year", "y", "", "Copyright year (default: current year)")
	AddLicenseCmd.Flags().BoolP("force", "f", false, "Overwrite existing LICENSE file")
}

func generateLicense(licenseType, author, year string) (string, error) {
	switch licenseType {
	case "mit":
		return fmt.Sprintf(`MIT License

Copyright (c) %s %s

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`, year, author), nil

	case "apache2":
		return fmt.Sprintf(`Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

Copyright %s %s

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`, year, author), nil

	case "gpl3":
		return fmt.Sprintf(`GNU GENERAL PUBLIC LICENSE
Version 3, 29 June 2007

Copyright (C) %s %s

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
`, year, author), nil

	case "bsd3":
		return fmt.Sprintf(`BSD 3-Clause License

Copyright (c) %s, %s
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`, year, author), nil

	default:
		return "", fmt.Errorf("unsupported license type: %s", licenseType)
	}
}
