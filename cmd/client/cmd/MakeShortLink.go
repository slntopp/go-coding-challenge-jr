/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"challenge/pkg/api"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"github.com/spf13/cobra"
)

// MakeShortLinkCmd represents the MakeShortLink command
var MakeShortLinkCmd = &cobra.Command{
	Use:   "MakeShortLink",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MakeShortLink called")
		if len(args) != 1 {
			log.Fatal("incorrect number of arguments")
		}

		longURL := args[0]

		conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}

		c := api.NewChallengeServiceClient(conn)
		resp, err := c.MakeShortLink(context.Background(), &api.Link{Data: longURL})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("short link: %s", resp.GetData())
	},
}

func init() {
	rootCmd.AddCommand(MakeShortLinkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// MakeShortLinkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// MakeShortLinkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
