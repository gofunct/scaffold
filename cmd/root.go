// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gorpc",
		Short: "A generator for Grpc based Applications",
		Long: `GoRPC is a library for Go that empowers grpc applications.
This application is a tool to generate the needed files
to quickly create a GRPC application.`,
	}
)

// Execute executes the root command.
func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(AddCmd)
	rootCmd.AddCommand(InitCmd)
}