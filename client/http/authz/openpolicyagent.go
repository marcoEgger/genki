package authz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/marcoEgger/genki/logger"
	"github.com/marcoEgger/genki/metadata"
)

const (
	OpenPolicyAgentUrlConfigKey = "opa-url"
	OpenPolicyAgentDefaultUrl   = "http://localhost:8181/v1/data/authz/"
)

type opa struct {
	url string
}

func NewOpenPolicyAgentClient(opaUrl string) *opa {
	return &opa{url: opaUrl}
}

func (auth *opa) Authorize(ctx context.Context, resourceId, action interface{}, externalData interface{}) error {
	log := logger.WithMetadata(ctx)
	if auth.url == "" {
		return fmt.Errorf("empty opa url, cannot query policy")
	}

	payload := map[string]interface{}{
		"resource": resourceId,
		"action":   action,
		"external": externalData,
	}

	if metadata.GetFromContext(ctx, metadata.InternalKey) == "true" {
		payload["internal"] = true
	}

	if metadata.GetFromContext(ctx, metadata.M2MKey) == "true" {
		payload["m2m"] = true
		payload["accounts"] = strings.Split(metadata.GetFromContext(ctx, metadata.AccountIDsKey), ",")
		payload["roles"] = strings.Split(metadata.GetFromContext(ctx, metadata.RolesKey), ",")
	} else if metadata.GetFromContext(ctx, metadata.UserIDKey) != "" {
		payload["user"] = metadata.GetFromContext(ctx, metadata.UserIDKey)
		payload["account"] = metadata.GetFromContext(ctx, metadata.AccountIDKey)
		payload["roles"] = metadata.GetFromContext(ctx, metadata.RolesKey)
		payload["type"] = metadata.GetFromContext(ctx, metadata.TypeKey)
		payload["subType"] = metadata.GetFromContext(ctx, metadata.SubTypeKey)
		payload["email"] = metadata.GetFromContext(ctx, metadata.EmailKey)
	} else {
		payload["roles"] = []string{"anonymous"}
	}

	input := map[string]interface{}{
		"input": payload,
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		return errors.Wrap(err, "unable to marshal payload to JSON")
	}

	log.Debugf("sending authorization request: %v", string(jsonData))
	reqStart := time.Now()
	resp, err := http.Post(auth.url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "Authorize failed")
	}
	defer resp.Body.Close()

	responseJSON := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&responseJSON); err != nil {
		return errors.Wrap(err, "unable to decode json in response body")
	}
	log.Debugf("received authorization response after %v: %v", time.Since(reqStart), responseJSON)

	if data, ok := responseJSON["result"]; ok {
		if data == true {
			return nil
		}
		return fmt.Errorf("request is not authorized")
	}

	return fmt.Errorf("authorization response did not contain a result")
}

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("openpolicyagent", pflag.ContinueOnError)
	fs.String(
		OpenPolicyAgentUrlConfigKey,
		OpenPolicyAgentDefaultUrl,
		"http url to the OPA instance for this service (sidecar)",
	)

	return fs
}
