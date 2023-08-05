package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-go-golems/go-go-labs/cmd/mastoid/pkg"
	"github.com/go-go-golems/go-go-labs/cmd/mastoid/pkg/render"
	"github.com/go-go-golems/go-go-labs/cmd/mastoid/pkg/render/html"
	"github.com/go-go-golems/go-go-labs/cmd/mastoid/pkg/render/plaintext"
	"github.com/mattn/go-mastodon"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var ThreadCmd = &cobra.Command{
	Use:   "thread",
	Short: "Retrieves a thread from a Mastodon instance",
	Run: func(cmd *cobra.Command, args []string) {
		statusID, _ := cmd.Flags().GetString("status-id")
		verbose, _ := cmd.Flags().GetBool("verbose")
		withHtml, _ := cmd.Flags().GetBool("withHtml")
		withJson, _ := cmd.Flags().GetBool("json")
		withHeader, _ := cmd.Flags().GetBool("with-header")

		// extract statusID from URL if we have a URL
		if strings.Contains(statusID, "http") {
			statusID = strings.Split(statusID, "/")[4]
		}

		ctx := context.Background()

		credentials, err := pkg.LoadCredentials()
		cobra.CheckErr(err)

		client, err := pkg.CreateClientAndAuthenticate(ctx, credentials)
		cobra.CheckErr(err)

		status, err := client.GetStatus(ctx, mastodon.ID(statusID))
		if err != nil {
			log.Error().Err(err).Str("statusId", statusID).Msg("Could not get status")
		}
		cobra.CheckErr(err)

		context, err := client.GetStatusContext(ctx, status.ID)
		cobra.CheckErr(err)

		thread := &pkg.Thread{
			Nodes: map[mastodon.ID]*pkg.Node{},
		}

		thread.AddStatus(status)

		thread.AddContextAndGetMissingIDs(status.ID, context)

		printNode := func(node *pkg.Node, depth int) error {
			fmt.Printf(
				"%s%s (parent: %s)",
				strings.Repeat("  ", depth), node.Status.ID, node.Status.InReplyToID,
			)
			if node.Status != nil {
				// print the first 20 characters of the status
				l := len(node.Status.Content)
				if l > 20 {
					l = 20
				}
				fmt.Printf(" %s", node.Status.Content[0:l])
			}
			fmt.Println()
			return nil
		}
		err = thread.WalkBreadthFirst(printNode)
		cobra.CheckErr(err)

		err = thread.WalkDepthFirst(printNode)

		return

		if withJson {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			err = encoder.Encode(context)
			cobra.CheckErr(err)
			return
		}

		var renderer render.Renderer

		if withHtml {
			renderer = html.NewRenderer(
				html.WithVerbose(verbose),
				html.WithHeader(withHeader),
			)
		} else {
			renderer = plaintext.NewRenderer(
				plaintext.WithVerbose(verbose),
				plaintext.WithHeader(withHeader),
			)
		}

		err = renderer.RenderThread(os.Stdout, status, context)
		cobra.CheckErr(err)
	},
}