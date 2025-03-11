package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	cryptossh "golang.org/x/crypto/ssh"
)

func CloneRepo(repositoryURL string, outputFolder string, privateKey []byte, depth int) (err error) {
	// manage authentication
	var gitAuth transport.AuthMethod
	if strings.HasPrefix(repositoryURL, "http") {
		gitAuth = nil // TODO: support for authentication with token
	} else {
		gitAuth, err = ssh.NewPublicKeys("git", privateKey, "")
		if err != nil {
			return fmt.Errorf("git authentication failure %s", err)
		}
		gitAuth.(*ssh.PublicKeys).HostKeyCallback = cryptossh.InsecureIgnoreHostKey()
	}

	_, err = git.PlainClone(outputFolder, false, &git.CloneOptions{
		URL:             repositoryURL,
		InsecureSkipTLS: true,
		Depth:           depth, // retrieve only latest commit
		SingleBranch:    true,
		Auth:            gitAuth,
	})

	if err != nil {
		if strings.HasPrefix(repositoryURL, "http") {
			return fmt.Errorf("failed to clone remote repository %s", err)
		} else {
			return fmt.Errorf("have you added the Codebox SSH public key to the remote Git server? %s", err)
		}
	}

	return nil
}
