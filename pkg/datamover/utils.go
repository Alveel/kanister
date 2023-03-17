package datamover

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/kanisterio/kanister/pkg/field"
	"github.com/kanisterio/kanister/pkg/kopia"
	kopiacmd "github.com/kanisterio/kanister/pkg/kopia/command"
	"github.com/kanisterio/kanister/pkg/kopia/repository"
	"github.com/kanisterio/kanister/pkg/kopia/snapshot"
	"github.com/kanisterio/kanister/pkg/location"
	"github.com/kanisterio/kanister/pkg/log"
	"github.com/kanisterio/kanister/pkg/output"
	"github.com/kanisterio/kanister/pkg/param"
)

const (
	usePipeParam = `-`
)

func targetWriter(target string) (io.Writer, error) {
	if target != usePipeParam {
		return os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0755)
	}
	return os.Stdout, nil
}

func locationPull(ctx context.Context, p *param.Profile, path string, target io.Writer) error {
	return location.Read(ctx, target, *p, path)
}

// kopiaLocationPull pulls the data from a kopia snapshot into the given target
func kopiaLocationPull(ctx context.Context, backupID, path, targetPath, password string) error {
	switch targetPath {
	case usePipeParam:
		return snapshot.Read(ctx, os.Stdout, backupID, path, password)
	default:
		return snapshot.ReadFile(ctx, backupID, targetPath, password)
	}
}

// connectToKopiaServer connects to the kopia server with given creds
func connectToKopiaServer(ctx context.Context, kp *param.Profile) error {
	contentCacheSize := kopia.GetDataStoreGeneralContentCacheSize(kp.Credential.KopiaServerSecret.ConnectOptions)
	metadataCacheSize := kopia.GetDataStoreGeneralMetadataCacheSize(kp.Credential.KopiaServerSecret.ConnectOptions)
	return repository.ConnectToAPIServer(
		ctx,
		kp.Credential.KopiaServerSecret.Cert,
		kp.Credential.KopiaServerSecret.Password,
		kp.Credential.KopiaServerSecret.Hostname,
		kp.Location.Endpoint,
		kp.Credential.KopiaServerSecret.Username,
		contentCacheSize,
		metadataCacheSize,
	)
}

// connectToKopiaRepositoryServer connects to the kopia server with given repository server CR
func connectToKopiaRepositoryServer(ctx context.Context, rs *param.RepositoryServer) (string, error) {
	contentCacheMB, metadataCacheMB := kopiacmd.GetCacheSizeSettingsForSnapshot()
	hostname, userPassphrase, certData, err := secretsFromRepositoryServerCR(rs)
	if err != nil {
		return "", errors.Wrap(err, "Error Retrieving Connection Data from Repository Server")
	}
	return userPassphrase, repository.ConnectToAPIServer(
		ctx,
		certData,
		userPassphrase,
		hostname,
		rs.Address,
		rs.Username,
		contentCacheMB,
		metadataCacheMB,
	)
}

func secretsFromRepositoryServerCR(rs *param.RepositoryServer) (string, string, string, error) {
	userCredJSON, err := json.Marshal(rs.Credentials.ServerUserAccess.Data)
	if err != nil {
		return "", "", "", errors.Wrap(err, "Error Unmarshalling User Credentials")
	}
	certJSON, err := json.Marshal(rs.Credentials.ServerTLS.Data)
	if err != nil {
		return "", "", "", errors.Wrap(err, "Error Unmarshalling Certificate")
	}
	hostname, userPassphrase, err := hostNameAndUserPassPhraseFromRepoServer(string(userCredJSON))
	if err != nil {
		return "", "", "", errors.Wrap(err, "Error Getting Hostname/User Passphrase from User credentials")
	}
	return hostname, userPassphrase, string(certJSON), err
}

func hostNameAndUserPassPhraseFromRepoServer(userCreds string) (string, string, error) {
	var userAccessMap map[string]string
	if err := json.Unmarshal([]byte(userCreds), &userAccessMap); err != nil {
		return "", "", errors.Wrap(err, "Failed to unmarshal User Credentials Data")
	}

	var userPassPhrase string
	var hostName string
	for key, val := range userAccessMap {
		hostName = key
		userPassPhrase = val
	}
	decodedUserPassPhrase, err := base64.StdEncoding.DecodeString(userPassPhrase)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to Decode User Passphrase")
	}
	return hostName, string(decodedUserPassPhrase), nil

}

// kopiaLocationPush pushes the data from the source using a kopia snapshot
func kopiaLocationPush(ctx context.Context, path, outputName, sourcePath, password string) error {
	var snapInfo *snapshot.SnapshotInfo
	var err error
	switch sourcePath {
	case usePipeParam:
		snapInfo, err = snapshot.Write(ctx, os.Stdin, path, password)
		log.Print("---- PipeParam ----", field.M{
			"Source Path":      os.Stdin,
			"Destination Path": path,
		})
	default:
		snapInfo, err = snapshot.WriteFile(ctx, path, sourcePath, password)
	}
	if err != nil {
		return errors.Wrap(err, "Failed to push data using kopia")
	}
	log.Print("---- Paths ----", field.M{
		"Source Path":      sourcePath,
		"Destination Path": path,
	})
	snapInfoJSON, err := snapshot.MarshalKopiaSnapshot(snapInfo)
	if err != nil {
		return err
	}

	return output.PrintOutput(outputName, snapInfoJSON)
}

func sourceReader(source string) (io.Reader, error) {
	if source != usePipeParam {
		return os.Open(source)
	}
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Stat stdin")
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return nil, errors.New("Stdin must be piped when the source parameter is \"-\"")
	}
	return os.Stdin, nil
}

func locationPush(ctx context.Context, p *param.Profile, path string, source io.Reader) error {
	return location.Write(ctx, source, *p, path)
}

// kopiaLocationDelete deletes the kopia snapshot with given backupID
func kopiaLocationDelete(ctx context.Context, backupID, path, password string) error {
	return snapshot.Delete(ctx, backupID, path, password)
}

func locationDelete(ctx context.Context, p *param.Profile, path string) error {
	return location.Delete(ctx, *p, path)
}