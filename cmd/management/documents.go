package management

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDocumentsCmd(getBaseURL func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "documents",
		Short: "Document management operations",
	}
	cmd.AddCommand(newDocumentPushCmd(getBaseURL))
	return cmd
}

func newDocumentPushCmd(getBaseURL func() string) *cobra.Command {
	var uid, docType, docNumber, filename, fileExt, upload, docExpire, metadata string
	var updateLabels bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push a document to Barong",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newManagementClient(getBaseURL())
			if err != nil {
				return err
			}
			if err := client.PushDocument(uid, docType, docNumber, filename, fileExt, upload, docExpire, updateLabels, metadata); err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Document pushed successfully")
			return nil
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "User UID (required)")
	cmd.Flags().StringVar(&docType, "doc-type", "", "Document type (required)")
	cmd.Flags().StringVar(&docNumber, "doc-number", "", "Document number (required)")
	cmd.Flags().StringVar(&filename, "filename", "", "Document filename (required)")
	cmd.Flags().StringVar(&fileExt, "file-ext", "", "Document file extension (required)")
	cmd.Flags().StringVar(&upload, "upload", "", "Base64 encoded document content (required)")
	cmd.Flags().StringVar(&docExpire, "doc-expire", "", "Document expiration date")
	cmd.Flags().BoolVar(&updateLabels, "update-labels", true, "Create or update user label (default true)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "Additional metadata as JSON string")
	_ = cmd.MarkFlagRequired("uid")
	_ = cmd.MarkFlagRequired("doc-type")
	_ = cmd.MarkFlagRequired("doc-number")
	_ = cmd.MarkFlagRequired("filename")
	_ = cmd.MarkFlagRequired("file-ext")
	_ = cmd.MarkFlagRequired("upload")
	return cmd
}
