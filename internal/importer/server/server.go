package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	library_proto "github.com/noitcelfer7/library-proto/gen/go/proto/library"

	"library_importer/internal/importer/config"
	"library_importer/internal/importer/sqlite"
)

func Serve(config *config.Config, cc library_proto.DataExchangeServiceClient, ctx context.Context) {
	http.HandleFunc("/upload", uploadHandler(cc, ctx))

	addr := net.JoinHostPort(config.Http.Server.Host, config.Http.Server.Port)

	http.ListenAndServe(addr, nil)
}

func uploadHandler(cc library_proto.DataExchangeServiceClient, ctx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "uploadHandler Error: request != http.MethodPost", http.StatusMethodNotAllowed)

			return
		}

		file, _, err := request.FormFile("file")

		if err != nil {
			http.Error(writer, fmt.Sprintf("request.FormFile Error: %v", err), http.StatusBadRequest)

			return
		}

		defer file.Close()

		temp, err := os.CreateTemp("", "temp-*")

		if err != nil {
			http.Error(writer, fmt.Sprintf("os.CreateTemp Error: %v", err), http.StatusInternalServerError)

			return
		}

		defer temp.Close()

		_, err = io.Copy(temp, file)

		if err != nil {
			http.Error(writer, fmt.Sprintf("io.Copy Error: %v", err), http.StatusInternalServerError)

			return
		}

		records, err := sqlite.Parse(temp.Name())

		if err != nil {
			http.Error(writer, fmt.Sprintf("sqlite.Parse Error: %v", err), http.StatusInternalServerError)

			return
		}

		for _, record := range records {
			issueReturnDate := ""

			if record.IssueReturnDate.Valid {
				issueReturnDate = record.IssueReturnDate.String
			}

			response, err := cc.Exchange(ctx,
				&library_proto.ExchangeRequest{
					AuthorFirstName: record.AuthorFirstName,
					AuthorLastName:  record.AuthorLastName,

					BookIsbn:  record.BookIsbn,
					BookTitle: record.BookTitle,

					GenreTitle: record.GenreTitle,

					IssueDate:       record.IssueDate,
					IssuePeriod:     record.IssuePeriod,
					IssueReturnDate: &issueReturnDate,

					ReaderFirstName:   record.ReaderFirstName,
					ReaderLastName:    record.ReaderLastName,
					ReaderPhoneNumber: record.ReaderPhoneNumber,
				})

			if err != nil {
				http.Error(writer, fmt.Sprintf("cc.Send Error: %v", err), http.StatusInternalServerError)

				return
			}

			if !response.IsSuccessful {
				http.Error(writer, "cc.Send Error: !response.IsSuccessful", http.StatusInternalServerError)

				return
			}
		}
	}
}
